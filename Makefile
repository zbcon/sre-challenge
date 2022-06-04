deploy:
	./deploy.sh

build:
	docker build -t invoice-app invoice-app/
	docker build -t payment-provider payment-provider/