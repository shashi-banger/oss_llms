#!/usr/bin/env python3

import openai

# Configure the OpenAI client to use the local LiteLLM endpoint
client = openai.OpenAI(
    api_key="sk-sb123",  # Use the LITELLM_MASTER_KEY from docker-compose
    base_url="http://localhost:4000/v1"  # LiteLLM endpoint
)

def get_embeddings(texts):
    """Get embeddings for a list of texts"""
    response = client.embeddings.create(
        model="qwen-embedding",  # Model name from litellm_config.yaml
        input=texts
    )
    return response

def main():
    # Example texts to embed
    texts = [
        "Hello, how are you?",
        "The weather is nice today.",
        "Machine learning is fascinating.",
        "Python is a great programming language."
    ]
    
    print("Getting embeddings for texts...")
    
    try:
        # Get embeddings
        response = get_embeddings(texts)
        
        print(f"Model used: {response.model}")
        print(f"Number of embeddings: {len(response.data)}")
        print(f"Embedding dimension: {len(response.data[0].embedding)}")
        
        # Print first few dimensions of each embedding
        for i, embedding_obj in enumerate(response.data):
            embedding = embedding_obj.embedding
            print(f"\nText {i+1}: '{texts[i]}'")
            print(f"Embedding (first 5 dims): {embedding[:5]}, {len(embedding)}")
            
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()