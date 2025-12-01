# Ephemeral Channels (Values) Example

## Background

In stateful applications, not all data needs to be persisted throughout the entire conversation history. Some data is "ephemeral"—it is relevant only for the immediate next step or within a specific super-step (parallel execution block) and should be discarded afterwards. Examples include temporary search results, intermediate reasoning steps, or flags that trigger immediate actions but shouldn't confuse future turns. **Ephemeral Channels** provide a way to manage this lifecycle automatically.

## Features

*   **Automatic Cleanup**: Channels marked as ephemeral are automatically cleared from the state after the step completes.
*   **Scope Isolation**: Prevents temporary data from leaking into future execution steps, reducing context pollution.
*   **Configurable**: Defined via `RegisterChannel` with an `isEphemeral` flag.

## Implementation Principle

This feature is implemented via the `CleaningStateSchema` interface in `graph/schema.go`.
1.  `MapSchema` maintains a set of `EphemeralKeys`.
2.  During the graph execution loop (`InvokeWithConfig`), after all nodes in the current step have executed and their results merged:
3.  The `Cleanup(state)` method is called.
4.  This method returns a new state map with all keys present in `EphemeralKeys` removed.

## Code Walkthrough

In `main.go`:

1.  **Schema Definition**:
    ```go
    schema := graph.NewMapSchema()
    // Register "temp_data" as ephemeral
    schema.RegisterChannel("temp_data", graph.OverwriteReducer, true)
    ```
    We explicitly tell the schema that `temp_data` should be treated as ephemeral.

2.  **Producer Node**:
    Sets `temp_data` to "secret_code_123".

3.  **Consumer Node**:
    Checks for `temp_data`. Since the producer runs in Step 1, and the consumer runs in Step 2 (sequentially), the cleanup logic triggers between them.
    *Note: In LangGraph semantics, ephemeral values are cleared *after* the step. If Producer -> Consumer is a direct transition, they might be considered separate steps.*

4.  **Verification**:
    The output shows that the consumer (running in the next step) does **not** see the `temp_data`, confirming it was cleaned up.

## How to Run

```bash
go run main.go
```
