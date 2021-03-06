-- 1. 传入参数
local uid = tostring(ARGV[1])

local iphoneStockKey = tostring(KEYS[1])
local luckyGuysKey = tostring(KEYS[2])

-- local luckyGuysKey = "promote:iphone:luckyguys"
-- local iphoneStockKey = "promote:iphone:stock"


-- 2. 检查 uid 是否已经中奖用户
local exists = redis.call("SISMEMBER", luckyGuysKey, uid)
--- demo: https://segmentfault.com/a/1190000018070172
---  这里需要转换成数字进行判断， 貌似 redis-lua 没有 toboolean() 这种方法
-- if not exists then
--     return -2
-- end
if tonumber(exists) == 1  then 
    -- 2.1 已经中奖
    return -2
end

-- 3. 检查库存数量
local remain = redis.call("GET",iphoneStockKey)
local n = tonumber(remain)

--- 3.1  return -1 : 没开始
if not n then
    return -1
end
--- 3.2  return 0 : 结束了

if n == 0 then
    return 0
end

-- 4. 修改库存
--- 4.2 将用户加入到中奖列表
redis.call("SADD",luckyGuysKey,uid)
--- 4.1 降低库存
redis.call("DECR",iphoneStockKey)



