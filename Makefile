deploy:
	./deploy.sh

build:
	docker build -t invoice-app invoice-app/
	docker build -t payment-provider payment-provider/

invoice:
	curl -X GET `minikube service invoice-svc --url`/invoices

pay:
	curl -X POST `minikube service invoice-svc --url`/invoices/pay