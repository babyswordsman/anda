/*

the backend server - 
* maintains the chat history for each user as long-term memory
* rewrites the conversational query
* retrieves from Web search APIs and/or self-hosted internal indexes (Vearch by default)
* reranks the retrieval result
* generats the answer by calling LLMs
* suggests related questions

*/

package main
