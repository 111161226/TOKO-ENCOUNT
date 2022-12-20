.PHONY: dev
dev:
	docker-compose up --build

.PHONY: down
down:
	docker-compose down

.PHONY: swagger
swagger:
	rm -rf docs/swagger-ui
	git clone --depth 1 https://github.com/swagger-api/swagger-ui docs/swagger-ui
	cd docs/swagger-ui && git apply ../swagger-ui-diff.patch
