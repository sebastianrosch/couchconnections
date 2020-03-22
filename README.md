# Couch Connections

Presentations in your living room

# Run the app

## Run the Angular app
```sh
cd public/app
npm install && npm start
```

## Run the backend
```sh
go run cmd/couchconnections-api/main.go
```

# Deploy the app

## Requirements
- Heroku CLI

## Deploy
1. `cd public/app && npm install && npm run build && cd -`
2. `heroku container:push web -a couchconnections`
3. `heroku container:release web -a couchconnections` 