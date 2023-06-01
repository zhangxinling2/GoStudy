--解锁要解自己的锁，所以要进行判断
if redis.call('GET',KEYS[1]) == ARGV[1] then
    return redis.call('DEL',KEYS[1])
else
    return 0
end