package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type PhotoResponse struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}

type Photo struct {
	ID          string   `json:"id"`
	Images      ImageSet `json:"urls"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Likes       int      `json:"likes"`
	User        User     `json:"user"`
}

type User struct {
	Name string `json:"name"`
}

type ImageSet struct {
	Full string `json:"full"`
}

const (
	baseURL = "https://api.unsplash.com"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Не вдалось завантажити .env:", err)
	}

	apiKey := os.Getenv("AccessKey")
	if apiKey == "" {
		fmt.Println("помилка: ключ доступа не знайдений. перконайтесь, що в .env є стрічка AccessKey=ВАШ_КЛЮЧ")
		return
	}

	fmt.Print("Введите тему для поиска картинок: ")
	reader := bufio.NewScanner(os.Stdin)
	ok := reader.Scan()
	if !ok {
		fmt.Println("Помилка ввода")
		return
	}
	query := reader.Text()

	resp, err := search(query, apiKey)
	if err != nil {
		fmt.Println("помилка при запросі:", err)
		return
	}

	if len(resp.Results) == 0 {
		fmt.Println("зображення по запиту не найдені.")
		return
	}

	fmt.Printf("Найдено %d зображень по темі %s", len(resp.Results), query)

	fmt.Printf("Тема: %s\n", query)
	for idx, photo := range resp.Results {
		fmt.Printf("зображення  №%d\n", idx+1)
		fmt.Printf("Автор: %s\n", photo.User.Name)
		if photo.Description != "" {
			fmt.Printf("Назва: %s\n", photo.Description)
		} else {
			fmt.Printf("Назва: %s\n", photo.Slug)
		}
		fmt.Printf("Лайків: %d\n", photo.Likes)
		fmt.Printf("URL: %s\n", photo.Images.Full)
		fmt.Println("------------------------------------------------")
	}
}

func search(query string, apiKey string) (PhotoResponse, error) {
	url := fmt.Sprintf("%s/search/photos?query=%s&client_id=%s", baseURL, query, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return PhotoResponse{}, fmt.Errorf("не вдалось виконати запрос: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return PhotoResponse{}, fmt.Errorf("забагато запросів: %s", resp.Status)
	} else if resp.StatusCode != 200 {
		return PhotoResponse{}, fmt.Errorf("сервер повернув помилку: %s", resp.Status)
	}

	var r PhotoResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return PhotoResponse{}, fmt.Errorf("помилка декодувания JSON: %v", err)
	}

	return r, nil
}
