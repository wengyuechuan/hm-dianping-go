package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
	"time"
)

// GetVoucherList 获取优惠券列表
func GetVoucherList(shopId uint) *utils.Result {
	var vouchers []models.Voucher
	err := dao.DB.Where("shop_id = ? AND status = 1", shopId).Find(&vouchers).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(vouchers)
}

// AddSeckillVoucherRequest 添加秒杀券请求结构
type AddSeckillVoucherRequest struct {
	ShopID      uint      `json:"shopId" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	SubTitle    string    `json:"subTitle"`
	Rules       string    `json:"rules"`
	PayValue    int64     `json:"payValue" binding:"required"`
	ActualValue int64     `json:"actualValue" binding:"required"`
	Stock       int       `json:"stock" binding:"required,min=1"`
	BeginTime   time.Time `json:"beginTime" binding:"required"`
	EndTime     time.Time `json:"endTime" binding:"required"`
}

// AddSeckillVoucher 添加秒杀券
func AddSeckillVoucher(req *AddSeckillVoucherRequest) *utils.Result {
	// 验证时间逻辑
	if req.EndTime.Before(req.BeginTime) {
		return utils.ErrorResult("结束时间不能早于开始时间")
	}
	
	if req.BeginTime.Before(time.Now()) {
		return utils.ErrorResult("开始时间不能早于当前时间")
	}

	// 验证价格逻辑
	if req.PayValue <= req.ActualValue {
		return utils.ErrorResult("支付金额必须大于实际价值")
	}

	// 开启事务
	tx := dao.DB.Begin()
	if tx.Error != nil {
		return utils.ErrorResult("事务开启失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建普通优惠券
	voucher := &models.Voucher{
		ShopID:      req.ShopID,
		Title:       req.Title,
		SubTitle:    req.SubTitle,
		Rules:       req.Rules,
		PayValue:    req.PayValue,
		ActualValue: req.ActualValue,
		Type:        1, // 1-秒杀券
		Status:      1, // 1-上架
		Stock:       req.Stock,
		BeginTime:   &req.BeginTime,
		EndTime:     &req.EndTime,
	}

	if err := tx.Create(voucher).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResult("创建优惠券失败")
	}

	// 2. 创建秒杀券记录
	seckillVoucher := &models.SeckillVoucher{
		VoucherID:  voucher.ID,
		Stock:      req.Stock,
		CreateTime: time.Now(),
		BeginTime:  req.BeginTime,
		EndTime:    req.EndTime,
		UpdateTime: time.Now(),
	}

	if err := tx.Create(seckillVoucher).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResult("创建秒杀券失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return utils.ErrorResult("事务提交失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"voucherId": voucher.ID,
		"message":   "秒杀券创建成功",
	})
}

// GetSeckillVoucher 获取秒杀券详情
func GetSeckillVoucher(voucherID uint) *utils.Result {
	// 获取优惠券基本信息
	var voucher models.Voucher
	if err := dao.DB.First(&voucher, voucherID).Error; err != nil {
		return utils.ErrorResult("优惠券不存在")
	}

	// 检查是否为秒杀券
	if voucher.Type != 1 {
		return utils.ErrorResult("该优惠券不是秒杀券")
	}

	// 获取秒杀券详细信息
	seckillVoucher, err := dao.GetSeckillVoucherByID(voucherID)
	if err != nil {
		return utils.ErrorResult("秒杀券信息获取失败")
	}

	// 组合返回数据
	result := map[string]interface{}{
		"voucher":        voucher,
		"seckillVoucher": seckillVoucher,
	}

	return utils.SuccessResultWithData(result)
}