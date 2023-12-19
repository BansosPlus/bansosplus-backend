# Bansos-plus backend

Backend application for bansosplus app

## Description

Based on the project plan, we build our backend application using Go programming language. We also use several packages that are available on Go programming language in our backend application, the packages are Echo to build REST API, GORM as object relation mapping, and go-qrcode to generate QR code of accepted Bansos registrations. We also use several services that are available on GCP to deploy our application, we use a compute engine to deploy our backend application, cloud SQL to implement our MySQL databases, cloud storage to store our Bansos image. We also use docker to containerize our backend application so we can deploy our application easily, we also implement CI/CD to our backend application using github CI/CD. 

## Requirements

1. Go
2. GORM
3. Echo

## How to Build & Run?

1. Run go mod tidy command on this repo

   ```
   go mod tidy
   ```
2. Run go run server.go command on this repo

   ```
   go run server.go
   ```

## Documentation

https://api.postman.com/collections/21473149-1132e122-18f8-43b9-a359-4847ffca5aea?access_key=PMAT-01HFAWWJDSHBKMK78AVC1ZQ865

## Contributors

| ID          | Name                   |
| ----------- | ---------------------- |
| C002BSY3034 | Ghazian Tsabit Alkamil |
| C002BSY3035 | Willy Wilsen           |
