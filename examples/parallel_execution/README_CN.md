# 并行执行示例

本示例演示 LangGraphGo 并行执行节点的能力。

## 概述

当多个节点共享同一个起始节点（扇出）时，LangGraphGo 会自动并发执行它们。然后，结果将使用配置的 Schema 或 Reducer 合并到状态中。

本示例展示了一个简单的扇出/扇入模式。

## 用法

```bash
go run main.go
```
