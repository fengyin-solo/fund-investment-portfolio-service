# Bug

请求对象归还复用池后仍被异步审计持有，下一租户改写对象造成身份与标签混合。

# 触发

先处理 tenant-a 并延后审计，再处理 tenant-b，最后用同步信号同时放行两条审计。

# 错误信息

`async audit mixed request identities: []audit.Entry{audit.Entry{Tenant:"tenant-b", Tag:"audit:tenant-a"}, audit.Entry{Tenant:"tenant-b", Tag:"audit:tenant-a"}}`
