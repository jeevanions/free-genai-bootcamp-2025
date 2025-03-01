# Setting Up Qdrant for the Language Listening Comprehension App

This guide will help you set up Qdrant, the vector database used for the RAG (Retrieval-Augmented Generation) functionality in our application.

## Prerequisites

1. **Docker**: Qdrant runs in a Docker container, so you need to have Docker installed on your system.
   - [Install Docker for Mac](https://docs.docker.com/desktop/install/mac-install/)
   - [Install Docker for Windows](https://docs.docker.com/desktop/install/windows-install/)
   - [Install Docker for Linux](https://docs.docker.com/desktop/install/linux-install/)

2. **Docker Compose**: Used to manage the Qdrant service.
   - Docker Desktop for Mac and Windows includes Docker Compose
   - For Linux, follow the [Docker Compose installation guide](https://docs.docker.com/compose/install/linux/)

## Starting Qdrant

### Using the Provided Scripts

We've included scripts to make it easy to start and stop Qdrant:

1. Make sure Docker is running on your system
2. Start Qdrant:
   ```bash
   ./scripts/start_qdrant.sh
   ```
3. Stop Qdrant when you're done:
   ```bash
   ./scripts/stop_qdrant.sh
   ```

### Manual Setup

If you prefer to manage Docker Compose manually:

1. Start Qdrant:
   ```bash
   docker compose up -d qdrant
   ```
2. Stop Qdrant:
   ```bash
   docker compose down
   ```

## Verifying Qdrant is Running

1. Check if the Qdrant container is running:
   ```bash
   docker ps | grep qdrant
   ```

2. Access the Qdrant dashboard at: http://localhost:6333/dashboard

3. Test the Qdrant API:
   ```bash
   curl http://localhost:6333/health
   ```
   You should get a response with status code 200. If that doesn't work, you can also try:
   ```bash
   curl http://localhost:6333/collections
   ```

## Troubleshooting

### Docker is not running

If you see an error like:
```
Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
```

Make sure to start Docker Desktop or the Docker service on your system.

### Port conflicts

If port 6333 or 6334 is already in use, you can modify the `docker-compose.yml` file to use different ports:

```yaml
services:
  qdrant:
    ports:
      - "6335:6333"  # Change 6335 to any available port
      - "6336:6334"  # Change 6336 to any available port
```

Then update your `.env` file to use the new port:
```
QDRANT_URL=http://localhost:6335
```

### Data persistence

Qdrant data is stored in a Docker volume named `qdrant_storage`. This ensures your data persists between container restarts.

To remove all data and start fresh:
```bash
docker compose down -v
```

## Using Qdrant in the Application

Once Qdrant is running, the application will automatically connect to it when you start the app:

```bash
python -m frontend.ui.app
```

If Qdrant is not available, the application will still run, but the RAG functionality will be limited.

## Further Information

- [Qdrant Documentation](https://qdrant.tech/documentation/)
- [Docker Documentation](https://docs.docker.com/)
