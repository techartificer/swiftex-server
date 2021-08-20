# swiftex-server
SwiftEx is an Open Source Delivery Application

### Prerequisite
- Go
- MongoDB
- Redis
- Firebase Admin SDK Private Key

### Development
First we need to clone this repository from github, install dependency and create .env file
```sh
$ git clone https://github.com/techartificer/swiftex-server
$ go mod download
$ cp example.env .env
```

We are using **Firebase Phone Auth** to verify phone number, hence to setup this project you need **Firebase Admin SDK Private Key**. You can generate from [here](https://console.firebase.google.com/)

Change necessary environment variable with valid varible and start application by running command
```sh
$ make run
```

## Environment Variable

| Variable Name            | Value                            |
|--------------------------|------------------------------------------------------------|
| `SERVER_HOST`            | 0.0.0.0                                                    |
| `SERVER_PORT`            | 4141                                                       |
| `SERVER_BCRYPT_COST`     | 10                                                         |
| `SERVER_NAME`            | swiftex                                                    |
| `SERVER_ENV`             | development                                                |
| `JWT_SECRET`             | sdfwfc23435                                                |
| `JWT_TTL`                | 3000                                                       |
| `JWT_REFRESH_TTL`        | 30000                                                      |
| `FIREBASE`               | {"type":"service_account",...}                             |
| `REDIS_HOST`             | localhost:6379                                             |
| `REDIS_PASWORD`          | password                                                   |
| `MONGO_URL`              | mongodb+srv://user:password@ewas.zqlyd.mongodb.net/example |
| `MONGO_DB_NAME`          | example                                                    |