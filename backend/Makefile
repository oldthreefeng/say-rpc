REGISTRY?=registry-vpc.cn-shanghai.aliyuncs.com/louisehong/say
build:
	GOOS=linux go build -o app
	docker build -t $(REGISTRY) .
	rm -f app
push:
	docker push $(REGISTRY) 
