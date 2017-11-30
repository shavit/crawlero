build:
	docker build -t itstommy/crawlero .

start_dev:
	docker run --rm \
		--env-file ${PWD}/.env \
		--net kirra_network \
		-v ${PWD}:/go/src/github.com/shavit/crawlero \
		-ti itstommy/crawlero
