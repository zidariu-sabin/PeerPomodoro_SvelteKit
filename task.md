# Implementation Plan: PeerPomodoro

This document outlines the step-by-step implementation plan for the PeerPomodoro project, following an Agile approach. The immediate priority is resolving the communication protocol mismatch between the SvelteKit frontend and the Go backend.

## Phase 1: Foundation & Connectivity ("Walking Skeleton")
**Goal:** Establish a successful, persistent WebSocket connection between the SvelteKit frontend and the Go backend.

- [x] **1.1 Cleanup Frontend Dependencies**
    - Remove `socket.io-client` and `socket.io` from `PeerPomodoro_Sveltekit/package.json`.
    - Install `svelte-sonner` or a similar notification library (optional but useful for connection status).

- [x] **1.2 Backend Connection Verification**
    - Ensure the Go server's `/connect` endpoint is correctly upgrading HTTP requests to WebSockets.
    - Verify CORS settings in `main.go` allow connections from the SvelteKit dev server (`http://localhost:5173`).

- [x] **1.3 Frontend WebSocket Store**
    - Create a Svelte 5 reactive store (using `$state`) in `src/lib/stores/connectionStore.svelte.ts`.
    - Implement native `WebSocket` logic to connect to `ws://localhost:8080/connect`.
    - Handle standard events: `onopen`, `onmessage`, `onclose`, `onerror`.

- [x] **1.4 "Hello World" Integration Test**
    - Create a simple component on the homepage that displays the connection status (Connected/Disconnected).
    - Send a test message from Frontend -> Backend and log it on the server console.

## Phase 2: Session Management (MVP Core)
**Goal:** Allow users to initiate and join isolated sessions (Sessions).

- [x] **2.1 Backend Session Infrastructure**
    - Define a `Session` struct in Go to hold connected clients and timer state.
    - Implement a `SessionManager` to handle creating and retrieving sessions.
    - Define JSON message structures for `create_session` and `join_Session`.
    - **Update:** Adopted UUIDv4 for Session IDs.
    - Create `/create_Session` endpoint logic to instantiate a new Session and return its UUID.

- [x] **2.2 Frontend Session UI**
    - Update the homepage (`/`) to feature "Create Session" form.
    - Implement navigation: On successful session creation, redirect user to `/timer/[sessionId]`.
    - Upon navigating to `/timer/[sessionId]`, create a WebSocket connection.
- [x] **2.3 Wiring Session Logic**
    - Frontend: Upon timer form submit, request `/create_Session` 
    - Backend: Respond with a unique `Session_id` (UUID).
    - Frontend: Upon navigating to `/timer/[sessionId]`, send `join_Session` with the UUID.
    - Backend: Add the client's connection to the specific Session's broadcast list.

## Phase 3: Timer Synchronization
**Goal:** Enable real-time, synchronized timer control across all clients in a Session.

- [x] **3.1 Backend Timer Logic**
    - Extend `Session` struct to track `TimerState` (running/paused, time remaining, last updated timestamp).
    - Implement message handlers for `start_timer`, `pause_timer`, `reset_timer`.
    - Create a ticker/loop (or efficient timestamp comparison) to broadcast time updates.

- [x] **3.2 Frontend Timer Component**
    - Build a `PomodoroTimer.svelte` component.
    - Bind "Start", "Pause", "Reset" buttons to WebSocket messages.
    - Subscribe to `timer_update` messages from the server to update the UI display.

- [x] **3.3 Handling Late Joiners**
    - Ensure that when a user joins an active Session, they immediately receive the current timer state (syncing them up with the group).

## Phase 4: Polish & UX Improvements
**Goal:** Make the application robust and pleasant to use.

- [x] **4.1 Graceful Disconnection**
    - Backend: Remove clients from Sessions on disconnect.
    - Frontend: Show a "Reconnecting..." state if the server drops.

- [ ] **4.2 Visuals & Feedback**
    - Create a user Dashboard to display session information about joined users and their connection state.
    - Keep the already implemented design and make it consistent.

- [x] **4.3 Notifications**
    - Add browser notifications for timer completion and phase changes.

- [x] **4.4 Session Cleanup**
    - Implement background job to remove inactive sessions (Garbage Collector pattern).
    - Add LastActivity tracking to sessions.


