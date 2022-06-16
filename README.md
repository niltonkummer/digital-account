# Internal Transfer

Transferência entre contas

## Features

- Criação de conta
- Autenticação
- Transferência entre contas

## Pré-requisitos
- Docker
- Docker-compose
- heroku
- jq (opcional)


## Rodando local

```shell
docker-compose up --build
```

## Web API

#### Heroku: `https://digital-banking.herokuapp.com`
#### Local: `http://localhost:8080`

### Auth
`/login`

`Request`
```shell
curl /login --data '{"cpf":"83948554072","secret":"12345"}' -s -o /dev/stdout
```
`Response`
```json
{"code":200,"expire":"2022-06-16T21:33:56Z","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTU0MTUyMzYsImp0aSI6MSwib3JpZ19pYXQiOjE2NTU0MTE2MzZ9.kPuUlZNo65qgcIrcZ1Q2PpxeGF4LoNQut1fXTyiG2Tk"}%      
```
Modo alternativo em variável
```shell
TOKEN=$(curl https://digital-banking.herokuapp.com/login --data '{"cpf":"83948554072","secret":"12345"}' -s -o /dev/stdout | jq -r ".token")
echo $TOKEN
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTU0MTUyNjUsImp0aSI6MSwib3JpZ19pYXQiOjE2NTU0MTE2NjV9.No_owV_smQvvaOuQQTx8L54WzKcKfFjNq1ICl0HkVi4
```

### Accounts
`GET /api/accounts`

`Request`
```shell
curl https://digital-banking.herokuapp.com/api/accounts -H "Authorization: Bearer ${TOKEN}" -s -o /dev/stdout 
```

`Response`
```json
[
  {"id":1,"balance":"-100","type":0,"user":{"id":1,"name":"Nilton","cpf":"839******72"}},
  {"id":4,"balance":"100","type":0,"user":{"id":2,"name":"Carlos","cpf":"123******10"}}
]
```

`GET /api/accounts/:account_id/balance`

`Request`
```shell
curl https://digital-banking.herokuapp.com/api/accounts/1/balance -H "Authorization: Bearer ${TOKEN}" -s -o /dev/stdout 
```

`Response`
```json
{"balance":"-91821673.14"}% 
```

`POST /api/accounts`

`Request`
```shell
curl -POST https://digital-banking.herokuapp.com/api/accounts -X POST --data '{"cpf":"12345678910","secret":"12345","name":"Carlos"}' -s -o /dev/stdout  
```

`Response`
```json
{"balance":"-91821673.14"}% 
```

### Transfers

`POST /api/transfers`
```shell
--data '{"account_destination_id":2, "amount":102.10}'
```



## Author
- Nilton Henrique Kummer - [niltonkummer](https://github.com/niltonkummer)

## License
Copyright © 2020 [niltonkummer](https://github.com/niltonkummer).
This project is [MIT](LICENSE) licensed.