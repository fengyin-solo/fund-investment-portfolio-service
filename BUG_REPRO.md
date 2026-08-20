# Bug

批处理在循环中延后关闭句柄，先触发资源上限；登记完成又吞掉业务错误，审计只记录开始状态。

# 触发

按 fund-a、fund-b、fund-c、blocked-fund 的顺序处理，句柄上限设为 2。

# 错误信息

`business failure was replaced: statement handle limit reached`

`batch stopped before valid items committed: []string{"fund-a", "fund-b"}`

`audit and committed items diverged: []string{"started:fund-a", "started:fund-b"}`
