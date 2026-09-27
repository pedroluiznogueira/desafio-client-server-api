# Client-Server API

`server.go` expõe `GET /cotacao` na porta 8080: busca o câmbio USD-BRL na AwesomeAPI (timeout
200ms) e persiste em SQLite (timeout 10ms — propositalmente apertado; se estourar, o erro é
logado no servidor mas a cotação é respondida ao cliente normalmente).

`client.go` chama `/cotacao` (timeout 300ms), extrai o campo `bid` e grava em `cotacao.txt` no
formato `Dólar: {valor}`.

## Rodando

Em dois terminais:

```
go run server.go
```

```
go run client.go
```

`cotacao.txt` é criado no diretório atual. Se o timeout do cliente (300ms) ou o do servidor
(200ms/10ms) estourar, o erro é logado no console correspondente.
