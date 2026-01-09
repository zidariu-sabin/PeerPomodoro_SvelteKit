# PeerPomodoro Project Context

## Project Overview
PeerPomodoro is a collaborative Pomodoro timer application designed to allow users to create and join timer sessions. The project is structured as a monorepo with a distinct separation between the frontend and backend.

- **Frontend:** Built with SvelteKit (using Svelte 5) and styled with Tailwind CSS.
- **Backend:** Built with Go, utilizing `gorilla/websocket` for real-time communication.

**Note:** There appears to be an architectural divergence. The SvelteKit `README.md` and `package.json` reference a Node.js server with `socket.io` and `socket.io-client`. However, the presence of `PeerPomodoroBackend_Go` and its use of `gorilla/websocket` suggests a move towards a Go-based backend using standard WebSockets. Be aware of this potential protocol mismatch when developing.

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

The frontend is a standard SvelteKit application.

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

**Build:**
```bash
pnpm build
```

**Key Configuration:**
-   `vite.config.ts`: Vite configuration.
-   `svelte.config.js`: SvelteKit configuration.
-   `tailwind.config.ts`: Tailwind CSS configuration.

### Backend (`PeerPomodoroBackend_Go`)

The backend is a Go application acting as a WebSocket server.

**Running the Server:**
```bash
cd PeerPomodoroBackend_Go
# Set allowed origins if developing locally (default fallback exists in code)
export ALLOWED_ORIGINS="http://localhost:5173" 
go run server/main.go
```
The server defaults to listening on `localhost:8080`.

**Endpoints:**
-   `/connect`: WebSocket connection endpoint.

## Development Conventions

-   **Frontend Style:** The project uses Tailwind CSS for styling. Code formatting is handled by Prettier and linting by ESLint.
-   **Backend Structure:**
    -   `server/`: Contains the entry point (`main.go`).
    -   `internal/`: Likely contains the core business logic (though `services/` was imported in `main.go`, standard Go layout often puts app code in `internal/`). *Correction based on `main.go` import:* The code imports `PeerPomodoroBackend_Go/services/connection_manager`, suggesting a structure where services are at the root or explicitly named packages.
    -   `go.mod`: Manages Go dependencies.

## Architecture Notes
-   **Real-time Communication:** The backend uses `gorilla/websocket` with a custom `connectionManager`. The frontend currently lists `socket.io-client` as a dependency. To make them communicate, the frontend will likely need to switch to utilizing the native `WebSocket` API or a library compatible with standard WebSockets, removing the `socket.io-client` dependency.



## Important User Notes

Kepp in mind that the socke.io will not be used for the sveltekit frontend application  as we decided to use gorilla/websockets.
