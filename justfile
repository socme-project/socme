# Build front & backend
build:
  mkdir -p "./build"
  cd back && go build -o "../build/backend" "./cmd/main.go"
  cd front && bun run build && mv "./dist" "../build/frontend"
  echo "Build completed successfully."

# Run dev mode for both front & backend. Live reload enabled.
dev:
  cd front && bun run dev
