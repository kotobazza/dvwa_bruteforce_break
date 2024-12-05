package taskSolution

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

func getCSRFToken(client *http.Client, loginURL string) (string, error) {
	resp, err := client.Get(loginURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get login page: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var csrfToken string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == "user_token" {
					for _, attr := range n.Attr {
						if attr.Key == "value" {
							csrfToken = attr.Val
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if csrfToken == "" {
		return "", fmt.Errorf("CSRF token not found")
	}

	fmt.Println("CSRF: ", csrfToken)
	return csrfToken, nil
}

func login(client *http.Client, loginURL, username, password, csrfToken string) (string, error) {
	loginData := url.Values{}
	loginData.Set("username", username)
	loginData.Set("password", password)
	loginData.Set("Login", "Login")
	loginData.Set("user_token", csrfToken)

	req, err := http.NewRequest("POST", loginURL, bytes.NewBufferString(loginData.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error during login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		if bytes.Contains(body, []byte("Welcome")) {
			fmt.Println("Вход выполнен успешно!")
			//fmt.Println(string(body))

			for _, cookie := range client.Jar.Cookies(&url.URL{Scheme: "http", Host: "192.168.55.36"}) {
				fmt.Println("Cookie: ", cookie.Name, "Value: ", cookie.Value)

				if cookie.Name == "PHPSESSID" {
					return cookie.Value, nil
				}
			}
		} else {
			return "", fmt.Errorf("login failed: unexpected response body")
		}
	} else {
		return "", fmt.Errorf("login failed, status: %s", resp.Status)
	}

	return "", nil
}

func GetSessionToken(loginURL string, username string, password string) string {

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}

	csrfToken, err := getCSRFToken(client, loginURL)
	if err != nil {
		fmt.Println("Error getting CSRF token:", err)
		return ""
	}

	phpsessid, err := login(client, loginURL, username, password, csrfToken)
	if err != nil {
		fmt.Println("Error during login:", err)
		return ""
	}

	fmt.Println("PHPSESSID:", phpsessid)

	return phpsessid

}
