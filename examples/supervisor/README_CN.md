# Supervisor Agent 示例

本示例演示如何使用 Supervisor 实现多 Agent 系统。

## 概述

在 Supervisor 模式中，一个中央 "Supervisor" Agent 将任务路由给专门的 Worker Agent。Worker 执行任务并向 Supervisor 汇报，Supervisor 随后决定下一步操作或结束工作流。

本示例展示了如何使用 `prebuilt.CreateSupervisor` 来编排多个 Agent。

## 用法

```bash
go run main.go
```
