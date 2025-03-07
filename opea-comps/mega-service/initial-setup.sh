# Run these commands manually in the sequence provided

sudo apt-get update
sudo apt-get install tesseract-ocr -y
sudo apt-get install libtesseract-dev -y
sudo apt-get install poppler-utils -y
sudo apt install python3-pip
sudo apt install python3.12-venv

wget https://raw.githubusercontent.com/opea-project/GenAIExamples/refs/heads/main/ChatQnA/docker_compose/install_docker.sh
chmod +x install_docker.sh
./install_docker.sh

sudo apt  install docker-compose

python3 -m venv .venv

cd .venv/bin
source activate
cd ../..

pip install -r requirements.txt


# For ubuntu
export HOST_IP=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1)

# for mac
export HOST_IP=$(ipconfig getifaddr en0)

# Set secrets & keys for external services
export HF_TOKEN=<hf_token>

# Modles to be used
export EMBEDDING_MODEL_ID="BAAI/bge-base-en-v1.5"
export EMBEDDING_DIMENSION=768
export QDRANT_EMBED_DIMENSION=768
export LLM_MODEL_ID="llama3"
export RERANK_MODEL_ID="BAAI/bge-reranker-base"

# Embedding Collection Name

export COLLECTION_NAME="PDF_COLLECTION"
export INDEX_NAME="PDF_INDEX"
export QDRANT_INDEX_NAME="PDF_COLLECTION"

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
export EMBEDDING_SERVICE_ENDPOINT="${HOST_IP}:${EMBEDDING_SERVICE_PORT}"
export LLM_SERVICE_ENDPOINT="${HOST_IP}:${LLM_SERVICE_PORT}"
export VECTORDB_QDRANT_SERVICE_ENDPOINT="${HOST_IP}:${VECTORDB_QDRANT_SERVICE_PORT}"
export DATAPREP_SERVICE_ENDPOINT="${HOST_IP}:${DATAPREP_MICROSERVICE_PORT}"
export RETRIEVER_SERVICE_ENDPOINT="${HOST_IP}:${RETRIEVER_PORT}"
export RERANKER_SERVICE_ENDPOINT="${HOST_IP}:${RERANKER_PORT}"

export TAG=1.2 # Latest update to the Dataprep is broken

# Disable open telemetry for now
export OTEL_SDK_DISABLED=true

sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker


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



# To Inject Qdrant using microservice

curl -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "files=@./EthicalAI.pdf" \
    -F "chunk_size=1500" \
    -F "chunk_overlap=100" \
    http://${DATAPREP_SERVICE_ENDPOINT}/v1/dataprep/ingest

# To Test the Qdrant via microservice



# To Test Embedding Service

curl -X POST \
-H 'Content-Type: application/json' \
-d '{"input":"What is Deep Learning?"}' \
http://${EMBEDDING_SERVICE_ENDPOINT}/v1/embeddings


# To Test retriver
export test_embedding=$(python3 -c "import random; embedding = [random.uniform(-1, 1) for _ in range(768)]; print(embedding)")

curl http://${RETRIEVER_SERVICE_ENDPOINT}/v1/retrieval \
-X POST \
-d "{\"text\":\"What is Inductive Biases?\",\"embedding\":${test_embedding}}" \
-H 'Content-Type: application/json'


curl http://${RETRIEVER_SERVICE_ENDPOINT}/v1/health_check \
  -X GET \
  -H 'Content-Type: application/json'

# To Test reranker

curl http://${RERANKER_SERVICE_ENDPOINT}/v1/reranking \
-X POST \
-d '{"initial_query":"What is Deep Learning?", "retrieved_docs": [{"text":"Deep Learning is not..."}, {"text":"Deep learning is..."}]}' \
-H 'Content-Type: application/json'

# Setup Needed to run the chat app

export MEGA_SERVICE_PORT=8888
export EMBEDDING_SERVICE_HOST_IP="${HOST_IP}"
export EMBEDDING_SERVICE_PORT=8007
export LLM_SERVICE_HOST_IP="${HOST_IP}"
export LLM_SERVICE_PORT=8008
export RETRIEVER_SERVICE_HOST_IP="${HOST_IP}"
export RETRIEVER__SERVICE_PORT=8006
export RERANKER_SERVICE_HOST_IP="${HOST_IP}"
export RERANKER_SERVICE_PORT=8005


curl http://${HOST_IP}:${MEGA_SERVICE_PORT}/v1/chatqna \
    -H "Content-Type: application/json" \
    -d '{
        "messages": "What is Inductive Biases?"
    }'



# To Run retriver from locally

export QDRANT_HOST=${your_qdrant_host_ip}
export QDRANT_PORT=6333
export EMBED_DIMENSION=${your_embedding_dimension}
export INDEX_NAME=${your_index_name}
export TEI_EMBEDDING_ENDPOINT="http://${your_ip}:6060"
export RETRIEVER_COMPONENT_NAME="OPEA_RETRIEVER_QDRANT"