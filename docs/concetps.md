# Overall design

The key elements of conversational search: question understanding, conversation flow, and answer generation. 

Laiwen is implemented with LLMs, chat history, and RAG.

## chat history

a CRUD API managing resources: memories, question/answer messages, summaries, et al.

## retrieval augmented generation. 

search from the index (general or vertical) and chat history, and send the result as context to the LLM. 

