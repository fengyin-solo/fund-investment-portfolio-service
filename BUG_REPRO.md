# Bug

已取消的报价请求继续重试，下一个正常请求又复用旧取消状态。

# 触发

先提交一个已取消的 F001 报价请求，随后使用新的背景上下文再次查询 F001。

# 错误信息

`canceled quote reached pricing 3 times`

`fresh request inherited the earlier cancellation: value=0 err=context canceled`
