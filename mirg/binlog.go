package mirg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func Binlog() {
	// 打开文件
	filePath := "/Users/linlong/Desktop/data2/recovery.log"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return
	}
	defer file.Close()

	// 初始化计数器
	count := 0
	var list = make([][]string, 0)
	var item = make([]string, 0)

	// 使用bufio.Reader逐行读取，这样可以处理超长行
	reader := bufio.NewReader(file)
	var line strings.Builder

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// 处理文件末尾的最后一行
				if line.Len() > 0 {
					currentLine := line.String()
					// 移除可能的换行符
					currentLine = strings.TrimSuffix(currentLine, "\r")
					currentLine = strings.TrimSuffix(currentLine, "\n")

					if strings.HasPrefix(currentLine, "### DELETE FROM `zhubaoe`.`zby_old_stock`") {
						count++
						if len(item) > 0 {
							list = append(list, item)
						}
						item = make([]string, 0)
					} else {
						item = append(item, currentLine)
					}
				}
				break
			} else {
				fmt.Printf("读取文件时出错: %v\n", err)
				return
			}
		}

		// 检查是否是行结束符
		if char == '\n' {
			currentLine := line.String()
			// 移除可能的回车符
			currentLine = strings.TrimSuffix(currentLine, "\r")

			if strings.HasPrefix(currentLine, "### DELETE FROM `zhubaoe`.`zby_old_stock`") {
				count++

				if len(item) > 0 {
					list = append(list, item)
				}
				item = make([]string, 0)
			} else {
				item = append(item, currentLine)
			}

			// 清空行缓冲区
			line.Reset()
		} else {
			// 将字符添加到当前行
			line.WriteRune(char)
		}
	}

	// 添加最后一组数据（如果有的话）
	if len(item) > 0 {
		list = append(list, item)
	}

	// 打印结果
	fmt.Printf("匹配到的DELETE FROM语句数量: %d\n", count)
	fmt.Printf("分组数量: %d\n", len(list))

	datas := make([]*ZbyOldStock, 0)

	//如果需要查看分组内容，可以取消下面的注释
	for _, group := range list {
		if len(group) == 82 {
			d := aaa(group)
			datas = append(datas, d)
		}

		if len(group) != 82 {
			d := aaa(group[:82])
			datas = append(datas, d)
		}
	}

	//var configPath = "./config/online/db.yaml"
	//var config orm.Config
	//configurator.NewYaml().MustLoad(configPath, &config)
	//db, err := orm.NewMySQLGorm(&config)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for index, d := range datas {
	//	fmt.Println(index)
	//
	//	err := db.Table("zby_old_stock").Create(d).Error
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
}

type ZbyOldStock struct {
	ID                        int     `json:"id"`
	GoodsTypeID               *string `json:"goods_type_id"`
	GoodsTypeIDOld            *string `json:"goods_type_id_old"`
	GoodsTypeName             *string `json:"goods_type_name"`
	MerchantID                *int    `json:"merchant_id"`
	MerchantShopID            *int    `json:"merchant_shop_id"`
	AdminUserID               *int    `json:"admin_user_id"`
	AdminUserName             *string `json:"admin_user_name"`
	ExportOrderSn             *string `json:"export_order_sn"`
	ClassifyID                *int    `json:"classify_id"`
	ClassifyIDOld             *int    `json:"classify_id_old"`
	ClassifyName              *string `json:"classify_name"`
	QualityID                 *int    `json:"quality_id"`
	QualityIDOld              *int    `json:"quality_id_old"`
	QualityName               *string `json:"quality_name"`
	BrandID                   *int    `json:"brand_id"`
	BrandIDOld                *int    `json:"brand_id_old"`
	BrandName                 *string `json:"brand_name"`
	IsOriginal                *int    `json:"is_original"`
	IsOriginalType            *int    `json:"is_original_type"`
	SourceType                *int    `json:"source_type"`
	SourceGoodsTypeName       *string `json:"source_goods_type_name"`
	SourceOrderID             *int    `json:"source_order_id"`
	SourceOrderSn             *string `json:"source_order_sn"`
	ExchangeGoodsID           *int    `json:"exchange_goods_id"`
	SourceGoodsID             *int    `json:"source_goods_id"`
	GoodsName                 *string `json:"goods_name"`
	CostChargeType            *int    `json:"cost_charge_type"`
	DeliveryChargeType        *int    `json:"delivery_charge_type"`
	MainStoneName             *string `json:"main_stone_name"`
	MainStoneTypeID           *int    `json:"main_stone_type_id"`
	MainStoneTypeName         *string `json:"main_stone_type_name"`
	MainStoneWeight           *string `json:"main_stone_weight"`
	ViceStoneTypeName         *string `json:"vice_stone_type_name"`
	ViceStoneWeight           *string `json:"vice_stone_weight"`
	CerNumber                 *string `json:"cer_number"`
	RecyleGoldPrice           *int    `json:"recyle_gold_price"`
	OriginalGoldPrice         *int    `json:"original_gold_price"`
	StockCode                 *string `json:"stock_code"`
	Weight                    *int    `json:"weight"`
	GoldWeight                *int    `json:"gold_weight"`
	GoldWeightRecycle         *int    `json:"gold_weight_recycle"`
	GoldWeightLossRate        *string `json:"gold_weight_loss_rate"`
	RecycleLabourPrice        *int    `json:"recycle_labour_price"`
	RecycleLabourAmount       *int    `json:"recycle_labour_amount"`
	GoodsNum                  *int    `json:"goods_num"`
	GoodsAmount               *int    `json:"goods_amount"`
	OriginalAmount            *int    `json:"original_amount"`
	OldOriginalLabelPrice     *int    `json:"old_original_label_price"`
	RecycleWeightDiff         *int    `json:"recycle_weight_diff"`
	NewOldDiffDay             *int    `json:"new_old_diff_day"`
	OldOriginalLabelWeight    *int    `json:"old_original_label_weight"`
	OldOriginalCostPrice      *int    `json:"old_original_cost_price"`
	RecycleAmount             *int64  `json:"recycle_amount"`
	DeductAmount              *int    `json:"deduct_amount"`
	CostAmount                *int    `json:"cost_amount"`
	CutNewCostAmount          *int    `json:"cut_new_cost_amount"`
	OldProfitAmount           *int    `json:"old_profit_amount"`
	OldValuedAmount           *int64  `json:"old_valued_amount"`
	GoldType                  *int    `json:"gold_type"`
	ShiftStatus               *int    `json:"shift_status"`
	Remarks                   *string `json:"remarks"`
	Status                    *int    `json:"status"`
	CreateTime                *int    `json:"create_time"`
	UpdateTime                *int    `json:"update_time"`
	DeleteTime                *int    `json:"delete_time"`
	AcountTime                *int    `json:"acount_time"`
	SaleDatetime              *string `json:"sale_datetime"`
	IsOverExchange            *int    `json:"is_over_exchange"`
	PicURL                    *string `json:"pic_url"`
	PicURLList                *string `json:"pic_url_list"`
	OldStockTypeID            *int    `json:"old_stock_type_id"`
	OldStockTypeName          *string `json:"old_stock_type_name"`
	FinancialCutNewCostAmount *int    `json:"financial_cut_new_cost_amount"`
	FinancialCostAmount       *int    `json:"financial_cost_amount"`
	FinancialOldProfitAmount  *int    `json:"financial_old_profit_amount"`
	MarketCutNewCostAmount    *int    `json:"market_cut_new_cost_amount"`
	MarketCostAmount          *int    `json:"market_cost_amount"`
	MarketOldProfitAmount     *int    `json:"market_old_profit_amount"`
	SettleMerchantShopID      *int    `json:"settle_merchant_shop_id"`
	InlaidJewelryAmount       *int    `json:"inlaid_jewelry_amount"`
}

func aaa(list []string) *ZbyOldStock {
	currentRecord := &ZbyOldStock{}

	for _, line := range list {
		if line == "### WHERE" {
			continue
		}

		if strings.HasPrefix(line, "###   @") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			fieldPart := strings.TrimSpace(parts[0])
			valueStr := parts[1]

			if !strings.HasPrefix(fieldPart, "###") {
				continue
			}
			fieldIndexStr := strings.Replace(fieldPart, " ", "", 10)
			fieldIndexStr = strings.Replace(fieldIndexStr, "#", "", 10)
			fieldIndexStr = strings.Replace(fieldIndexStr, "@", "", 10)
			fieldIndex, err := strconv.Atoi(fieldIndexStr)
			if err != nil {
				continue
			}

			switch fieldIndex {
			case 1:
				if id := parseIntValue(valueStr); id != nil {
					currentRecord.ID = *id
				}
			case 2:
				currentRecord.GoodsTypeID = parseValue(valueStr)
			case 3:
				currentRecord.GoodsTypeIDOld = parseValue(valueStr)
			case 4:
				currentRecord.GoodsTypeName = parseValue(valueStr)
			case 5:
				currentRecord.MerchantID = parseIntValue(valueStr)
			case 6:
				currentRecord.MerchantShopID = parseIntValue(valueStr)
			case 7:
				currentRecord.AdminUserID = parseIntValue(valueStr)
			case 8:
				currentRecord.AdminUserName = parseValue(valueStr)
			case 9:
				currentRecord.ExportOrderSn = parseValue(valueStr)
			case 10:
				currentRecord.ClassifyID = parseIntValue(valueStr)
			case 11:
				currentRecord.ClassifyIDOld = parseIntValue(valueStr)
			case 12:
				currentRecord.ClassifyName = parseValue(valueStr)
			case 13:
				currentRecord.QualityID = parseIntValue(valueStr)
			case 14:
				currentRecord.QualityIDOld = parseIntValue(valueStr)
			case 15:
				currentRecord.QualityName = parseValue(valueStr)
			case 16:
				currentRecord.BrandID = parseIntValue(valueStr)
			case 17:
				currentRecord.BrandIDOld = parseIntValue(valueStr)
			case 18:
				currentRecord.BrandName = parseValue(valueStr)
			case 19:
				currentRecord.IsOriginal = parseIntValue(valueStr)
			case 20:
				currentRecord.IsOriginalType = parseIntValue(valueStr)
			case 21:
				currentRecord.SourceType = parseIntValue(valueStr)
			case 22:
				currentRecord.SourceGoodsTypeName = parseValue(valueStr)
			case 23:
				currentRecord.SourceOrderID = parseIntValue(valueStr)
			case 24:
				currentRecord.SourceOrderSn = parseValue(valueStr)
			case 25:
				currentRecord.ExchangeGoodsID = parseIntValue(valueStr)
			case 26:
				currentRecord.SourceGoodsID = parseIntValue(valueStr)
			case 27:
				currentRecord.GoodsName = parseValue(valueStr)
			case 28:
				currentRecord.CostChargeType = parseIntValue(valueStr)
			case 29:
				currentRecord.DeliveryChargeType = parseIntValue(valueStr)
			case 30:
				currentRecord.MainStoneName = parseValue(valueStr)
			case 31:
				currentRecord.MainStoneTypeID = parseIntValue(valueStr)
			case 32:
				currentRecord.MainStoneTypeName = parseValue(valueStr)
			case 33:
				currentRecord.MainStoneWeight = parseValue(valueStr)
			case 34:
				currentRecord.ViceStoneTypeName = parseValue(valueStr)
			case 35:
				currentRecord.ViceStoneWeight = parseValue(valueStr)
			case 36:
				currentRecord.CerNumber = parseValue(valueStr)
			case 37:
				currentRecord.RecyleGoldPrice = parseIntValue(valueStr)
			case 38:
				currentRecord.OriginalGoldPrice = parseIntValue(valueStr)
			case 39:
				currentRecord.StockCode = parseValue(valueStr)
			case 40:
				currentRecord.Weight = parseIntValue(valueStr)
			case 41:
				currentRecord.GoldWeight = parseIntValue(valueStr)
			case 42:
				currentRecord.GoldWeightRecycle = parseIntValue(valueStr)
			case 43:
				currentRecord.GoldWeightLossRate = parseValue(valueStr)
			case 44:
				currentRecord.RecycleLabourPrice = parseIntValue(valueStr)
			case 45:
				currentRecord.RecycleLabourAmount = parseIntValue(valueStr)
			case 46:
				currentRecord.GoodsNum = parseIntValue(valueStr)
			case 47:
				currentRecord.GoodsAmount = parseIntValue(valueStr)
			case 48:
				currentRecord.OriginalAmount = parseIntValue(valueStr)
			case 49:
				currentRecord.OldOriginalLabelPrice = parseIntValue(valueStr)
			case 50:
				currentRecord.RecycleWeightDiff = parseIntValue(valueStr)
			case 51:
				currentRecord.NewOldDiffDay = parseIntValue(valueStr)
			case 52:
				currentRecord.OldOriginalLabelWeight = parseIntValue(valueStr)
			case 53:
				currentRecord.OldOriginalCostPrice = parseIntValue(valueStr)
			case 54:
				currentRecord.RecycleAmount = parseInt64Value(valueStr)
			case 55:
				currentRecord.DeductAmount = parseIntValue(valueStr)
			case 56:
				currentRecord.CostAmount = parseIntValue(valueStr)
			case 57:
				currentRecord.CutNewCostAmount = parseIntValue(valueStr)
			case 58:
				currentRecord.OldProfitAmount = parseIntValue(valueStr)
			case 59:
				currentRecord.OldValuedAmount = parseInt64Value(valueStr)
			case 60:
				currentRecord.GoldType = parseIntValue(valueStr)
			case 61:
				currentRecord.ShiftStatus = parseIntValue(valueStr)
			case 62:
				currentRecord.Remarks = parseValue(valueStr)
			case 63:
				currentRecord.Status = parseIntValue(valueStr)
			case 64:
				currentRecord.CreateTime = parseIntValue(valueStr)
			case 65:
				currentRecord.UpdateTime = parseIntValue(valueStr)
			case 66:
				currentRecord.DeleteTime = parseIntValue(valueStr)
			case 67:
				currentRecord.AcountTime = parseIntValue(valueStr)
			case 68:
				currentRecord.SaleDatetime = parseValue(valueStr)
			case 69:
				currentRecord.IsOverExchange = parseIntValue(valueStr)
			case 70:
				currentRecord.PicURL = parseValue(valueStr)
			case 71:
				currentRecord.PicURLList = parseValue(valueStr)
			case 72:
				currentRecord.OldStockTypeID = parseIntValue(valueStr)
			case 73:
				currentRecord.OldStockTypeName = parseValue(valueStr)
			case 74:
				currentRecord.FinancialCutNewCostAmount = parseIntValue(valueStr)
			case 75:
				currentRecord.FinancialCostAmount = parseIntValue(valueStr)
			case 76:
				currentRecord.FinancialOldProfitAmount = parseIntValue(valueStr)
			case 77:
				currentRecord.MarketCutNewCostAmount = parseIntValue(valueStr)
			case 78:
				currentRecord.MarketCostAmount = parseIntValue(valueStr)
			case 79:
				currentRecord.MarketOldProfitAmount = parseIntValue(valueStr)
			case 80:
				currentRecord.SettleMerchantShopID = parseIntValue(valueStr)
			case 81:
				currentRecord.InlaidJewelryAmount = parseIntValue(valueStr)
			}

		} else {
			fmt.Println(line)
		}

	}

	return currentRecord
}

// parseValue 解析字符串值为对应的类型
func parseValue(value string) *string {
	if value == "NULL" {
		return nil
	}
	trimmed := strings.Trim(value, "'")
	return &trimmed
}

// parseIntValue 解析字符串值为int类型
func parseIntValue(value string) *int {
	if value == "NULL" {
		return nil
	}
	if i, err := strconv.Atoi(value); err == nil {
		return &i
	}
	return nil
}

// parseInt64Value 解析字符串值为int64类型
func parseInt64Value(value string) *int64 {
	if value == "NULL" {
		return nil
	}
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return &i
	}
	return nil
}

// parseDecimalValue 解析字符串值为decimal类型（以字符串形式存储）
func parseDecimalValue(value string) *string {
	if value == "NULL" {
		return nil
	}
	trimmed := strings.Trim(value, "'")
	return &trimmed
}
