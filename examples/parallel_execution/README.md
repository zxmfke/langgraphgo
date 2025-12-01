# Parallel Execution Example

This example demonstrates LangGraphGo's ability to execute nodes in parallel.

## Overview

When multiple nodes share the same starting node (fan-out), LangGraphGo automatically executes them concurrently. The results are then merged into the state using the configured schema or reducers.

This example shows a simple fan-out/fan-in pattern.

## Usage

```bash
go run main.go
```
