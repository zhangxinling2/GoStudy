--lock时先判断
local val = redis.call("GET",KEYS[1])
--如果没有这个key，那么就可以加锁
if val == false then
    return redis.call("SET",KEYS[1],ARGV[1],"EX",ARGV[2])
--如果已经持有了这个key，刷新
else if val == ARGV[1] then
    redis.call("expire",KEY[1],ARGV[2])
    return "Ok"
--别人持有锁
else
    return ""
end