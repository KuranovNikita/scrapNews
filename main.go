package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"scrapNews/internal/client"
	"scrapNews/internal/config"
	"scrapNews/internal/models"
	"scrapNews/storage"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	godotenv.Load(".env")

	cp := os.Getenv("CONFIG_PATH")

	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("start config path", slog.String("CONFIG_PATH", cp))
	log.Info("start url shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	client := client.New(cfg.Telegram.TgBotHost, cfg.Telegram.Token, log)

	// _, err := sql.Open("postgres", cfg.DB.DBURL)
	// if err != nil {
	// 	log.Error("Can't connect to database:", err)
	// } else {
	// 	log.Info("PG works!!!!")
	// }

	// TEST queries
	storageEx, err := storage.New("postgres", cfg.DB.DBURL)
	if err != nil {
		log.Error("Can't connect to database:", err)
	} else {
		log.Info("PG works!!!!")
	}
	userParams := storage.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "TestName",
		Type:      models.Telegram,
	}

	userNew, err := storageEx.SaveUser(userParams)
	if errors.Is(err, storage.ErrEmptyUser) {
		log.Info("User already make")
	} else if err != nil {
		log.Error("Can't make :", err)
	} else {
		log.Info("SaveUser works!!!!")
		telegramUserParams := storage.CreateTelegramUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      "TestName",
			ChatID:    "123456",
			UserID:    userNew.ID,
			Active:    true,
		}
		_, err := storageEx.SaveTelegramUser(telegramUserParams)
		if errors.Is(err, storage.ErrEmptyUser) {
			log.Info("User already make")
		} else if err != nil {
			log.Error("Can't make SaveTelegramUser:", err)
		} else {
			log.Info("SaveUserTelegram works!!!!")
		}
	}

	telUser, err := storageEx.GetTelegramUserByChatId("123456")
	if err != nil {
		log.Error("Can't make :", err)
	}
	fmt.Println("!!!!!!!!!!NaME IS %s", telUser.Name)

	UpdateTelegramUserActiveParams := storage.UpdateTelegramUserActiveParams{Active: false, ChatID: "123456"}
	err = storageEx.UpdateTelegramUserActive(UpdateTelegramUserActiveParams)
	if err != nil {
		log.Error("Can't make UpdateTelegramUserActive:", err)
	} else {
		log.Info("UpdateTelegramUserActive WORK")
	}

	createSiteParseParams := storage.CreateSiteParseParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		Name:          "SiteParse",
		UrlSite:       "http:site",
		Type:          "type",
		LastFetchedAt: time.Now().UTC(),
	}

	SiteParse, err := storageEx.CreateSiteParse(createSiteParseParams)
	if err != nil {
		log.Error("Can't make CreateSiteParse:", err)
	} else {
		log.Info("CreateSiteParse WORK")
	}

	_, err = storageEx.GetAllSiteParses()
	if err != nil {
		log.Error("Can't make GetAllSiteParses:", err)
	} else {
		log.Info("GetAllSiteParses WORK")
	}

	_, err = storageEx.GetSiteParseById(SiteParse.ID)
	if err != nil {
		log.Error("Can't make GetSiteParseById:", err)
	} else {
		log.Info("GetSiteParseById WORK")
	}

	_, err = storageEx.GetSiteParseByName(SiteParse.Name)
	if err != nil {
		log.Error("Can't make GetSiteParseByName:", err)
	} else {
		log.Info("GetSiteParseByName WORK")
	}

	createSiteParseParamsFollows := storage.CreateSiteParseParamsFollows{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		UserID:      userNew.ID,
		SiteParseID: SiteParse.ID,
		Active:      false,
	}
	// fmt.Printf("%v", createSiteParseParamsFollows)
	SiteParseFollows, err := storageEx.CreateSiteParseFollows(createSiteParseParamsFollows)
	if err != nil {
		log.Error("Can't make CreateSiteParseFollows:", err)
	} else {
		log.Info("CreateSiteParseFollows WORK")
	}

	_, err = storageEx.GetSiteParseFollowsByUserID(SiteParseFollows.UserID)
	if err != nil {
		log.Error("Can't make GetSiteParseFollowsByUserID:", err)
	} else {
		log.Info("GetSiteParseFollowsByUserID WORK")
	}

	updateSiteParseActiveParamsFollows := storage.UpdateSiteParseActiveParamsFollows{
		Active:      true,
		UpdatedAt:   time.Now().UTC(),
		UserID:      userNew.ID,
		SiteParseID: SiteParse.ID,
	}
	err = storageEx.UpdateSiteParseFollowsActive(updateSiteParseActiveParamsFollows)
	if err != nil {
		log.Error("Can't make UpdateSiteParseFollowsActive:", err)
	} else {
		log.Info("UpdateSiteParseFollowsActive WORK")
	}

	// TEST queries  END

	// messenger := models.Messenger{Name: "telegram"}

	// sendingMessage := models.SendingMessage{
	// 	ID:        "1246725945",
	// 	Text:      "hello testing222",
	// 	Messenger: messenger,
	// }
	// client.Send(sendingMessage)
	webHookTelegram := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Error("Error reading request body: %v", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			log.Error("Error parsing JSON: %v", err)
			return
		}
		var messenger = models.Messenger{}
		messenger.Name = models.Telegram
		client.WebhookEvent(data, messenger)
	}
	// button1 := models.Button{Text: "Кнопка 1"}
	// button2 := models.Button{Text: "Кнопка 2"}
	// sendingMessage2 := models.SendingMessage{
	// 	ID:        "1246725945",
	// 	Text:      "hello testing buttons",
	// 	Messenger: messenger,
	// 	Buttons:   []models.Button{button1, button2},
	// }
	// testTelegram.Send(sendingMessage2)
	router := chi.NewRouter()

	router.Post("/tgwebhook", webHookTelegram)

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Info("stopping server")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// func webHookTelegram(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	defer r.Body.Close()

// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		log.Printf("Error reading request body: %v", err)
// 		return
// 	}

// 	var data map[string]interface{} // Используем map[string]interface{} для динамической структуры
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		log.Printf("Error parsing JSON: %v", err)
// 		return
// 	}
// 	fmt.Println(data) ///!!!!!!
// 	// Извлекаем текст сообщения
// 	message, ok := data["message"].(map[string]interface{})
// 	if !ok {
// 		http.Error(w, "Invalid message format", http.StatusBadRequest)
// 		log.Printf("Error extracting 'message' field")
// 		return
// 	}

// 	formattedJSON, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		log.Printf("Error formatting JSON: %v", err)
// 		return
// 	}

// 	fmt.Println(string(formattedJSON))

// 	text, ok := message["text"].(string)
// 	if !ok {
// 		http.Error(w, "Invalid message format", http.StatusBadRequest)
// 		log.Printf("Error extracting 'text' field")
// 		return
// 	}

// 	// Выводим декодированный текст
// 	fmt.Println("Текст сообщения:", text)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(body)

// }
