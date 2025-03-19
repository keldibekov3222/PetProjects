package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tebeka/selenium"
)

const (
	chromeDriverPath = "C:/chromedriver/chromedriver.exe" // Укажите правильный путь
	seleniumPort     = 4444
	baseURL          = "https://kuper.ru"
)

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Запуск Selenium WebDriver
	service, err := selenium.NewChromeDriverService(chromeDriverPath, seleniumPort)
	if err != nil {
		log.Fatalf("Ошибка при запуске ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Настройки для браузера
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", seleniumPort))
	if err != nil {
		log.Fatalf("Ошибка при подключении к WebDriver: %v", err)
	}
	defer wd.Quit()

	// Открываем страницу категорий
	categoriesURL := baseURL + "/categories"
	if err := wd.Get(categoriesURL); err != nil {
		log.Fatalf("Ошибка при открытии страницы категорий: %v", err)
	}

	// Ждем загрузки страницы
	time.Sleep(5 * time.Second)

	// Получаем все ссылки на категории
	categoryLinks, err := wd.FindElements(selenium.ByCSSSelector, "a[href^='/categories/']")
	if err != nil {
		log.Fatalf("Ошибка при поиске ссылок на категории: %v", err)
	}

	// Собираем URL категорий
	var categoryURLs []string
	for _, category := range categoryLinks {
		categoryURL, err := category.GetAttribute("href")
		if err != nil {
			log.Printf("Ошибка при получении ссылки на категорию: %v", err)
			continue
		}
		viewCategory := baseURL + categoryURL
		categoryURLs = append(categoryURLs, viewCategory)

		fmt.Printf("\nКатегория: %s\n", viewCategory)
	}

	// Выбираем две случайные категории
}
