# job-avito-assignment
Задание на стажировку в Авито

### Общая Информация
- Основные требования находятся по ссылке: https://github.com/avito-tech/job-backend-trainee-assignment/
- Используется: Docker, Make, Golang, PostgreSQL
- Внешние библиотеки: ZapLogger, gorrila/mux
- Сервер работает на 9000 порту
- Путь до конфига /internal/app/config.go
- Выполнил: Андронов Дмитрий

### How to Run | Build

**Run Mode**
- Для запуска используется технология Docker
- Склонировать репозиторий и перейти в корень
- sudo docker build -t dev/avito-job .
- sudo docker run -p 9000:9000 --name dev -t dev/avito-job


## API / примеры запросов

###Начисление  и списание средств** 
 Запрос:
 
 ```bash
 curl --header "Content-Type: application/json" \
   --request POST \
   --data '{"user_id": user_1, "sum" : sum_1}' \
   http://localhost:9000/balance/change
 ```
Если sum > 0 - операция зачисления, иначе списания. sum != 0

###Получение информации о балансе
 Запрос:
 
 ```bash
 curl --header "Content-Type: application/json" \
   --request GET \
   --data '{"user_id": user_1}' \
   http://localhost:9000/balance/get
 ```

### Перевод денег между пользователями 
 Запрос:
 
 ```bash
 curl --header "Content-Type: application/json" \
   --request GET \
   --data '{"sender": user_1, "receiver" : user_2, "sum": 1000}' \
   http://localhost:9000/balance/transaction
 ```
 Где sum > 0



