# Практическая работа 3
Программа для параллельного перебора пароля на странице DVWA/Brute Force. 
Запуск: 
```bash
go run main.go --ip <IP адрес DVWM> --user <username> --pass <password>
```
Аргументы:
+ --ip
  + IP-адрес DVWM (как Docker, так и виртуальная машина)
+ --user
  + Имя пользователя, которое используется для получения доступа к всем задачам внутри DVWM
+ --pass
  + Пароль для того же пользователя


+ В самом коде есть предложение о выполнении полного перебора. Однако такая стратегия достаточно неэффективна в процессе перебора кредов из-за выполнения сетевых функций
+ Вместо полного перебора предполагается использование популярных кредов
  + Списки `unix_users.txt` и `unix_passwords.txt` внутри [Metasploit Framework](https://github.com/rapid7/metasploit-framework) отлично помогли в переборе и нашли нужные креды



## Порядок выполнения запросов
1. Выполнение запроса GET для получения страницы входа в DVWA
2. Выполнение запроса POST для входа в систему и получения PHPSESSID
3. Выполнение множества запросов GET к странице `DVWA/vulnerabilites/brute/`
4. Проверка каждого ответа на отсутствие текста о неправильном пароле