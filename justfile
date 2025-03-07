# Run frontend with vite for hot reload
dev-frontend:
  cd frontend && pnpm run dev

# Run backend with air for hot reload
dev-backend:
  cd backend && air

# Build frontend
run-frontend:
  cd frontend && pnpm run build
  cd frontend && pnpm run preview

# Build backend
run-backend:
  cd backend && go build -o backend
  cd backend && ./backend

# Install dependencies
install:
  cd frontend && pnpm install
  cd backend && go mod tidy
  go install github.com/air-verse/air@latest
