# Copilot Instructions for everything-base

## Project Overview
- **everything-base** is a cross-platform desktop app using [Wails](https://wails.io/) (Go backend + modern JS frontend).
- The frontend is built with Vite, React, TypeScript, TailwindCSS, and shadcn/ui components (see `frontend/`).
- The backend is written in Go, with main entry points in `main.go` and `app.go`.
- The project structure separates Go backend logic, frontend assets, and build scripts/assets.

## Key Components & Data Flow
- **Go Backend**: Main logic in `app.go`, `main.go`, and subfolders (`db/`, `process/`, `switchfs/`, etc.).
  - The `App` struct in `app.go` is the main bridge for backend/frontend communication.
  - Backend methods (e.g., `Greet`) are exposed to the frontend via Wails.
- **Frontend**: Located in `frontend/`.
  - Uses Vite for dev/build, React for UI, TailwindCSS for styling, and shadcn/ui for components.
  - Entry: `frontend/src/main.tsx`, main app: `frontend/src/App.tsx`.
- **Build Assets**: `build/` contains platform-specific files for packaging (icons, manifests, installer scripts).

## Developer Workflows
- **Live Development**:
  1. In project root: `wails dev` (starts Go backend with hot reload)
  2. In `frontend/`: `npm run dev` (starts Vite dev server)
  3. Access frontend at http://localhost:34115
- **Build Production App**:
  - Run `wails build` in project root (builds Go backend, bundles frontend, creates distributables)
- **Frontend Component Additions**:
  - Use shadcn's CLI: `npx shadcn-ui@latest add [component]`
- **Scripts**:
  - See `scripts/` for platform-specific build helpers (e.g., `build-windows.sh`, `build-macos.sh`).

## Project-Specific Conventions
- **Backend/Frontend Bridge**: Only methods on the `App` struct are exposed to the frontend.
- **Frontend**: Use TypeScript, Tailwind, and shadcn/ui patterns. Place new components in `frontend/src/components/`.
- **Database/Storage**: Go code in `db/` handles local persistence; see `persistentDB.go` for patterns.
- **Switch File System**: `switchfs/` contains logic for Nintendo Switch file formats.

## Integration Points
- **Wails**: Handles backend/frontend communication and app lifecycle.
- **shadcn/ui**: For UI components, follow https://ui.shadcn.com/docs/cli#add for usage.
- **Platform Packaging**: Custom icons, manifests, and installer scripts in `build/`.

## Examples
- Exposing a Go method to the frontend: see `Greet` in `app.go`.
- Adding a new UI component: use shadcn CLI, then import in `frontend/src/components/`.

## References
- Main entry: `main.go`, `app.go`
- Frontend entry: `frontend/src/main.tsx`, `frontend/src/App.tsx`
- Build config: `wails.json`, `frontend/vite.config.ts`, `frontend/tailwind.config.js`
- Platform assets: `build/`

---

_If you are unsure about a workflow or pattern, check the README.md or the relevant directory for examples._
