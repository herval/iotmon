
build:
	docker build . -t iot:latest

logs:
	flyctl logs

deploy:
	flyctl launch

secrets:
	cat .env | sed -e '/^#/d' | xargs flyctl secrets set