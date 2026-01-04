package mirg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/persistence/cache"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/sliceconv"
	"math"
	"os"
	"time"
)

type BonusModel struct {
	Telephone  string  `json:"telephone"`
	Bonus      float64 `json:"bonus"`
	TotalBonus float64 `json:"total_bonus"`
}

var BonusHeader = []excel.Header{
	{Key: "telephone", Name: "手机号", Type: "string"},
	{Key: "bonus", Name: "积分", Type: "float"},
	{Key: "total_bonus", Name: "累计积分", Type: "float"},
}

func loadExcel(path, sheet string, isAdd, isChange bool) ([]*BonusModel, error) {
	var data []*BonusModel
	r := excel.NewRead(path)
	defer r.Close()
	err := r.Read(sheet, BonusHeader, &data)

	if err != nil {
		return nil, err
	}

	for _, d := range data {
		if isChange {
			d.Bonus = -d.Bonus
		}

		d.Bonus = roundToNDecimal(d.Bonus, 2)
		if isAdd {
			d.TotalBonus = d.Bonus
		}
	}

	return data, nil
}

func roundToNDecimal(f float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(f*pow) / pow
}

// ====================================================================

type Customer struct {
	Id              int32  `json:"id"`
	Name            string `json:"name"`
	MerchantId      int32  `json:"merchant_id"`
	BrandId         int32  `json:"brand_id"`
	MerchantShopId  int32  `json:"merchant_shop_id"`
	Telephone       string `json:"telephone"`
	SuffixTelephone string `json:"suffix_telephone"`
	CreateTime      int32  `json:"create_time"`
	UpdateTime      int32  `json:"update_time"`
	Status          int32  `json:"status"`
	BirthdayDt      string `json:"birthday_dt"`
	MarryDate       string `json:"marry_date"`
	LastBuyTime     string `json:"last_buy_time"`
	AddDatetime     string `json:"add_datetime"`
	OriginChannel   int32  `json:"origin_channel"`
	OriginPlatform  int32  `json:"origin_platform"`
	Level           int32  `json:"level"`
}

type CustomerBpLog struct {
	Id             int32 `json:"id"`
	MerchantId     int32 `json:"merchant_id"`
	BrandId        int32 `json:"brand_id"`
	MerchantShopId int32 `json:"merchant_shop_id"`

	CustomerId   int32  `json:"customer_id"`
	CustomerName string `json:"customer_name"`

	PointTotal   float64 `json:"point_total"`
	Point        float64 `json:"point"`
	Type         int32   `json:"type"`
	Comment      string  `json:"comment"`
	CreateTime   int32   `json:"create_time"`
	ActivityName string  `json:"activity_name"`
	ActivityId   int32   `json:"activity_id"`
	AdminUserId  int32   `json:"admin_user_id"`

	OrderId      int32  `json:"order_id"`
	OrderSn      string `json:"order_sn"`
	OrderGoodsId int32  `json:"order_goods_id"`
	StockCode    string `json:"stock_code"`
}

//====================================================================

func CustomerBpAction(accessToken int) {
	if accessToken != time.Now().Day() {
		panic(errors.New("check accessToken"))
	}

	merchantId := 0
	BrandId := 0
	Path := "/Users/linlong/Desktop/a.xlsx"
	sheet := "Sheet1"
	isAdd := true     // 新增 or 扣减
	isChange := false // excel中数据是否需要加上负数

	// 获取表格信息
	data, err := loadExcel(Path, sheet, isAdd, isChange)
	if err != nil {
		panic(err)
	}
	tels := sliceconv.Extract(data, func(d *BonusModel) string {
		return d.Telephone
	})
	telHash := sliceconv.Map(data, func(d *BonusModel) (key string, value *BonusModel) {
		return d.Telephone, d
	})

	// 获取会员信息
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad("./config/online/db.yaml", &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}
	customers := make([]*Customer, 0, len(tels))
	for _, chuck := range sliceconv.Chunk(tels, 10000) {
		chuckCustomers := make([]*Customer, 0, 10000)

		//query := orm.NewBoolQuery().Must(
		//	orm.Eq("merchant_id", merchantId),
		//	orm.Eq("brand_id", BrandId),
		//	orm.Eq("status", 1),
		//	orm.In("telephone", chuck),
		//)
		//if query.IsEmpty() {
		//	panic(errors.New("query is empty"))
		//}
		//sql, args := query.Do()
		//if err := db.Where(sql, args...).Find(&chuckCustomers).Error; err != nil {
		//	panic(err)
		//}

		if err := db.Where("merchant_id = ?", merchantId).
			Where("brand_id = ?", BrandId).
			Where("status = ?", 1).
			Where("telephone in (?)", chuck).
			Find(&chuckCustomers).Error; err != nil {
			panic(err)
		}
		customers = append(customers, chuckCustomers...)
	}

	// 转换成新的结构体
	var Type int32 = 2
	Comment := "手工录入(积分扣减)"
	if isAdd {
		Comment = "手工录入(积分增加)"
	}
	customerBpLog := sliceconv.Change(customers, func(customer *Customer) *CustomerBpLog {
		return &CustomerBpLog{
			MerchantId:     customer.MerchantId,
			BrandId:        customer.BrandId,
			MerchantShopId: customer.MerchantShopId,
			CustomerId:     customer.Id,
			CustomerName:   customer.Name,
			Type:           Type,
			Comment:        Comment,
			PointTotal:     telHash[customer.Telephone].TotalBonus,
			Point:          telHash[customer.Telephone].Bonus,
		}
	})

	// 发送消息
	ctx := context.Background()
	var redisConfig cache.Config
	configurator.NewYaml().MustLoad("./config/online/redis.yaml", &redisConfig)
	mq := cache.NewRedis(&redisConfig)
	for _, item := range customerBpLog {
		b, err := json.Marshal(&item)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(mq.LPush(ctx, "mq_pipeline_bonus", string(b)))
	}
}

func CustomerJson() {
	b, err := os.ReadFile("/Users/linlong/Desktop/insert_data_fixed.json")
	if err != nil {
		return
	}

	var customerBpLog []*CustomerBpLog
	err = json.Unmarshal(b, &customerBpLog)
	if err != nil {
		return
	}

	// 发送消息
	ctx := context.Background()
	var redisConfig cache.Config
	configurator.NewYaml().MustLoad("./config/online/redis.yaml", &redisConfig)
	mq := cache.NewRedis(&redisConfig)
	for _, item := range customerBpLog {
		b, err := json.Marshal(&item)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(mq.LPush(ctx, "mq_pipeline_bonus", string(b)))
	}
}

func CustomerBpNum(accessToken int) {
	if accessToken != time.Now().Day() {
		panic(errors.New("check accessToken"))
	}

	merchantId := 0
	isAdd := true // 新增 or 扣减
	//isChange := false // excel中数据是否需要加上负数

	// 获取会员信息
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad("./config/online/db_read.yaml", &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}
	customers := make([]*Customer, 0, 100000)
	if err := db.Where("merchant_id = ?", merchantId).
		Where("status = ?", 1).
		Find(&customers).Error; err != nil {
		panic(err)
	}

	// 转换成新的结构体
	var Type int32 = 2
	Comment := "手工录入(积分扣减)"
	if isAdd {
		Comment = "手工录入(积分增加)"
	}
	customerBpLog := sliceconv.Change(customers, func(customer *Customer) *CustomerBpLog {
		return &CustomerBpLog{
			MerchantId:     customer.MerchantId,
			BrandId:        customer.BrandId,
			MerchantShopId: customer.MerchantShopId,
			CustomerId:     customer.Id,
			CustomerName:   customer.Name,
			Type:           Type,
			Comment:        Comment,
			PointTotal:     188,
			Point:          188,
		}
	})

	//return

	// 发送消息
	ctx := context.Background()
	var redisConfig cache.Config
	configurator.NewYaml().MustLoad("./config/online/redis.yaml", &redisConfig)
	mq := cache.NewRedis(&redisConfig)
	for index, item := range customerBpLog {
		b, err := json.Marshal(&item)
		if err != nil {
			fmt.Println(err)
			continue
		}
		res, err := mq.LPush(ctx, "mq_pipeline_bonus", string(b))
		fmt.Println(index, res, err)
	}
}
