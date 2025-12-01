# 人在回路 (Human-in-the-loop) 示例

本示例演示如何使用 LangGraphGo 实现 **人在回路 (Human-in-the-loop, HITL)** 工作流。

## 1. 背景

在许多实际的 Agent 应用中，完全自动化并不总是理想或可能的。你可能需要人工介入来：
- **批准** 关键操作（例如：部署代码、发送邮件）。
- **提供输入**，即 Agent 无法获取的信息。
- **修正** Agent 的推理或状态，然后再继续执行。

LangGraphGo 通过 **中断 (Interrupt)** 机制支持这一功能。你可以配置图在特定节点之前或之后暂停执行，允许外部系统（或人类）在恢复执行前检查并修改状态。

## 2. 核心概念

- **InterruptBefore**: 一个配置选项，指示图在进入指定节点 *之前* 停止执行。
- **GraphInterrupt**: 当图暂停时返回的一种特定错误类型。它包含当前状态和停止位置的节点信息。
- **ResumeFrom**: 一个配置选项，用于从特定节点重新开始执行，通常是之前被中断的那个节点。

## 3. 工作原理

1.  **定义图**: 标准的图定义，包含节点和边。
2.  **初始运行**: 调用图，并将 `InterruptBefore` 设置为目标节点（例如 "human_approval"）。
3.  **暂停与检查**: 图执行到目标节点前会停止，并返回 `GraphInterrupt` 错误。当前状态被保留。
4.  **人工交互**: 应用程序捕获中断，将状态展示给人类（此处为模拟），并根据输入更新状态（例如设置 `Approved = true`）。
5.  **恢复执行**: 使用 *更新后的状态* 再次调用图，并将 `ResumeFrom` 设置为被中断的节点。图将从该点继续执行。

## 4. 代码亮点

### 设置中断
```go
config := &graph.Config{
    InterruptBefore: []string{"human_approval"},
}
// 当执行到 "human_approval" 时，Invoke 调用将返回 GraphInterrupt 错误
res, err := runnable.InvokeWithConfig(ctx, initialState, config)
```

### 处理中断
```go
var interrupt *graph.GraphInterrupt
if errors.As(err, &interrupt) {
    // 获取中断时刻的状态
    currentState := interrupt.State.(State)
    // ... 展示给人类 ...
}
```

### 恢复执行
```go
// 使用人工输入更新状态
currentState.Approved = true 

// 恢复配置
resumeConfig := &graph.Config{
    ResumeFrom: []string{"human_approval"},
}
// 使用修改后的状态继续执行
finalRes, err := runnable.InvokeWithConfig(ctx, currentState, resumeConfig)
```

## 5. 运行示例

```bash
go run main.go
```

**预期输出:**
```text
=== Starting Workflow (Phase 1) ===
[Process] Processing request: Deploy to Production
Workflow interrupted at node: human_approval
Current State: {Input:Deploy to Production Approved:false Output:Processed: Deploy to Production}

=== Human Interaction ===
Reviewing request...
Approving request...

=== Resuming Workflow (Phase 2) ===
[Human] Request APPROVED.
[Finalize] Final output: Processed: Deploy to Production (Approved)
Workflow completed successfully.
Final Result: {Input:Deploy to Production Approved:true Output:Processed: Deploy to Production (Approved)}
```
