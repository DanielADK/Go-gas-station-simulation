# Použijte oficiální obraz Go jako základní obraz
FROM golang:latest

# Nastavte pracovní adresář v kontejneru
WORKDIR /app

# Zkopírujte zdrojový kód projektu do kontejneru
COPY . .

# Stáhněte závislosti projektu (pokud existují)
RUN go mod download

# Kompilujte aplikaci pro produkční nasazení
RUN go build -o main cmd/*

# Spustí aplikaci při spuštění kontejneru
CMD ["./main"]