Did a bit of digging to run the LLM models in a container and make use of the Mac M1 GPU

# Setup

1. Install libkrun-efi and libkrun

```bash
brew tap slp/krunkit
brew install krunkit
```
2. Install Podman desktop application and create a machin with libkrun to access the Mac M1 GPU 

3. Use podman cli 

Prefix each podman machine command with CONTAINERS_MACHINE_PROVIDER=libkrun. For example, CONTAINERS_MACHINE_PROVIDER=libkrun podman machine ls.
Otherwise you wont be able to interact with the podman machine enabled with libkrun.


4. Check the GPU access

 CONTAINERS_MACHINE_PROVIDER=libkrun podman machine ssh


 https://podman-desktop.io/docs/podman/gpu

 

# Usins GPU optimised Azure VM - Ubuntu

Install docker 

https://docs.docker.com/engine/install/linux-postinstall/

export host_ip=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1)


sudo docker compose up


sudo curl --noproxy "*" http://${host_ip}:8008/api/generate -d '{
  "model": "llama3",
  "prompt":"Why is the sky blue?"
}'

{"error":"model 'llama3' not found"}

docker compose -f ollama-compose.yaml up -d
sudo docker exec ollama-server ollama pull llama3


curl http://localhost:8008/api/tags


Create the docker group (if it doesn't exist):
sudo groupadd docker
Add your user to the docker group:
sudo usermod -aG docker $USER
Apply the new group membership:
newgrp docker
Create a .env file to set the environment variables (place it in the same directory as your docker-compose.yaml):
Verify the setup:
docker compose up -d
After these changes, you should be able to run Docker commands without sudo. If you still get permission errors, try:

Check your docker socket permissions:
Logout and login again to ensure all group changes take effect.




-----


git clone https://github.com/opea-project/GenAIComps.git

export OPEA_GENAICOMPS_ROOT=$(pwd)/GenAIComps

# Build the microservice docker
cd ${OPEA_GENAICOMPS_ROOT}


export host_ip=$(ip -4 addr show eth0 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')

export LLM_ENDPOINT_PORT=8008
export TEXTGEN_PORT=9000
export host_ip=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1)
export HF_TOKEN=${HF_TOKEN}
export LLM_ENDPOINT="http://${host_ip}:${LLM_ENDPOINT_PORT}"
export LLM_MODEL_ID="llama3:latest"

<!-- export LLM_MODEL_ID="meta-llama/Llama-2-70b" -->



export service_name="textgen-service-ollama"

cd ../../deployment/docker_compose/
docker compose -f compose_text-generation.yaml up ${service_name} -d

docker compose -f compose_text-generation.yaml down ${service_name}



curl http://${host_ip}:${TEXTGEN_PORT}/v1/health_check\
  -X GET \
  -H 'Content-Type: application/json'

{"Service Title":"opea_service@llm","Service Description":"OPEA Microservice Infrastructure"}

# References

https://medium.com/@andreask_75652/gpu-accelerated-containers-for-m1-m2-m3-macs-237556e5fe0b

https://sinrega.org/2024-03-06-enabling-containers-gpu-macos/


https://docs.docker.com/engine/install/ubuntu/