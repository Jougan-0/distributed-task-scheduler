.PHONY: build run k8s-deploy helm-install monitor release clean

build:
	docker-compose build

run:
	docker-compose up -d

k8s-deploy:
	kubectl apply -f k8s/manifest/

helm-install:
	helm install task-scheduler k8s/helm/

monitor:
	kubectl port-forward svc/grafana 3000:3000 &
	kubectl port-forward svc/prometheus 9090:9090 &
	kubectl port-forward svc/elasticsearch 9200:9200 &
	echo "Grafana: http://localhost:3000"
	echo "Prometheus: http://localhost:9090"
	echo "Elasticsearch: http://localhost:9200"

release:
	docker build -t myrepo/task-scheduler:latest .
	docker push myrepo/task-scheduler:latest
	cd k8s/helm && helm package .
	cd k8s/helm && helm push task-scheduler myrepo/helm-charts
	echo "Docker image and Helm chart pushed successfully."

clean:
	docker-compose down -v
	kubectl delete -f k8s/manifest/
	helm uninstall task-scheduler