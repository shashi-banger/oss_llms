# Text Embedding Service

A containerized text embedding service using the Qwen3-Embedding-0.6B model served through vLLM and accessed via LiteLLM proxy.

## Architecture

- **vLLM Service**: Serves the Qwen/Qwen3-Embedding-0.6B model using vLLM OpenAI-compatible server
- **LiteLLM Proxy**: Provides a unified API interface, proxying requests to the vLLM backend
- **PostgreSQL Database**: Stores LiteLLM configuration and usage data

## Prerequisites

- Docker and Docker Compose
- NVIDIA GPU with CUDA support
- Hugging Face token stored in `~/.hugging_face/sb-hf-read-token`

## Quick Start

1. **Start the services:**
   ```bash
   cd text_embedding
   HUGGING_FACE_HUB_TOKEN=`cat ~/.hugging_face/sb-hf-read-token` docker-compose up -d
   ```

2. **Wait for services to be ready** (usually takes 1-2 minutes for all services to become healthy)

3. **Test with Python:**
   ```bash
   python test_embedding.py
   ```

4. **Test with Go:**
   ```bash
   go run test_embedding.go
   ```

## API Usage

The embedding service exposes an OpenAI-compatible API:

- **Endpoint**: `http://localhost:4000/v1/embeddings`
- **Model name**: `qwen-embedding`
- **Authentication**: Bearer token `sk-sb123`

### Example Request

```bash
curl -X POST "http://localhost:4000/v1/embeddings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-sb123" \
  -d '{
    "model": "qwen-embedding",
    "input": ["Hello, how are you?"]
  }'
```

### Example Response

```json
{
  "object": "list",
  "data": [
    {
      "object": "embedding",
      "embedding": [0.1234, -0.5678, ...],
      "index": 0
    }
  ],
  "model": "Qwen/Qwen3-Embedding-0.6B",
  "usage": {
    "prompt_tokens": 5,
    "total_tokens": 5
  }
}
```

## Test Scripts

### Python Test (`test_embedding.py`)

Tests the embedding API using the OpenAI Python client:
- Sends 4 sample texts for embedding
- Displays model info, embedding dimensions, and first 5 values of each embedding

### Go Test (`test_embedding.go`)

Tests the embedding API using native Go HTTP client:
- Sends 4 sample texts for embedding
- Displays model info, embedding dimensions, token usage, and first 5 values

## Service Management

### Start services
```bash
HUGGING_FACE_HUB_TOKEN=`cat ~/.hugging_face/sb-hf-read-token` docker-compose up -d
```

### Stop services
```bash
docker-compose down
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f litellm
docker-compose logs -f vllm
docker-compose logs -f postgres
```

### Check service status
```bash
docker-compose ps
```

## Troubleshooting

### Services not starting
- Ensure Docker has access to GPU: `docker run --rm --gpus all nvidia/cuda:11.0-base nvidia-smi`
- Check if ports 4000, 5433, 8000 are available
- Verify Hugging Face token is valid and has read access

### Authentication errors
- The master key `sk-sb123` is hardcoded in `litellm_config.yaml`
- LiteLLM requires a database connection for authentication

### Performance issues
- vLLM requires significant GPU memory (recommended: 8GB+ VRAM)
- First request may be slower due to model loading

## Configuration Files

- `docker-compose.yml`: Service orchestration and environment setup
- `litellm_config.yaml`: LiteLLM model routing and API settings
- `test_embedding.py`: Python test client
- `test_embedding.go`: Go test client

## Model Information

- **Model**: Qwen/Qwen3-Embedding-0.6B
- **Embedding Dimension**: 1024
- **Provider**: Hugging Face Transformers via vLLM
- **License**: Check model page on Hugging Face for licensing terms