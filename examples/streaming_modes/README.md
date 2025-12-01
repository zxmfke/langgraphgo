# Streaming Modes Example

## Background

For long-running LLM chains or complex agent workflows, waiting for the final result is often unacceptable. Users expect real-time feedback. LangGraph Go supports a robust streaming architecture that allows you to subscribe to events as they happen. Different use cases require different "views" of the stream—sometimes you want every token, sometimes just the final output of each node. **Streaming Modes** allow you to configure this granularity.

## Features

*   **`StreamModeUpdates`**: Emits the output of each node as it completes. Useful for showing progress (e.g., "Step 1 done", "Tool executed").
*   **`StreamModeValues`**: Emits the full graph state after each step. Useful for debugging or UIs that render the entire context.
*   **`StreamModeMessages`**: (Planned) Emits LLM tokens for typewriter effects.
*   **`StreamModeDebug`**: Emits all internal events for deep inspection.

## Implementation Principle

The streaming logic is handled by `StreamingRunnable` in `graph/streaming.go`.
1.  It injects a `StreamingListener` (which implements `GraphCallbackHandler` and `NodeListener`) into the execution context.
2.  As the graph executes, nodes and the graph engine emit events (`NodeEventComplete`, `OnGraphStep`, etc.).
3.  The `StreamingListener.shouldEmit` method filters these events based on the configured `StreamMode`.
4.  Allowed events are sent to a Go channel returned to the caller.

## Code Walkthrough

In `main.go`:

1.  **Configuration**:
    ```go
    g.SetStreamConfig(graph.StreamConfig{
        Mode: graph.StreamModeUpdates,
    })
    ```
    We configure the graph to stream in `updates` mode.

2.  **Execution**:
    ```go
    streamResult := runnable.Stream(context.Background(), nil)
    for event := range streamResult.Events { ... }
    ```
    We iterate over the `Events` channel.

3.  **Output**:
    You will see events printed in real-time as `step_1` and `step_2` complete, rather than waiting for the entire graph to finish.

## How to Run

```bash
go run main.go
```
