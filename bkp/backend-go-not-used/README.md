# Italian Learning Platform Backend

## Setup

1. Install Go 1.21 or later
2. Install Mage:
   ```bash
   go install github.com/magefile/mage@v1.15.0
   ```
3. Install Swagger tools:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```
4. Add Go's bin directory to your PATH (add to ~/.bashrc or ~/.zshrc):
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
5. Reload your shell:
   ```bash
   source ~/.zshrc  # or source ~/.bashrc
   ```
6. Copy `.env.example` to `.env` and adjust values if needed
7. Create required directories:
   ```bash
   mkdir -p bin
   ```
8. Update dependencies and generate Swagger docs:
   ```bash
   mage tidy
   mage docs
   ```

## Development

- Build: `mage build`
- Run: `mage run`
- Test: `mage test`
- Update Dependencies: `mage tidy`
- Generate API Docs: `mage docs`

## API Documentation

The API documentation is available through Swagger UI when the server is running:

1. Start the server:
   ```bash
   mage run
   ```

2. Access Swagger UI:
   ```
   http://localhost:8080/swagger/index.html
   ```

## Verify Installation

1. Build the application:
   ```bash
   mage build
   ```

2. Run the tests:
   ```bash
   mage test
   ```

3. Start the server:
   ```bash
   mage run
   ```

4. Test the health endpoint:
   ```bash
   curl http://localhost:8080/health
   ```
   
   Expected response:
   ```json
   {"status":"healthy"}
   ```

## API Endpoints

- Health Check: `GET /health`
- API Documentation: `GET /swagger/*any`

## Development Workflow

1. Add Swagger annotations to your handlers
2. Generate Swagger documentation:
   ```bash
   mage docs
   ```
3. Implement your endpoint
4. Build and run:
   ```bash
   mage build
   mage run
   ```
5. Check your new endpoint in Swagger UI

## Swagger Annotation Example

```go
// HealthCheck godoc
// @Summary Check API health
// @Description Get the health status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
    })
}
```
