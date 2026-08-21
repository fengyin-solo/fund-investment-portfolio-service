# 基金定投服务

纯 Go 标准库实现的基金定投后端服务，金额单位一律为 **分（int64）**，避免浮点精度问题。

## 运行方式

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## 实体清单

| 实体 | 说明 |
|------|------|
| Fund | 基金（代码、名称、类型、风险等级、净值） |
| Account | 账户（持有人、余额，单位为分） |
| Plan | 定投计划（金额、周期、状态、下次执行日） |
| Transaction | 交易（定投买入 / 赎回，份额、单价、状态） |
| Holding | 持仓（总份额、成本、市值） |
| NavHistory | 净值历史 |
| Dividend | 分红（现金 / 再投资） |

## 状态机

- **Plan**: active → paused → active；active → closed；paused → closed
- **Transaction**: pending → confirmed / failed

## API 列表

### Fund
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/funds | 创建基金 |
| GET | /api/funds | 基金列表（支持 type/status/keyword 筛选） |
| GET | /api/funds/{id} | 查看基金 |
| PUT | /api/funds/{id} | 更新基金 |
| DELETE | /api/funds/{id} | 删除基金 |

### Account
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/accounts | 创建账户 |
| GET | /api/accounts | 账户列表（支持 keyword 筛选） |
| GET | /api/accounts/{id} | 查看账户 |
| PUT | /api/accounts/{id} | 更新账户 |
| DELETE | /api/accounts/{id} | 删除账户 |

### Plan
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/plans | 创建定投计划 |
| GET | /api/plans | 计划列表（支持 status/interval/account_id 筛选） |
| GET | /api/plans/{id} | 查看计划 |
| PUT | /api/plans/{id} | 更新计划 |
| POST | /api/plans/{id}/status | 变更计划状态 |
| POST | /api/plans/{id}/execute | 执行定投（扣款、生成交易、更新持仓） |
| DELETE | /api/plans/{id} | 删除计划 |

### Transaction
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/transactions | 创建交易记录 |
| GET | /api/transactions | 交易列表（支持 account_id/fund_id/type/status/plan_id/时间范围 筛选） |
| GET | /api/transactions/{id} | 查看交易 |
| PUT | /api/transactions/{id} | 更新交易 |
| POST | /api/transactions/redeem | 赎回（校验持仓份额、更新余额） |
| DELETE | /api/transactions/{id} | 删除交易 |

### Holding
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/holdings | 创建持仓 |
| GET | /api/holdings | 持仓列表（支持 account_id/fund_id 筛选） |
| GET | /api/holdings/{id} | 查看持仓 |
| PUT | /api/holdings/{id} | 更新持仓 |
| DELETE | /api/holdings/{id} | 删除持仓 |

### NavHistory
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/nav-history | 创建净值历史 |
| GET | /api/nav-history | 净值历史列表（支持 fund_id 筛选） |
| GET | /api/nav-history/{id} | 查看净值历史 |
| PUT | /api/nav-history/{id} | 更新净值历史 |
| DELETE | /api/nav-history/{id} | 删除净值历史 |

### Dividend
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/dividends | 创建分红 |
| GET | /api/dividends | 分红列表（支持 account_id/fund_id/type 筛选） |
| GET | /api/dividends/{id} | 查看分红 |
| PUT | /api/dividends/{id} | 更新分红 |
| DELETE | /api/dividends/{id} | 删除分红 |

### Stats（聚合查询）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/holdings/{account_id} | 账户持仓汇总（按基金分组，含市值/成本/浮动盈亏） |
| GET | /api/stats/profit/{account_id} | 账户收益统计 |
| GET | /api/stats/fund-types | 基金按类型统计 |
| GET | /api/stats/fund-risks | 基金按风险等级统计 |
| GET | /api/stats/top-funds?n=10 | TOP N 热门基金（按累计定投金额） |

## 统一响应格式

```json
{"code":0,"message":"ok","data":...}
```

错误码映射：
- 400: ValidationError
- 404: ErrNotFound
- 409: ErrConflict
- 500: 其他内部错误

## 测试

```bash
go test ./...
```

测试覆盖：
- Store CRUD + 冲突/不存在断言
- Service 定投执行、余额不足 failed、赎回、计划状态机、聚合查询、跨实体校验、分页、时间筛选
