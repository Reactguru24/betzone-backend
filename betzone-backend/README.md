# Betzone Backend API

A Go-based REST API backend for the Betzone sports betting application, built with Gin web framework.

## Project Structure

```
betzone-backend/
├── config/          # Configuration management
├── handlers/        # HTTP request handlers
├── models/          # Data models and structures
├── services/        # External API clients and business logic
├── utils/           # Utility functions
├── main.go          # Application entry point
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── Makefile         # Build and task automation
├── Dockerfile       # Container configuration
└── .env             # Environment variables
```

## Prerequisites

- Go 1.21 or higher
- Make (optional but recommended)

## Installation

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Setup environment variables:**
   - The `.env` file is already configured with Betkraft API credentials

## Development

### Run Development Server

```bash
make dev
```

Or directly with Go:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Build for Production

```bash
make build
```

This creates an executable in `bin/betzone-backend`

### Run Production Build

```bash
make run
```

## API Endpoints

### Health Check
- `GET /health` - Check API health status

### Games
- `GET /api/v1/games` - Get all available games
- `GET /api/v1/games/:id` - Get specific game details

### Bets
- `POST /api/v1/bets` - Create a new bet
- `GET /api/v1/bets` - Get user bets
- `GET /api/v1/bets/:id` - Get specific bet

### Odds
- `GET /api/v1/odds/:gameId` - Get odds for a specific game

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `ENVIRONMENT` | Environment mode | development |
| `BETKRAFT_BASE_URL` | Betkraft API base URL | https://api.staging.betkraft.co.uk |
| `API_KEY` | Betkraft API key | - |
| `APP_KEY` | Betkraft app key | - |

## Docker Support

### Build Docker Image
```bash
make docker-build
```

### Run Docker Container
```bash
make docker-run
```

## Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the built application
make dev           # Run in development mode
make test          # Run tests
make clean         # Clean build artifacts
make deps          # Download dependencies
make lint          # Format code
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

## Testing

Run tests with:
```bash
make test
```

Or directly:
```bash
go test ./...
```

## Code Style

The project uses Go's standard formatting. Run linter with:
```bash
make lint
```

## Dependencies

- **gin-gonic/gin** - HTTP web framework
- **joho/godotenv** - Environment variable loader

## TODO

- [ ] Implement Betkraft API integration
- [ ] Add database integration
- [ ] Implement JWTauthentication
- [ ] Add request validation
- [ ] Add error handling middleware
- [ ] Add logging middleware
- [ ] Write unit tests
- [ ] Add API documentation (Swagger)
- [ ] Implement rate limiting
- [ ] Add caching layer

## Contributing

1. Create feature branches
2. Keep code clean and documented
3. Write tests for new features
4. Submit pull requests

## License

MIT

