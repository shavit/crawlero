build:
	docker build -t itstommy/crawlero .
	docker run --rm \
		-e GOOS=${GOOS} \
		-ti itstommy/crawlero \
			go build -o build/listen cmd/main.go

start_dev:
	docker run --rm \
		--env-file ${PWD}/.env \
		--name crawlero_1 \
		--net kirra_network \
		-v ${PWD}:/go/src/github.com/shavit/crawlero \
		-ti itstommy/crawlero bash
