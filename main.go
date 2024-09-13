package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"scrapNews/internal/client"
	"scrapNews/internal/config"
	"scrapNews/internal/models"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
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
