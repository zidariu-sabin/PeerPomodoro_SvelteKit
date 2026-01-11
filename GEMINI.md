# PeerPomodoro Project Context

## Project Overview
PeerPomodoro is a collaborative Pomodoro timer application designed to allow users to create and join timer sessions. The project is structured as a monorepo with a distinct separation between the frontend and backend.

- **Frontend:** Built with SvelteKit (using Svelte 5) and styled with Tailwind CSS.
- **Backend:** Built with Go, utilizing `gorilla/websocket` for real-time communication.

## Directory Structure

```text
PeerPomodoro/
├── PeerPomodoro_Sveltekit/    # SvelteKit Frontend application
├── PeerPomodoroBackend_Go/    # Go Backend application
└── GEMINI.md                  # This context file
```

## Getting Started

### Prerequisites
-   **Node.js & pnpm:** Required for the frontend.
-   **Go (1.25+):** Required for the backend.

### Frontend (`PeerPomodoro_Sveltekit`)

The frontend is a standard SvelteKit application using Svelte 5 runes.

**Installation:**
```bash
cd PeerPomodoro_Sveltekit
pnpm install
```

**Development Server:**
```bash
pnpm dev
```
Runs on `http://localhost:5173` (default).

### Backend (`PeerPomodoroBackend_Go`)

The backend is a Go application acting as a WebSocket server.

**Configuration:**
-   Create a `.env` file (copy from `.example.env`).
-   `ALLOWED_ORIGINS`: Comma-separated list of allowed origins (e.g., `http://localhost:5173`).

**Running the Server:**
```bash
cd PeerPomodoroBackend_Go
go run server/main.go
```
The server listens on `localhost:8080`.

**Endpoints:**
-   `/connect`: WebSocket connection endpoint.
-   `/create-session`: HTTP POST endpoint for initial session creation.

## Development Conventions

-   **Frontend Style:** Tailwind CSS for styling. Svelte 5 runes for state management. Native `WebSocket` API for communication.
-   **Frontend Routing:** 
    -   `/`: Homepage with session creation.
    -   `/timer`: Local-only Pomodoro timer.
    -   `/timer/[slug]`: Collaborative session timer (uses WebSocket).
-   **Backend Architecture:**
    -   `server/`: Application entry point (`main.go`).
    -   `internal/domain/`: Core entities and business rules (Timer, Session, Client, Message, Protocol).
    -   `internal/ports/`: Interface definitions for external systems (e.g., `SessionRepository`).
    -   `internal/service/`: Business logic orchestration (e.g., `SessionService`).
    -   `internal/transport/http/`: HTTP endpoint handlers (e.g., `CreateSession`).
    -   `internal/transport/websocket/`: WebSocket infrastructure (Hub, Client).
    -   `internal/adapters/`: Implementation of ports (e.g., `in_memory/SessionRepository`).

## Architecture Notes
-   **Real-time Communication:** 
    -   **Backend:** Uses `gorilla/websocket`. A session-aware `Hub` manages scoped broadcasting to specific sessions.
    -   **Frontend:** Uses native `WebSocket` API integrated into a Svelte 5 reactive store (`connectionStore.svelte.ts`).
-   **Timer Modes:**
    -   **Solo Mode:** The timer runs locally using `pomodoroStore.svelte.ts`. No server connection is required.
    -   **Collaborative Mode:** When a session is active (`connection.session.id` exists), timer actions are delegated to the backend via `connectionStore.svelte.ts`.
    -   **Abstraction:** `src/lib/controllers/timerController.ts` handles the switching logic, keeping the UI components agnostic of the mode.
-   **Session Management:**
    -   **Lifecycle:** Sessions are created via an HTTP POST request and joined via a WebSocket `join_session` message.
    -   **Storage:** Decoupled via the Repository pattern. Currently uses an in-memory adapter but is designed for easy transition to a database (SQL/NoSQL).
    -   **Identification:** Sessions are identified by UUIDv4 strings.
-   **Communication Protocol:**
    -   All WebSocket messages follow a strict JSON structure defined by the `domain.Message` type, containing a `type` string and a `payload` raw JSON object.
-   **Concurrency:** The backend Hub uses a non-blocking broadcast pattern and session-scoped maps to ensure efficient message routing.
