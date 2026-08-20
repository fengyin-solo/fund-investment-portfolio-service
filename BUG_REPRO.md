# Bug

解析缓冲区在两批导入间复用，第一批缓存和待导出内容被第二批覆盖。

# 触发

先接收 batch-a/ALPHA，再接收 batch-b/BRAVO，随后读取 batch-a 并刷新导出队列。

# 错误信息

`cached first batch changed after second import: parser.Record{ID:"batch-a", Payload:[]uint8{0x42, 0x52, 0x41, 0x56, 0x4f}}`

`export mixed the two import payloads: []string{"BRAVO", "BRAVO"}`
