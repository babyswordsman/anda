/*

the backend server - 
* maintains the chat history for each user as long-term memory
* judges to search or not
* rewrites the conversational query
* retrieves from Web search APIs and/or self-hosted internal indexes (Vearch by default)
* reranks the retrieval result
* generates the answer by calling LLMs
* suggests related questions
* can be planning and multi-step reasoning

*/

package main
