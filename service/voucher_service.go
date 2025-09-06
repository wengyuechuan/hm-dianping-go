package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
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