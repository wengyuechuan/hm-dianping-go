package service

import (
	"context"
	"fmt"
	"hm-dianping-go/dao"
	"hm-dianping-go/utils"
	"time"
)

// GetDailyUV 获取指定日期的UV统计
func GetDailyUV(ctx context.Context, date string) *utils.Result {
	// 验证日期格式
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return utils.ErrorResult("日期格式错误，请使用YYYY-MM-DD格式")
	}

	uvKey := fmt.Sprintf("uv:daily:%s", date)

	// 使用HyperLogLog获取UV数量
	count, err := dao.Redis.PFCount(ctx, uvKey).Result()
	if err != nil {
		return utils.ErrorResult("获取UV统计失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"date": date,
		"uv":   count,
	})
}

// GetTodayUV 获取今日UV统计
func GetTodayUV(ctx context.Context) *utils.Result {
	today := time.Now().Format("2006-01-02")
	return GetDailyUV(ctx, today)
}

// GetUVRange 获取指定日期范围的UV统计
func GetUVRange(ctx context.Context, startDate, endDate string) *utils.Result {
	// 验证日期格式
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return utils.ErrorResult("开始日期格式错误，请使用YYYY-MM-DD格式")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return utils.ErrorResult("结束日期格式错误，请使用YYYY-MM-DD格式")
	}

	if start.After(end) {
		return utils.ErrorResult("开始日期不能晚于结束日期")
	}

	// 限制查询范围，避免查询过多数据
	if end.Sub(start).Hours() > 24*30 { // 最多30天
		return utils.ErrorResult("查询范围不能超过30天")
	}

	var results []map[string]interface{}

	// 遍历日期范围
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		uvKey := fmt.Sprintf("uv:daily:%s", dateStr)

		count, err := dao.Redis.PFCount(ctx, uvKey).Result()
		if err != nil {
			// 如果某天的数据获取失败，记录为0
			count = 0
		}

		results = append(results, map[string]interface{}{
			"date": dateStr,
			"uv":   count,
		})
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"startDate": startDate,
		"endDate":   endDate,
		"data":      results,
	})
}

// GetRecentUV 获取最近N天的UV统计
func GetRecentUV(ctx context.Context, days int) *utils.Result {
	if days <= 0 || days > 30 {
		return utils.ErrorResult("天数必须在1-30之间")
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -(days - 1))

	return GetUVRange(ctx, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}

// GetUVSummary 获取UV统计摘要（今日、昨日、本周、本月）
func GetUVSummary(ctx context.Context) *utils.Result {
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	// 获取今日UV
	todayResult := GetDailyUV(ctx, today)
	var todayUV int64 = 0
	if todayResult.Success {
		if data, ok := todayResult.Data.(map[string]interface{}); ok {
			if uv, ok := data["uv"].(int64); ok {
				todayUV = uv
			}
		}
	}

	// 获取昨日UV
	yesterdayResult := GetDailyUV(ctx, yesterday)
	var yesterdayUV int64 = 0
	if yesterdayResult.Success {
		if data, ok := yesterdayResult.Data.(map[string]interface{}); ok {
			if uv, ok := data["uv"].(int64); ok {
				yesterdayUV = uv
			}
		}
	}

	// 获取本周UV（周一到今天）
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	if now.Weekday() == time.Sunday {
		weekStart = now.AddDate(0, 0, -6)
	}
	weekResult := GetUVRange(ctx, weekStart.Format("2006-01-02"), today)
	var weekUV int64 = 0
	if weekResult.Success {
		if data, ok := weekResult.Data.(map[string]interface{}); ok {
			if dataList, ok := data["data"].([]map[string]interface{}); ok {
				for _, item := range dataList {
					if uv, ok := item["uv"].(int64); ok {
						weekUV += uv
					}
				}
			}
		}
	}

	// 获取本月UV（月初到今天）
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthResult := GetUVRange(ctx, monthStart.Format("2006-01-02"), today)
	var monthUV int64 = 0
	if monthResult.Success {
		if data, ok := monthResult.Data.(map[string]interface{}); ok {
			if dataList, ok := data["data"].([]map[string]interface{}); ok {
				for _, item := range dataList {
					if uv, ok := item["uv"].(int64); ok {
						monthUV += uv
					}
				}
			}
		}
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"today":     todayUV,
		"yesterday": yesterdayUV,
		"thisWeek":  weekUV,
		"thisMonth": monthUV,
	})
}
