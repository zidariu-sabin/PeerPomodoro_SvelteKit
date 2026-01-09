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

## Development Conventions

-   **Frontend Style:** Tailwind CSS for styling. Svelte 5 runes for state management. Native `WebSocket` API for communication.
-   **Backend Architecture:**
    -   `server/`: Application entry point (`main.go`).
    -   `internal/domain/`: Core entities and business rules (Timer, Session, Client).
    -   `internal/service/`: Business logic orchestration (SessionManager).
    -   `internal/transport/websocket/`: WebSocket infrastructure (Hub, Client).
    -   `internal/adapters/`: External system integrations (e.g., in-memory repos).

## Architecture Notes
-   **Real-time Communication:** The system uses standard WebSockets. 
    -   **Backend:** Uses `gorilla/websocket`. A `Hub` manages broadcasting and connection lifecycle.
    -   **Frontend:** Uses native `WebSocket` API integrated into a Svelte 5 reactive store (`connectionStore.svelte.ts`).
-   **Concurrency:** The backend Hub uses a non-blocking broadcast pattern to protect against slow consumers.
