package hello

import (
	"encoding/json"
	"net/http"

	"github.com/wieceslaw/chat-go/internal/server/auth"
)

func Register(mux *http.ServeMux, middleware func(hanlder func(w http.ResponseWriter, r *http.Request)) http.Handler) {
	mux.Handle("GET /api/v1/hello", middleware(hello))
}

type HelloMessage struct {
	Message string
}

func hello(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := HelloMessage{Message: "Hello, " + user.Name + "!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
