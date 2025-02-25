# LLM Text Generation with Ollama

This project demonstrates how to set up and run Large Language Models using Ollama in a containerized environment, with a focus on text generation capabilities.

## Prerequisites

- Azure VM (Standard NV8as v4 - 8 vcpus, 28 GiB memory) with GPU support
- Docker and Docker Compose
- VS Code with Remote SSH extension

## Environment Setup

### System Requirements
1. Set up Docker permissions to run without sudo:
```bash
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker
```

2. Clone the repository:
```bash
git clone <repository-url>
cd /home/azureuser/free-genai-bootcamp-2025/opea-comps/
```

### Environment Variables
Set up the required environment variables:
```bash
export LLM_ENDPOINT_PORT=8008
export TEXTGEN_PORT=9000
export host_ip=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1)
export LLM_ENDPOINT="http://${host_ip}:${LLM_ENDPOINT_PORT}"
export LLM_MODEL_ID="llama3"
export service_name="textgen-service-ollama"
```

## Deployment

### 1. Deploy Ollama Server
```bash
docker compose -f ollama-compose.yaml up -d
```

### 2. Pull the LLM Model
```bash
docker exec ollama-server ollama pull llama2
```

### 3. Test Ollama Server
Verify the server is responding:
```bash
curl --noproxy "*" http://${host_ip}:8008/api/generate -d '{
  "model": "llama2",
  "prompt": "Why is the sky blue?"
}'
```

### 4. Deploy Text Generation Service
```bash
docker compose -f compose_text-generation.yaml up ${service_name} -d
```

## API Usage

### Health Check
```bash
curl http://${host_ip}:${TEXTGEN_PORT}/v1/health_check \
  -X GET \
  -H 'Content-Type: application/json'
```

### Text Generation Examples

#### Streaming Mode
```bash
curl http://${host_ip}:${TEXTGEN_PORT}/v1/chat/completions \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{
      "model": "llama2",
      "messages": [{"role": "user", "content": "What is Deep Learning?"}],
      "max_tokens": 100
    }'
```

#### Non-Streaming Mode
```bash
curl http://${host_ip}:${TEXTGEN_PORT}/v1/chat/completions \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{
      "model": "llama2",
      "messages": [{"role": "user", "content": "What is Deep Learning?"}],
      "max_tokens": 100,
      "stream": false
    }'
```

## Networking Configuration

### Default Bridge Network
When using the default bridge network:
- Container name resolution is not available
- Use container IP addresses for communication
- Alternative: Use custom networks for DNS resolution

### Custom Network Setup
To enable container name resolution, create a custom network:
```bash
docker network create llm-network
```

Update your docker-compose files to use the custom network:
```yaml
networks:
  llm-network:
    external: true
```

## Troubleshooting

1. **Permission Denied**: If you encounter permission errors, ensure your user is in the docker group and you've logged out and back in.

2. **Container Communication**: If containers can't communicate:
   - Verify both containers are on the same network
   - Check if the ports are correctly mapped
   - Ensure the host firewall allows the required ports

3. **Model Not Found**: Verify the model is pulled correctly:
```bash
docker exec ollama-server ollama list
```

## References

- [Docker Post-installation Steps](https://docs.docker.com/engine/install/linux-postinstall/)
- [Ollama Documentation](https://ollama.ai/docs)
- [Docker Networking Guide](https://docs.docker.com/network/)