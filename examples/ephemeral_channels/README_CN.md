# 临时通道 (Ephemeral Channels) 示例

## 背景

在有状态的应用中，并非所有数据都需要在整个对话历史中持久保存。有些数据是“临时的”——它仅对紧接的下一步或特定的超步（并行执行块）有效，之后应被丢弃。例如临时的搜索结果、中间推理步骤，或触发立即动作但不应混淆未来轮次的标志。**临时通道 (Ephemeral Channels)** 提供了一种自动管理这种生命周期的方法。

## 功能特性

*   **自动清理**: 标记为临时的通道在步骤完成后会自动从状态中清除。
*   **作用域隔离**: 防止临时数据泄漏到未来的执行步骤中，减少上下文污染。
*   **可配置**: 通过带有 `isEphemeral` 标志的 `RegisterChannel` 进行定义。

## 实现原理

该功能通过 `graph/schema.go` 中的 `CleaningStateSchema` 接口实现。
1.  `MapSchema` 维护一组 `EphemeralKeys`。
2.  在图执行循环 (`InvokeWithConfig`) 中，当当前步骤的所有节点执行完毕且结果合并后：
3.  调用 `Cleanup(state)` 方法。
4.  该方法返回一个新的状态映射，其中所有在 `EphemeralKeys` 中的键都被移除。

## 代码导读

在 `main.go` 中：

1.  **Schema 定义**:
    ```go
    schema := graph.NewMapSchema()
    // 注册 "temp_data" 为临时通道
    schema.RegisterChannel("temp_data", graph.OverwriteReducer, true)
    ```
    我们显式告诉 Schema `temp_data` 应被视为临时数据。

2.  **Producer 节点**:
    设置 `temp_data` 为 "secret_code_123"。

3.  **Consumer 节点**:
    检查 `temp_data`。由于 Producer 在步骤 1 运行，Consumer 在步骤 2 运行（顺序执行），清理逻辑会在它们之间触发。
    *注意：在 LangGraph 语义中，临时值在步骤*之后*被清除。如果 Producer -> Consumer 是直接转换，它们被视为不同的步骤。*

4.  **验证**:
    输出显示 Consumer（在下一步运行）**没有**看到 `temp_data`，确认它已被清理。

## 如何运行

```bash
go run main.go
```
