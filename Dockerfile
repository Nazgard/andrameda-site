# ============ Этап 1: Сборка приложения ============
FROM golang:1.25-alpine AS builder

# Устанавливаем необходимые инструменты
RUN apk add --no-cache git tzdata

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum (для кэширования зависимостей)
#COPY go.mod go.sum ./

# Загружаем зависимости
#RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем статический бинарник (без CGO, для максимальной совместимости с scratch)
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o site main.go

# ============ Этап 2: Финальный минимальный образ ============
FROM scratch

# Копируем только исполняемый файл из этапа сборки
COPY --from=builder /app/site /site
COPY --from=builder /app/template/* /template/
COPY --from=builder /app/static/* /static/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Указываем точку входа
ENTRYPOINT ["/site"]

# Порт, на котором работает приложение (опционально, для документации)
EXPOSE 8080