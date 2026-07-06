package mirg

import (
	"context"
	"fmt"
	"testing"

	"github.com/olivere/elastic/v7"
)

func TestGoodsStockUpdates(t *testing.T) {
	ctx := context.Background()
	client := NewEs()

	const (
		index      = "goods_stock"
		merchantID = 0
		code       = "123456654321"
	)

	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("merchant_id", merchantID),
		elastic.NewTermQuery("code", code),
	)

	searchRes, err := client.Search(index).Query(query).Size(10).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Found docs: %d\n", searchRes.TotalHits())

	if searchRes.TotalHits() == 0 {
		return
	}

	script := elastic.NewScript(`ctx._source.status = params.status`).
		Lang("painless").
		Param("status", 2)

	updateRes, err := client.UpdateByQuery(index).
		Query(query).
		Script(script).
		Refresh("true").
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Updated docs: %d\n", updateRes.Updated)
}
