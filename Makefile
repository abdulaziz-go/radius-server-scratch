docker-up:
	cd docker && docker compose up -d

docker-up-clean:
	cd docker && docker compose up -d --remove-orphans

docker-down:
	cd docker && docker compose down

docker-prune:
	docker system prune -a --volumes

run-tests:
	go test -v ./tests -run '^TestOrderedSuite$$'