local:
	docker-compose -f docker-compose.local.yml up --build
local-down:
	docker-compose -f docker-compose.local.yml down