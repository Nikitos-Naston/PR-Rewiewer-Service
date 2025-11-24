package main

import (
	"PRreviewService/config"
	"PRreviewService/internal/handlers"
	"PRreviewService/internal/messages"
	"PRreviewService/internal/repository"
	"PRreviewService/storage"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	messages.SendLogMessage(logrus.InfoLevel, "Configs load succesful", nil)

	db, err := storage.New(cfg)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "DataBase connection unsuccesfull", err)
		return
	}
	time.Sleep(5 * time.Second)
	err = initDB(db)
	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "DataBase init unsuccesfull", err)
		return
	}
	messages.SendLogMessage(logrus.ErrorLevel, "DataBase init succes", nil)
	defer db.Close()
	messages.SendLogMessage(logrus.InfoLevel, "Connected to DataBase succesfull", nil)

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPRRepository(db)

	h := handlers.NewHandler(teamRepo, userRepo, prRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /team/add", h.AddTeam)
	mux.HandleFunc("GET /team/get/{team_name}", h.GetTeam)
	mux.HandleFunc("POST /users/setIsActive", h.SetUserActive)
	mux.HandleFunc("GET /users/getReview/{user_id}", h.GetRewiesByUser)
	mux.HandleFunc("POST /pullRequest/create", h.CreatePR)
	mux.HandleFunc("POST /pullRequest/merge", h.MergePR)
	mux.HandleFunc("POST /pullRequest/reassign", h.RessignRewiewer)
	mux.HandleFunc("GET /status", h.Status)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
	messages.SendLogMessage(logrus.InfoLevel, "Starting server at PORT:"+cfg.Port, nil)

	if err := server.ListenAndServe(); err != nil {
		messages.SendLogMessage(logrus.InfoLevel, "Server start unsucces", err)
	}
	messages.SendLogMessage(logrus.InfoLevel, "Server Starts succesfull", nil)
}
