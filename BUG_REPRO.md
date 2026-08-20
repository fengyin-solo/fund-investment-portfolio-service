# Bug

业务拒绝与临时失败被当成同类错误，导致拒绝重复请求并留下记录，临时失败恢复后仍返回失败和重复记录。

# 触发

先提交 reject 模式投资，再用新的处理器提交 temporary 模式投资。

# 错误信息

`rejection retried or left an investment record: calls=3 entries=1`

`temporary failure still surfaced after recovery: investment unavailable: pricing service busy`
