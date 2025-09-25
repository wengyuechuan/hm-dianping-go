-- 1. 参数列表
-- 1.1 优惠券id
local voucherId = ARGV[1]
-- 1.2 用户id
local userId = ARGV[2]

-- 2. 数据key
-- 2.1 库存key
local stockKey = "cache:seckill_voucher:stock:" .. voucherId
-- 2.2 订单key
local orderKey = "cache:seckill_voucher:order:" .. voucherId


-- 3. 脚本业务
-- 3.1 判断库存是否充足
if(tonumber(redis.call('get', stockKey)) <= 0) then
    return 1
end
-- 3.2 判断用户是否重复下单
if(redis.call('sismember', orderKey, userId) == 1) then
    return 2
end

-- 3.3 扣减库存
redis.call('incrby', stockKey, -1)
-- 3.4 下单（保存用户）
redis.call('sadd', orderKey, userId)
return 0