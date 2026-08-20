# Bug

源端在第二项失败后，批处理没有传回源错误，而是等待截止时间并返回第一项的部分结果。

# 触发

导入 fund-a、fund-b、fund-c，并把固定失败位置设为第二项，使用 50ms 截止时间等待结果。

# 错误信息

`source failure was masked: context deadline exceeded`

`failed batch exposed partial results: []string{"fund-a"}`
