package main

import (
	"PRreviewService/config"
	"PRreviewService/internal/handlers"
	"PRreviewService/internal/messages"
	"PRreviewService/internal/repository"
	"PRreviewService/storage"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	messages.SendLogMessage("Configs load succesful", nil)

	db, err := storage.New(cfg)
	if err != nil {
		messages.SendLogMessage("DataBase connection unsuccesfull", err)
		return
	}

	defer db.Close()
	messages.SendLogMessage("Connected to DataBase succesfull", nil)

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPRRepository(db)
	statRepo := repository.NewStatRepository(db)

	h := handlers.NewHandler(teamRepo, userRepo, prRepo, statRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /team/add", h.AddTeam)
	mux.HandleFunc("GET /team/get", h.GetTeam)
	mux.HandleFunc("POST /users/setIsActive", h.SetUserActive)
	mux.HandleFunc("GET /users/getReview", h.GetRewiesByUser)
	mux.HandleFunc("POST /pullRequest/create", h.CreatePR)
	mux.HandleFunc("POST /pullRequest/merge", h.MergePR)
	mux.HandleFunc("POST /pullRequest/reassign", h.RessignRewiewer)
	mux.HandleFunc("GET /status", h.Status)
	mux.HandleFunc("GET /stats", h.ShowStat)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
	messages.SendLogMessage("Starting server at PORT:"+cfg.Port, nil)

	if err := server.ListenAndServe(); err != nil {
		messages.SendLogMessage("Server start unsucces", err)
	}
	messages.SendLogMessage("Server Starts succesfull", nil)
}
