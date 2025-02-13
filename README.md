# Avito Shop

Проект представляет собой сервис внутреннего магазин мерча, где сотрудники могут приобретать товары за монеты (coin). Каждому новому сотруднику выделяется 1000 монет, которые можно использовать для покупки товаров. Кроме того, монеты можно передавать другим сотрудникам в знак благодарности или как подарок.

*Стек проекта: Go, PostgreSQL*

## Запуск проекта

Запуск приложения: 

`docker-compose up --build`

Запуск миграций: 

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

## Проблемы с которыми я столкнулся при разработке проекта

1) **Миграция базы данных** <br>
    В техническом задании не указано, необходимо ли версионировать миграцию. В docker-compose.yaml *(представлен в дополнение к техническому заданию)* путь к миграции указан: `./migrations/init.sql`, что указывает на то, что миграция выполняется одним файлом init.sql. Несмотря на это при разработке проекта я использовал пакет для реализации миграций. **Я понимаю, что в таком проекте это избыточно, но проект я проектировал не как учебный, а как полноценный микросервис с возможностью дальнейшего развития функционала. В таком случае считаю необходимым полноценную реализацию миграции.**

## Список использованных пакетов

Пакеты, которые также использовались в проекте, но не упоминались выше:

- [github.com/jackc/pgx/v4/stdlib](https://github.com/jackc/pgx) - драйвер для работы с PostgreSQL
- ["github.com/jmoiron/sqlx"](https://github.com/jmoiron/sqlx) - расширение пакета database/sql