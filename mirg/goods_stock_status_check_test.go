package mirg

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/olivere/elastic/v7"
)

// TestCompareGoodsStockStatus 仅核对 zby_goods_stock 和 ES goods_stock 的 status，
// 不会修改 MySQL 或 Elasticsearch 中的数据。
//
// 使用 ES 文档 _source.id 与 zby_goods_stock.id 关联；不依赖 ES 的 _id。
// 范围固定为 2026-07-01 00:00:00（含）至 2026-07-15 23:00:00（不含），且 stock_type=1。
//
// 执行：go test -run '^TestCompareGoodsStockStatus$' -v -count=1
// 可选环境变量：
//
//	GOODS_STOCK_CHECK_OUTPUT_DIR  输出目录，默认 ./output/goods_stock_status_check
//	GOODS_STOCK_CHECK_RESET       设为 1 时忽略已有断点重新开始；建议同时使用新的输出目录
func TestCompareGoodsStockStatus(t *testing.T) {
	ctx := context.Background()
	const batchSize = 2000
	const startAt = "2026-07-23 00:00:00"
	const endAt = "2026-07-23 23:59:59"
	outputDir := goodsStockCheckOutputDir(t)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("create output directory: %v", err)
	}

	checkpointPath := filepath.Join(outputDir, "checkpoint.json")
	checkpoint, exists := readGoodsStockCheckpoint(t, checkpointPath)
	if os.Getenv("GOODS_STOCK_CHECK_RESET") == "1" {
		exists = false
	}
	if !exists {
		checkpoint = goodsStockCheckpoint{}
	}

	resultPath := filepath.Join(outputDir, "status_mismatch.csv")
	resultFile, err := os.OpenFile(resultPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatalf("open result file: %v", err)
	}
	defer resultFile.Close()

	writer := csv.NewWriter(resultFile)
	defer writer.Flush()
	fileInfo, err := resultFile.Stat()
	if err != nil {
		t.Fatalf("stat result file: %v", err)
	}
	if fileInfo.Size() == 0 {
		if err := writer.Write([]string{"id", "db_status", "es_status", "reason"}); err != nil {
			t.Fatalf("write CSV header: %v", err)
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			t.Fatalf("flush CSV header: %v", err)
		}
	}

	var config orm.Config
	configurator.NewYaml().MustLoad("./config/online/db_read.yaml", &config)
	readDB, err := orm.NewMySQLGorm(&config)
	if err != nil {
		t.Fatalf("connect read database: %v", err)
	}
	esClient := NewEs()

	var scanned, missing, mismatch, duplicate int64
	processBatch := func(rows []goodsStockStatusRow) {
		esStatuses, duplicateIDs, err := getGoodsStockESStatuses(ctx, esClient, rows)
		if err != nil {
			t.Fatalf("read ES for DB id %d-%d: %v", rows[0].ID, rows[len(rows)-1].ID, err)
		}
		for _, row := range rows {
			dbStatus := nullInt64CSV(row.Status)
			if duplicateIDs[row.ID] {
				duplicate++
				writeGoodsStockResult(t, writer, row.ID, dbStatus, "", "duplicate_in_es")
				continue
			}
			esStatus, found := esStatuses[row.ID]
			if !found {
				missing++
				writeGoodsStockResult(t, writer, row.ID, dbStatus, "", "missing_in_es")
				continue
			}
			if dbStatus != esStatus {
				mismatch++
				writeGoodsStockResult(t, writer, row.ID, dbStatus, esStatus, "status_mismatch")
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			t.Fatalf("flush CSV: %v", err)
		}
		lastRow := rows[len(rows)-1]
		checkpoint.LastUpdateTime = lastRow.UpdateTime
		checkpoint.LastID = lastRow.ID
		if err := writeGoodsStockCheckpoint(checkpointPath, checkpoint); err != nil {
			t.Fatalf("write checkpoint: %v", err)
		}
		scanned += int64(len(rows))
		if scanned%(int64(batchSize)*100) == 0 {
			t.Logf("progress: scanned=%d, last_update_time=%d, last_id=%d, missing=%d, mismatch=%d, duplicate=%d", scanned, checkpoint.LastUpdateTime, checkpoint.LastID, missing, mismatch, duplicate)
		}
		// 每批 ES 查询完成后暂停 300ms，避免连续请求对 ES 造成压力。
		time.Sleep(500 * time.Millisecond)
	}

	// 这里只执行一次 MySQL SELECT。Rows() 保持结果集游标打开，由 Go 逐行读取；
	// 每累计 500 条再批量查询一次 ES，因此不会将 260 多万条记录全部放进内存。
	dbRows, err := readDB.Table("zby_goods_stock").
		Select("id, status, update_time").
		//Where("id in (78277274,78277610,78277768,78278239,72833216) ").
		Where("merchant_id > 0 AND stock_type = 1 AND update_time >= UNIX_TIMESTAMP(?) AND update_time < UNIX_TIMESTAMP(?)", startAt, endAt).
		Where("(update_time > ? OR (update_time = ? AND id > ?))", checkpoint.LastUpdateTime, checkpoint.LastUpdateTime, checkpoint.LastID).
		Order("update_time ASC, id ASC").
		Rows()
	if err != nil {
		t.Fatalf("query database after update_time=%d, id=%d: %v", checkpoint.LastUpdateTime, checkpoint.LastID, err)
	}
	defer dbRows.Close()

	batch := make([]goodsStockStatusRow, 0, batchSize)
	for dbRows.Next() {
		var row goodsStockStatusRow
		if err := readDB.ScanRows(dbRows, &row); err != nil {
			t.Fatalf("scan database row: %v", err)
		}
		batch = append(batch, row)
		if len(batch) == batchSize {
			processBatch(batch)
			batch = batch[:0]
		}
	}
	if err := dbRows.Err(); err != nil {
		t.Fatalf("iterate database rows: %v", err)
	}
	if len(batch) > 0 {
		processBatch(batch)
	}
	t.Logf("completed: scanned=%d, missing=%d, mismatch=%d, duplicate=%d, range=[%s,%s), result=%s", scanned, missing, mismatch, duplicate, startAt, endAt, resultPath)
}

type goodsStockStatusRow struct {
	ID         int64         `gorm:"column:id"`
	Status     sql.NullInt64 `gorm:"column:status"`
	UpdateTime int64         `gorm:"column:update_time"`
}

type goodsStockCheckpoint struct {
	LastUpdateTime int64 `json:"last_update_time"`
	LastID         int64 `json:"last_id"`
}

func getGoodsStockESStatuses(ctx context.Context, client *es.ES, rows []goodsStockStatusRow) (map[int64]string, map[int64]bool, error) {
	ids := make([]interface{}, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	response, err := client.Search("goods_stock").
		Query(elastic.NewTermsQuery("id", ids...)).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("id", "status")).
		Size(len(rows) * 2).
		Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	statuses := make(map[int64]string, len(rows))
	duplicateIDs := make(map[int64]bool)
	for _, hit := range response.Hits.Hits {
		var source struct {
			ID     json.RawMessage `json:"id"`
			Status json.RawMessage `json:"status"`
		}
		if err := json.Unmarshal(hit.Source, &source); err != nil {
			return nil, nil, fmt.Errorf("unmarshal ES document %s: %w", hit.Id, err)
		}
		id, err := parseGoodsStockESID(source.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("read ES document %s id: %w", hit.Id, err)
		}
		if _, alreadyExists := statuses[id]; alreadyExists {
			duplicateIDs[id] = true
			continue
		}
		statuses[id] = strings.TrimSpace(string(source.Status))
	}
	return statuses, duplicateIDs, nil
}

func goodsStockCheckOutputDir(t *testing.T) string {
	dir := os.Getenv("GOODS_STOCK_CHECK_OUTPUT_DIR")
	if dir == "" {
		dir = "./output/goods_stock_status_check"
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("resolve output directory: %v", err)
	}
	return absDir
}

func readGoodsStockCheckpoint(t *testing.T, checkpointPath string) (goodsStockCheckpoint, bool) {
	content, err := os.ReadFile(checkpointPath)
	if os.IsNotExist(err) {
		return goodsStockCheckpoint{}, false
	}
	if err != nil {
		t.Fatalf("read checkpoint: %v", err)
	}
	var checkpoint goodsStockCheckpoint
	if err := json.Unmarshal(content, &checkpoint); err != nil || checkpoint.LastUpdateTime < 0 || checkpoint.LastID < 0 {
		t.Fatalf("invalid checkpoint %s", strings.TrimSpace(string(content)))
	}
	return checkpoint, true
}

func writeGoodsStockCheckpoint(checkpointPath string, checkpoint goodsStockCheckpoint) error {
	content, err := json.Marshal(checkpoint)
	if err != nil {
		return err
	}
	return os.WriteFile(checkpointPath, append(content, '\n'), 0o644)
}

func parseGoodsStockESID(raw json.RawMessage) (int64, error) {
	var id int64
	if err := json.Unmarshal(raw, &id); err == nil {
		return id, nil
	}
	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return 0, fmt.Errorf("expected numeric id, got %s", string(raw))
	}
	return strconv.ParseInt(text, 10, 64)
}

func nullInt64CSV(value sql.NullInt64) string {
	if !value.Valid {
		return "null"
	}
	return strconv.FormatInt(value.Int64, 10)
}

func writeGoodsStockResult(t *testing.T, writer *csv.Writer, id int64, dbStatus, esStatus, reason string) {
	if err := writer.Write([]string{strconv.FormatInt(id, 10), dbStatus, esStatus, reason}); err != nil {
		t.Fatalf("write CSV result: %v", err)
	}
}
