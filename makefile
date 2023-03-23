docker.build:
	docker build -f Dockerfile -t charlesbarkles/ipblockerservice . 

docker.push: docker.build
	docker push charlesbarkles/ipblockerservice:latest

docker.run: docker.build
	docker run -p 3333:3000 -p 3334:8080 charlesbarkles/ipblockerservice

kubernetes:
	kubectl create -f deployment.yaml
	kubectl expose deployment my-ipblockerservice --type=NodePort --name=ipblockerservice-service
	kubectl get svc

kubernetes.delete:
	kubectl delete service ipblockerservice-service
	kubectl delete deployment my-ipblockerservice

minikube.url:
	minikube service ipblockerservice-service --url

grpc.client:
	go run apis/grpc/client/client.go --addr "127.0.0.1:64159"

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative apis/grpc/proto/IPBlockerService.proto