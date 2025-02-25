
python3 -m venv venv

source venv/bin/activate


# What components we need

 - Vector Database
 - Embedding Model
 - LLM

# First getting the LLM into a Microservice

Starting the Ollam and run llama3 modle

```bash
# For ubuntu
export host_ip=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1)

# for mac
export host_ip=$(ipconfig getifaddr en0)

export LLM_ENDPOINT_PORT=8008
export LLM_ENDPOINT="http://${host_ip}:${LLM_ENDPOINT_PORT}"
export LLM_MODEL_ID="llama3"

docker-compose up -d 
docker exec ollama-server ollama list
docker exec ollama-server ollama pull llama3
docker exec ollama-server ollama list

curl ${LLM_ENDPOINT}/v1/chat/completions \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{
      "model": "llama3",
      "messages": [{"role": "user", "content": "What is Deep Learning?"}],
      "max_tokens": 100
    }'

export LLM_SERVICE_HOST_IP=$(ipconfig getifaddr en0)
export LLM_SERVICE_PORT=8008
export OTEL_SDK_DISABLED=true

``

## Test app 


```bash
curl http://${host_ip}:8888/v1/chatqna \
-H "Content-Type: application/json" \
-d '{
    "messages": [
        {
            "role": "user",
            "content": "What is the revenue of Nike in 2023?"
        }
    ]
}'

curl http://${host_ip}:8888/v1/chatqna \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "What is the revenue of Nike in 2023?"
      }
    ],
    "stream": false
  }'
```

# What are we building?

Featcure: Chat with pdf

User Journey 1: First user uploads a pdf document which will be processed offline and stored in our vector database

    This we can use one of the DataPrep microservices


User Journey 2: User initiates a chat with the application where relevant vectors are queried and sent to LLM with user questions for a response.


## Test Qdrant DataPrep Microservice

```bash
curl -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "files=@./file1.txt" \
    http://${host_ip}:6007/v1/dataprep/ingest
```

You can specify chunk_size and chunk_size by the following commands.

```bash
curl -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "files=@./file1.txt" \
    -F "chunk_size=1500" \
    -F "chunk_overlap=100" \
    http://${host_ip}:6007/v1/dataprep/ingest
```