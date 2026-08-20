# Bug

报价请求取消后仍启动下一次重试，错误退出的任务没有完成等待计数，关闭过程超时。

# 触发

用许可通道放行首次失败，让第二次调用进入后取消请求；先等待关闭，再放行在途调用。

# 错误信息

`shutdown waited for canceled quote work: context deadline exceeded`

`quote retry 3 started after cancellation`
