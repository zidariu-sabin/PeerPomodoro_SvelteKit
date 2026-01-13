# Peer Pomodoro

A collaborative Pomodoro timer application where users can create and join timer sessions to work together.

This project uses SvelteKit for the frontend and a Golang server with native WebSockets and gorilla/WebSockets for real-time communication.

## Implementation Guide

This project uses a Go backend for real-time communication via native WebSockets.

### 1. Prerequisites

-   **Frontend:** `pnpm install`
-   **Backend:** Go 1.25+ installed. Run the server in `../PeerPomodoroBackend_Go`.

### 2. Architecture

-   **Frontend:** SvelteKit (Svelte 5) with Tailwind CSS.
-   **Backend:** Go with `gorilla/websocket`.
-   **Protocol:** Native WebSockets using JSON payloads.

### 3. Connection Handling

State is managed via Svelte 5 reactive stores (`src/lib/stores/connectionStore.svelte.ts`).

-   **Store:** Manages the single WebSocket connection.
-   **Reactivity:** Uses `$state` runes to expose connection status and messages.

### 4. Development Workflow

1.  Start Backend: `cd ../PeerPomodoroBackend_Go && go run server/main.go`
2.  Start Frontend: `pnpm dev`
3.  Access: `http://localhost:5173`
