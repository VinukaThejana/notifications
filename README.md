# Notifications

A simple notification implementation built with Kafka and Go

## How to run it

### Prerequisites
- Git

- Docker


### Steps
- Clone the Repo
```shell
git clone git@github.com:VinukaThejana/notifications.git && cd notifications
```

- Run the Kafka Broker and the Kafka management UI with docker.
```shell
cd init/ && docker compose up -d
```

- Run the Producer to produce new notifications.
```shell
go run cmd/producer/producer.go
```

- Run the consumer to consume notifications in real time as the user submits them.
```shell
go run cmd/consumer/consumer.go
```
