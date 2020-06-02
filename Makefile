run:
	DATABASE_URL=$$(heroku config:get DATABASE_URL -a projectincubator-backend) \
	BUCKET_NAME=$$(heroku config:get BUCKET_ID -a projectincubator-backend) \
	GOOGLE_APPLICATION_CREDENTIALS='key.json' \
	PORT=8000 \
    go run .
