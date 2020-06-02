run:
	DATABASE_URL=$$(heroku config:get DATABASE_URL -a projectincubator-backend) \
	PORT=8000 \
	bash run.sh \
    go run .

