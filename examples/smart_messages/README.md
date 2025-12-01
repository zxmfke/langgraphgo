# Smart Message Merging Example

## Background

In modern LLM applications, especially those involving tool use (Function Calling), managing the message history is critical. A common pattern involves an LLM generating a "tool call" message, followed by the tool execution result. Sometimes, we want to **update** a previous message (e.g., replacing a "Thinking..." placeholder with a final answer) rather than appending a new one. Standard list appending is insufficient for this. The **AddMessages Reducer** provides intelligent merging logic.

## Features

*   **ID-Based Deduplication**: If a new message has the same ID as an existing one, it updates (overwrites) the existing message instead of appending.
*   **Automatic Appending**: Messages without IDs (or with new IDs) are appended to the end of the list.
*   **Flexibility**: Works with both `llms.MessageContent` (standard) and custom map-based message structures.

## Implementation Principle

The `AddMessages` function in `graph/add_messages.go` implements the reducer logic:
1.  It iterates through the new messages.
2.  For each message, it attempts to extract an ID (via `MessageWithID` interface, map key "id", or struct field "ID").
3.  If an ID is found and matches an existing message in the current state list, the existing message is replaced at its original index.
4.  Otherwise, the new message is appended.

## Code Walkthrough

In `main.go`:

1.  **Graph Setup**:
    ```go
    g := graph.NewMessagesStateGraph()
    ```
    This helper creates a graph where the "messages" key is pre-configured to use `AddMessages`.

2.  **AI Response Node**:
    Returns a message with `id: "msg_123"` and content "Thinking...".

3.  **AI Update Node**:
    Returns a message with the **same** `id: "msg_123"` but content "Hello...".

4.  **Result**:
    Because the IDs match, the final history contains only **one** AI message (the updated one), not two.

## How to Run

```bash
go run main.go
```
