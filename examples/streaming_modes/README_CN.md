# 流式模式 (Streaming Modes) 示例

## 背景

对于长时间运行的 LLM 链或复杂的 Agent 工作流，等待最终结果通常是不可接受的。用户期望获得实时反馈。LangGraph Go 支持强大的流式架构，允许你订阅发生的事件。不同的用例需要不同的流“视图”——有时你需要每个 Token，有时只需要每个节点的最终输出。**流式模式 (Streaming Modes)** 允许你配置这种粒度。

## 功能特性

*   **`StreamModeUpdates`**: 在每个节点完成时发射其输出。适用于显示进度（例如，“步骤 1 完成”，“工具已执行”）。
*   **`StreamModeValues`**: 在每一步后发射完整的图状态。适用于调试或渲染整个上下文的 UI。
*   **`StreamModeMessages`**: (计划中) 发射 LLM Token 以实现打字机效果。
*   **`StreamModeDebug`**: 发射所有内部事件以进行深度检查。

## 实现原理

流式逻辑由 `graph/streaming.go` 中的 `StreamingRunnable` 处理。
1.  它将一个 `StreamingListener`（实现了 `GraphCallbackHandler` 和 `NodeListener`）注入到执行上下文中。
2.  随着图的执行，节点和图引擎会发射事件（`NodeEventComplete`, `OnGraphStep` 等）。
3.  `StreamingListener.shouldEmit` 方法根据配置的 `StreamMode` 过滤这些事件。
4.  允许的事件被发送到返回给调用者的 Go 通道中。

## 代码导读

在 `main.go` 中：

1.  **配置**:
    ```go
    g.SetStreamConfig(graph.StreamConfig{
        Mode: graph.StreamModeUpdates,
    })
    ```
    我们将图配置为 `updates` 模式进行流式传输。

2.  **执行**:
    ```go
    streamResult := runnable.Stream(context.Background(), nil)
    for event := range streamResult.Events { ... }
    ```
    我们遍历 `Events` 通道。

3.  **输出**:
    你会看到事件随着 `step_1` 和 `step_2` 的完成实时打印，而不是等待整个图结束。

## 如何运行

```bash
go run main.go
```
