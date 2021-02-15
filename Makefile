create-database: ## Apply script to create database, user and schema for the service.
	USERNAME=postgres
	DB_HOST=localhost
	DB_PORT=5432
	
	psql postgres://$(USERNAME):$(PASSWORD)@$(DB_HOST):$(DB_PORT) \
	--variable=user_var=$(DB_USERNAME) \
	--variable=password_var=$(DB_PASSWORD) \
	-f scripts/schema.sql