# Bug

事件首次发布临时失败时重做整笔结算并更换事件键，造成重复写入、重复事件和多个审计状态。

# 触发

运行 order-9 结算，让事件总线第一次返回临时错误、第二次成功，再读取三个出口。

# 错误信息

`transaction was repeated or duplicated: begins=2 writes=2`

`published event was duplicated: calls=2 delivered=2`

`audit retained intermediate states: []string{"pending", "success"}`
