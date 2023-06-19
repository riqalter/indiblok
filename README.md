# Apakah Di Blok?

> Cek apakah suatu situs di blok atau tidak di jaringan telkom.

## Instalasi

```bash
git clone https://github.com/riqalter/indiblok.git
```

## Cara make

1. run langsung

```bash
go run main.go -d <domain> <domain> <domain>
```

2. kalo mau di build dulu

windows
```powershell
set CGO_ENABLED=0 && go build -o dibloktelkom.exe main.go
```

linux
```bash
CGO_ENABLED=0 go build -o dibloktelkom main.go
```

3. run hasil build

windows
```powershell
.\dibloktelkom.exe -d <domain> <domain> <domain>
```

linux
```bash
./dibloktelkom -d <domain> <domain> <domain>
```

## Contoh

```bash
go run main.go -d google.com reddit.com github.com
```

