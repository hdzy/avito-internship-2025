# Avito Shop

Проект представляет собой сервис внутреннего магазин мерча, где сотрудники могут приобретать товары за монеты (coin). Каждому новому сотруднику выделяется 1000 монет, которые можно использовать для покупки товаров. Кроме того, монеты можно передавать другим сотрудникам в знак благодарности или как подарок.

*Стек проекта: Go, PostgreSQL*

## Запуск проекта

Запуск приложения: 

`docker-compose up --build`

Запуск миграций (необходимо при изменении файлов миграции): 

`docker-compose run avito-shop-service /build -migrate`

## Структура проекта

Проект организован по принципам чистой архитектуры: 

- **/config** - конфигурационные файлы
- **/internal**
    - **/cmd** — точка входа в приложение
    - **/config** — конфигурация приложения
    - **/entity** — описание сущностей
    - **/handlers** — HTTP-обработчики
    - **/logger** — настройка логирования
    - **/migrations** — миграция базы данных
    - **/repository** — слой доступа к данным.
    - **/service** — бизнес-логика
- **/migrations** - файлы для миграции базы данных
- **/pkg** — вспомогательные библиотеки
- **/scripts** — make-файлы для запуска сервиса
- **/tests** — интеграционные и E2E-тесты


## Конфигурация

Пример файла конфигурации находится в каталоге config в корне проекта. Для корректной работы сервиса необходимо создать файл config.yaml в каталоге config в корне проекта и задать соответствующие настройки конфигурации.

Конфигурация загружается через [Viper](https://github.com/spf13/viper).

## Миграция базы данных

Для миграции базы данных используется [golang-migrate/migrate](https://github.com/golang-migrate/migrate).


Миграционные файлы хранятся в каталоге migrations и имеют формат:

- Таблица employees *содержит информацию о сотрудниках*:
  - <version\>_create_employees_table.up.sql
  - <version\>_create_employees_table.down.sql 
- Таблица transactions *содержит историю операций*:
  - <version\>_create_transactions_table.up.sql
  - <version\>_create_transactions_table.down.sql 
- Таблица transactions *содержит информацию о мерче*:
  - <version\>_create_merch_items_table.up.sql
  - <version\>_create_merch_items_table.down.sql


## Логирование 

Логирование реализовано через пакет slog. Вывод логов зависит от окружения. Соответствующие настройки конфигурации: `env`, `log_level`, `log_format`. 

## Тестирование

### Юнит-тесты:
- Покрытие: 71.9%
- Запуск: `go test -v ./...` (добавить флаг -cover для подсчета покрытия)

### Интеграционные тесты:

Запуск: `go test -v -tags=integration ./integration`

### Нагрузочные тесты:

Результат нагрузочного тестирования через Apache Bench эндпоинта /api/info: 

```
ab -n 10000 -c 1000 http://localhost:8080/api/info
This is ApacheBench, Version 2.3 <$Revision: 1923142 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /api/info
Document Length:        28 bytes

Concurrency Level:      1000
Time taken for tests:   0.528 seconds
Complete requests:      10000
Failed requests:        0
Non-2xx responses:      10000
Total transferred:      1880000 bytes
HTML transferred:       280000 bytes
Requests per second:    18944.88 [#/sec] (mean)
Time per request:       52.785 [ms] (mean)
Time per request:       0.053 [ms] (mean, across all concurrent requests)
Transfer rate:          3478.16 [Kbytes/sec] received

Connection Times (ms)
min  mean[+/-sd] median   max
Connect:        0    6   4.1      6      17
Processing:     6   45  38.6     31     183
Waiting:        0   42  38.7     29     179
Total:         15   50  39.9     36     193

Percentage of the requests served within a certain time (ms)
50%     36
66%     40
75%     45
80%     51
90%     99
95%    174
98%    181
99%    184
100%    193 (longest request)
```

## Проблемы с которыми я столкнулся при разработке проекта

1) **Миграция базы данных** <br>
    В техническом задании не указано, необходимо ли версионировать миграцию. В docker-compose.yaml *(представлен в дополнение к техническому заданию)* путь к миграции указан: `./migrations/init.sql`, что указывает на то, что миграция выполняется одним файлом init.sql. Несмотря на это при разработке проекта я использовал пакет для реализации миграций. **Я понимаю, что в таком проекте это избыточно, но проект я проектировал не как учебный, а как полноценный микросервис с возможностью дальнейшего развития функционала. В таком случае считаю необходимым полноценную реализацию миграции.**
2) **highload-реализация базы данных** <br>
  В техническом задании не указана необходимость шардирования, репликации или pgBouncer (например). В рамках разработки посчитал это излишнем и не показательным для моего уровня, знаний и опыта разработки на Go. Я знаю об этих инструментах (/паттернах) и в highload-микросервисе использование подобного действительно необходимо, чтобы избежать недоступности, задержки и потери данных.

## Список использованных пакетов

Пакеты, которые также использовались в проекте, но не упоминались выше:

- [github.com/jackc/pgx/v4/stdlib](https://github.com/jackc/pgx) - драйвер для работы с PostgreSQL
- [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) - расширение пакета database/sql
- [github.com/golang-jwt/jwt/v4](https://github.com/golang-jwt/jwt/v4) - пакет для работы с JWT-токенами
- [github.com/stretchr/testify](https://github.com/stretchr/testify) - пакет для тестирования
- [github.com/stretchr/testify/mock](https://github.com/stretchr/testify/mock) - пакет для моков
- [github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) - пакет для моков