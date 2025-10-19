package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	aiagent "github.com/tousart/browser/ai_agent"
	"github.com/tousart/browser/models"
)

type OpenAIApi struct {
	AIAgent aiagent.AIAgent
}

func CreateOpenAIApi(aiAgent aiagent.AIAgent) *OpenAIApi {
	return &OpenAIApi{
		AIAgent: aiAgent,
	}
}

func (op *OpenAIApi) postUsersRequestHandler(w http.ResponseWriter, r *http.Request) {
	var request models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("api openai postUsersRequestHandler error: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	answer, err := op.AIAgent.DoPromptWithContent(request.Message)
	if err != nil {
		log.Printf("api openai postUsersRequestHandler error: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(answer)
	if err != nil {
		log.Printf("api openai postUsersRequestHandler error: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (op *OpenAIApi) WithOpenAIHandlers(r *chi.Mux) {
	r.Post("/ask", op.postUsersRequestHandler)
}
