package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	baseURL      = flag.String("url", "http://localhost:8080", "Адрес API сервиса")
	teamsCount   = flag.Int("teams", 20, "Количество команд")
	usersPerTeam = flag.Int("users", 10, "Количество пользователей в команде")
	prsPerTeam   = flag.Int("prs", 5, "Количество Pull Request'ов на каждую команду")
)

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

type PR struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func main() {
	flag.Parse()
	log.Printf("Старт генерации на %s: %d команд, %d юзеров/ком, %d PR/ком",
		*baseURL, *teamsCount, *usersPerTeam, *prsPerTeam)

	client := &http.Client{Timeout: 10 * time.Second}

	for i := 1; i <= *teamsCount; i++ {
		teamName := fmt.Sprintf("Team_%d", i)
		members := make([]User, 0, *usersPerTeam)

		for j := 1; j <= *usersPerTeam; j++ {
			globalID := (i-1)*(*usersPerTeam) + j

			isActive := true
			if j%10 == 0 {
				isActive = false
			}

			members = append(members, User{
				UserID:   fmt.Sprintf("user_%d", globalID),
				UserName: fmt.Sprintf("User_%d", globalID),
				TeamName: teamName,
				IsActive: isActive,
			})
		}

		team := Team{TeamName: teamName, Members: members}

		if err := postJSON(client, "/team/add", team); err != nil {
			log.Printf("Ошибка команды %s: %v", teamName, err)
			continue
		}
		log.Printf("Team %s создана", teamName)

		for k := 1; k <= *prsPerTeam; k++ {

			author := members[(k-1)%len(members)]

			prID := fmt.Sprintf("pr_%d_%d", i, k)
			pr := PR{
				ID:       prID,
				Name:     fmt.Sprintf("Feature %d by %s", k, author.UserName),
				AuthorID: author.UserID,
			}

			if err := postJSON(client, "/pullRequest/create", pr); err != nil {
				log.Printf("Ошибка PR %s: %v", prID, err)
			} else {
				log.Printf("PR %s создан (автор: %s)", prID, author.UserID)
			}
		}
	}
	log.Println("Генерация данных полностью завершена!")
}

func postJSON(client *http.Client, endpoint string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(*baseURL+endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}
	return nil
}
