up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f db

psql:
	docker compose exec db psql -U postgres -d ecommerce