local:
	docker-compose -f docker-compose.local.yml up --build
local-down:
	docker-compose -f docker-compose.local.yml down
migrate_up:
	docker-compose -f docker-compose.local.yml exec api migrate -source file://migrations -database 'mysql://root:root@tcp(db:3306)/videos-app' up
test:
	go test -cover ./...