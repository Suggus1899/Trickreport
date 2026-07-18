#!/usr/bin/env bash
# Trickreport local dev script (bash)
# Starts backend and frontend concurrently. Press Ctrl+C to stop both.

set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo -e "\033[36mStarting Trickreport dev environment (local, no Docker)...\033[0m"
echo ""

# Check prerequisites
if ! command -v go &>/dev/null; then
    echo -e "\033[31mERROR: Go is not installed or not in PATH\033[0m"
    exit 1
fi
if ! command -v npm &>/dev/null; then
    echo -e "\033[31mERROR: npm is not installed or not in PATH\033[0m"
    exit 1
fi

# Check backend .env
if [ ! -f "$ROOT/backend/.env" ]; then
    echo -e "\033[33mWARNING: backend/.env not found. Copying from .env.example\033[0m"
    cp "$ROOT/backend/.env.example" "$ROOT/backend/.env"
    echo -e "\033[33mEdit backend/.env with your PostgreSQL credentials before running again.\033[0m"
    exit 1
fi

# Check frontend node_modules
if [ ! -d "$ROOT/frontend/node_modules" ]; then
    echo -e "\033[36mInstalling frontend dependencies...\033[0m"
    (cd "$ROOT/frontend" && npm install)
fi

# Cleanup on exit
cleanup() {
    echo ""
    echo -e "\033[33mStopping services...\033[0m"
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    wait $BACKEND_PID 2>/dev/null || true
    wait $FRONTEND_PID 2>/dev/null || true
    echo -e "\033[32mStopped.\033[0m"
}
trap cleanup EXIT INT TERM

# Start backend
echo -e "\033[32m[backend] Starting on http://localhost:8080 ...\033[0m"
(cd "$ROOT/backend" && go run cmd/api/main.go) &
BACKEND_PID=$!

# Start frontend
echo -e "\033[32m[frontend] Starting on http://localhost:4321 ...\033[0m"
(cd "$ROOT/frontend" && npm run dev) &
FRONTEND_PID=$!

echo ""
echo -e "\033[36mBoth services are running. Press Ctrl+C to stop.\033[0m"
echo ""

wait
