from langchain_community.llms import Ollama
from langchain_community.document_loaders import WebBaseLoader
from langchain.chains.summarize import load_summarize_chain

loader = WebBaseLoader("https://ollama.com/blog/run-llama2-uncensored-locally")
docs = loader.load()
print(f"docs: {docs}\n")

llm = Ollama(
    model="llama2:latest", verbose=True, timeout=10
)  # timeout parameter is seconds.
print(f"llm: {llm}\n")

chain = load_summarize_chain(llm, chain_type="stuff")
print(f"chain: {chain}\n")

result = chain.invoke(docs)
print(result)
