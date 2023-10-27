# docker-registry-auth-proxy

Прокси сервер для docker registry. Испольует basic auth в качестве способа авторизации. 
Данные о пользователях хранятся в конфиг файле по пути `/app/.authproxy/config.json`.
Путь до конфиг файла можно переопределить с помощью переменной окружения `AUTH-PROXY_ACCESS_CONFIG_PATH`.
Структура конфиг файла такая
```json
[
  {
    "username": "admin",
    "password": "auth-proxy-pass"
  }
]
```

### Build

```bash
brewkit build
docker build . -f docker/Dockerfile -t veresnikov/docker-registry-auth-proxy:dev
```
