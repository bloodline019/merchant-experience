# merchant-experience
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)  
[![merchant-experience](https://github.com/bloodline019/merchant-experience/actions/workflows/go.yml/badge.svg)](https://github.com/bloodline019/merchant-experience/actions/workflows/go.yml)  
Тестовое задание:
Задача разработать сервис, через который продавцы смогут передавать нам свои товары пачками в формате excel (xlsx).
UI делать не нужно, достаточно только API.

Сервис принимает на вход ссылку на файл и id продавца, к чьему аккаунту будут привязаны загружаемые товары. Сервис читает файл и сохраняет, либо обновляет товары в БД. Обновление будет происходить, если пара (id продавца, offer_id) уже есть у нас в базе. В ответ на запрос выдаёт краткую статистику: количество созданных товаров, обновлённых, удалённых и количество
строк с ошибками (например цена отрицательная, либо вообще не число).

Для проверки работоспособности сервиса нужно так же реализовать метод, с помощью которого можно будет достать список товаров
из базы. Метод должен принимать на вход id продавца, offer_id, подстрока названия товара (по тексту "теле" находились и
"телефоны", и "телевизоры"). Ни один параметр не является обязательным, все указанные параметры применяются через логический оператор "AND".

В каждой строке скачанного файла будет содержаться отдельный товар. Колонки в файле и соответствующие значения полей
товара следующие:

- offer_id уникальный идентификатор товара в системе продавца
- name название товара
- price цена в рублях
- quantity количество товара на складе продавца
- available true/false, в случае false продавец хочет удалить товар из нашей базы

# Использование

API реализует следующие методы:
1. POST /upload - загрузка Excel файла на сервер и его последующая обработка
2. POST /getGoods - получение списка товаров из базы согласно передаваемым параметрам

Для загрузки Excel файла на сервер и получения краткой статистики работы необходимо передать POST запрос по методу /upload на сервер c данными в формате JSON (Url, seller_id).  
Формат запроса:
```
curl -X POST -H "Content-Type: application/json" -d '{"Url": "https://cdn.discordapp.com/attachments/1061032871785136158/1061033596057567394/goods_initial.xlsx","seller_id": "1"}' http://localhost:8080/upload
```
Для получения списка товаров из базы необходимо передать POST запрос по методу /getGoods на сервер c данными в формате JSON (поддерживается любая комбинация параметров).  
Формат запроса:
```
curl -X POST -H "Content-Type: application/json" -d '{"offer_id": "", "seller_id": "", "goodSubstring": "Pro" }' http://localhost:8080/getGoods
```
# Запуск приложения через Docker:
```
docker compose up --build
```

# TODO:  
1) Обработка потенциальных ошибок
2) Заменить sql-запросы на ORM :heavy_check_mark:
3) Тесты
4) Контейнеризация (Docker) :heavy_check_mark:
5) Многопоточность
