# Peer Pomodoro

A collaborative Pomodoro timer application where users can create and join timer sessions to work together.

This project uses SvelteKit for the frontend and a Node.js server with Socket.IO for real-time communication.

## To-Do / Implementation Steps

Here is a step-by-step guide to building the backend and integrating it with your SvelteKit frontend.

### 1. Project Setup and Dependencies

-   **Install necessary packages:** You'll need `socket.io` for the server and `socket.io-client` for the client-side. You may also want `@types/ws` for type support in your custom server file.
    ```bash
    pnpm install socket.io socket.io-client @types/ws
    ```

### 2. Create a Custom Server for WebSocket Integration

SvelteKit uses Vite, and to handle WebSocket connections, you need to create a custom server that runs alongside SvelteKit's development server.

-   Create a new file in your project's root directory named `server.ts`.
-   In this file, you will:
    -   Create a standard Node.js `http` server.
    -   Initialize a `Socket.IO` server and attach it to the `http` server.
    -   Configure the SvelteKit `handler` to work as middleware for regular HTTP requests.
    -   Import and start your WebSocket connection logic from a separate file.

### 3. Implement the WebSocket Server Logic

To keep your code organized, it's best to manage the WebSocket logic in a separate module.

-   Create a new file at `src/lib/server/webSocketServer.ts`.
-   This module will export a function that takes the `Socket.IO` server instance as an argument.
-   **Core responsibilities of this module:**
    -   **Session Management:**
        -   Keep track of active timer "rooms." A simple in-memory `Map` is a good starting point, where the key is the `roomId` and the value is the timer's state.
        -   Each room's state should include details like time remaining and whether the timer is running.
    -   **Event Handling:**
        -   Listen for the initial `connection` event from new clients.
        -   Handle `disconnect` events to manage users leaving rooms.

### 4. Define the WebSocket Communication Protocol

A clear set of events is crucial for reliable communication between the client and server.

-   **Client-to-Server Events:**
    -   `createRoom`: Sent when a user wants to create a new timer session.
    -   `joinRoom(roomId)`: Sent when a user attempts to join an existing session.
    -   `timerStateChange(roomId, newState)`: Sent when a user starts, stops, or resets the timer. The `newState` object contains the updated timer information.
-   **Server-to-Client Events:**
    -   `roomCreated(roomId)`: Sent back to the client that created the room, providing the unique ID for the new session.
    -   `timerStateUpdate(newState)`: Broadcast to all clients within a specific room whenever the timer's state changes.
    -   `error(message)`: Sent to a client if an action fails (e.g., trying to join a room that doesn't exist).

### 5. Integrate WebSocket Client in the Svelte Frontend

-   Create a Svelte store (e.g., `src/lib/stores/socketStore.ts`) to act as a centralized place for managing the WebSocket connection and state on the client-side.
-   This store will:
    -   Establish and maintain a connection to the WebSocket server.
    -   Expose methods for sending events to the server, such as `createRoom` or `joinRoom`.
    -   Listen for events from the server (like `timerStateUpdate`) and update its own state, allowing your Svelte components to reactively update.

### 6. Update SvelteKit Routing and Pages

-   **`/` (The Home Page):**
    -   This page should feature a "Create a new session" button.
    -   Clicking this button will trigger the `createRoom` event via your `socketStore`.
    -   After the server responds with the `roomCreated(roomId)` event, programmatically navigate the user to the new timer page at `/timer/[roomId]`.
-   **`/timer/[slug]` (The Timer Page):**
    -   When this page loads, extract the `slug` (which is the `roomId`) from the URL.
    -   Use this `slug` to call the `joinRoom` method from your `socketStore`.
    -   The timer display and control buttons (play, pause) should be bound to the reactive state managed by your `socketStore`.
    -   User interactions (like clicking "pause") will send the `timerStateChange` event to the server, which then synchronizes the state across all clients in that room.
