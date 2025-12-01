# 自定义 Reducer 示例

本示例演示如何在 LangGraphGo 中使用自定义 Reducer 进行状态管理。

## 概述

Reducer 定义了如何将节点的更新合并到当前状态中。虽然 `AppendReducer` 常用于消息列表，但您可以为任何状态字段定义自定义逻辑。

本示例展示了如何实现一个自定义 Reducer 来聚合值或处理复杂的合并策略。

## 用法

```bash
go run main.go
```
