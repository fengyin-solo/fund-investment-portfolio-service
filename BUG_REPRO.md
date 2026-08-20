# Bug

行情中断后的组合构建返回错误，但半初始化组合仍进入缓存，后续查询可见只含债券的残留对象。

# 触发

先为 client-a 发起会中断的组合构建，收到错误后查询 client-a；再为 client-b 正常构建并查询。

# 错误信息

`failed preparation leaked a partial portfolio: &builder.Portfolio{ID:"client-a", Positions:map[string]int64{"bond":6000}, Ready:false}`
