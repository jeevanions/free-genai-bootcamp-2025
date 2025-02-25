# Run these commands manually in the sequence provided

apt-get update

wget https://raw.githubusercontent.com/opea-project/GenAIExamples/refs/heads/main/ChatQnA/docker_compose/install_docker.sh
chmod +x install_docker.sh
./install_docker.sh

pip install -r requirements.txt
apt-get install tesseract-ocr -y
apt-get install libtesseract-dev -y
apt-get install poppler-utils -y


# For ubuntu
export HOST_IP=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1)

# for mac
export HOST_IP=$(ipconfig getifaddr en0)

# Set secrets & keys for external services
export HF_TOKEN=<Hugging face token>

# Modles to be used
export EMBEDDING_MODEL_ID="BAAI/bge-large-en-v1.5"
export LLM_MODEL_ID="llama3"
export RERANK_MODEL_ID="BAAI/bge-reranker-base"

# Embedding Collection Name

export COLLECTION_NAME="PDF_COLLECTION"
export INDEX_NAME="PDF_INDEX"

# Service Ports
export QDRANT_PORT=6333
export RETRIEVER_PORT=8006
export RERANKER_PORT=8005
export LLM_SERVICE_PORT=8008
export EMBEDDING_SERVICE_PORT=8007
export VECTORDB_QDRANT_SERVICE_PORT=6007
export DATAPREP_MICROSERVICE_PORT=8009

# Service Endpoints
export QDRANT_HOST="${HOST_IP}"
export EMBEDDING_SERVICE_ENDPOINT="http://${HOST_IP}:${EMBEDDING_SERVICE_PORT}"
export LLM_SERVICE_ENDPOINT="${HOST_IP}:${LLM_SERVICE_PORT}"
export VECTORDB_QDRANT_SERVICE_ENDPOINT="${host_ip}:${VECTORDB_QDRANT_SERVICE_PORT}"


# Disable open telemetry for now
export OTEL_SDK_DISABLED=true


docker-compose up -d 

docker exec ollama-server ollama pull llama3

# To test the LLM 

curl http://${LLM_SERVICE_ENDPOINT}/v1/chat/completions \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{
      "model": "llama3",
      "messages": [{"role": "user", "content": "What is Deep Learning?"}],
      "max_tokens": 100
    }'

# To test Qudrant using microservice

curl -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "files=@./file1.txt" \
    -F "chunk_size=1500" \
    -F "chunk_overlap=100" \
    http://${VECTORDB_QDRANT_SERVICE_ENDPOINT}/v1/dataprep/ingest

# To Test Embedding Service





# Setup Needed to run the chat app

export MEGA_SERVICE_PORT=8000
export EMBEDDING_SERVICE_HOST_IP="${HOST_IP}"
export EMBEDDING_SERVICE_PORT=8007
export LLM_SERVICE_HOST_IP="${HOST_IP}"
export LLM_SERVICE_PORT=8008
export RETRIEVER_SERVICE_HOST_IP="${HOST_IP}"
export RETRIEVER__SERVICE_PORT=8006
export RERANKER_SERVICE_HOST_IP="${HOST_IP}"
export RERANKER_SERVICE_PORT=8005

