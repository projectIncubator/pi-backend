run:
	DATABASE_URL=$$(heroku config:get DATABASE_URL -a projectincubator-backend) \
	GOOGLE_APPLICATION_CREDENTIALS='key.json' \
	PORT=8000 \
    go run .

down:
	docker-compose down

up:
	docker-compose up -d
