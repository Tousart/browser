package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	ailogic "github.com/tousart/browser/ai_agent/ai_logic"
	"github.com/tousart/browser/api"
	"github.com/tousart/browser/usecase/service"
)

const (
	servicePort     = 8080
	chromDriverPort = 9515
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to download env file: %v", err)
	}

	apiKey := os.Getenv("OPEN_API_KEY")
	if apiKey == "" {
		log.Fatal("api key is empty")
	}

	chromDriverPath := os.Getenv("CHROMDRIVER_PATH")
	yandexBrowserPath := os.Getenv("YANDEX_BROWSER_PATH")

	// Подключаемся к яндекс браузеру через selenium
	serviceDriver, err := selenium.NewChromeDriverService(chromDriverPath, chromDriverPort)
	if err != nil {
		log.Fatalf("failed to create chrom driver service: %v", err)
	}
	defer serviceDriver.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--ignore-certificate-errors",
			"--disable-gpu",
			"--no-sandbox",
			"--disable-extensions",
			"--disable-dev-shm-usage",
			"--disable-web-security",
			"--start-maximized",
			`--user-data-dir=D:\second_ya_profile`,
			"--profile-directory=Default",
		},
		Path: yandexBrowserPath,
	}
	caps.AddChrome(chromeCaps)

	browserDriverURL := fmt.Sprintf("http://localhost:%d/wd/hub", chromDriverPort)
	webDriver, err := selenium.NewRemote(caps, browserDriverURL)
	if err != nil {
		log.Fatalf("failed to remote web driver: %v", err)
	}
	defer webDriver.Quit()

	// Сервис с обработкой веб-страниц в браузере
	mailService := service.CreateMailService(webDriver)

	// ИИ-агент
	aiAgent := ailogic.CreateAIAgent(apiKey, mailService)

	// API
	aiApi := api.CreateOpenAIApi(aiAgent)

	r := chi.NewRouter()
	aiApi.WithOpenAIHandlers(r)

	// Запуск
	log.Printf("Server has been started on %d\n", servicePort)
	if err := CreateAndRunServer(r, servicePort); err != nil {
		log.Fatalf("failed to create and run server: %v", err)
	}
}

func CreateAndRunServer(router *chi.Mux, port int) error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	return httpServer.ListenAndServe()
}
