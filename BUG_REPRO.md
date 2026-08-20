# Bug

结算重试使用不同副作用键，成功状态又被首轮延迟快照覆盖。

# 触发

让首次结算返回临时失败、第二次成功，再把首次 running 快照延迟写回。

# 错误信息

`settlement side effect ran 2 times`

`delayed callback rolled final state back: completed=model.Job{ID:"plan-7", State:"succeeded", Version:1} cached=model.Job{ID:"plan-7", State:"running", Version:1}`
