# Makefile for portal-data project
# Provides convenient commands for development and CI

.PHONY: help install test lint format security pre-commit clean docker build

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Installation and setup
install: ## Install project dependencies
	pip install --upgrade pip
	pip install poetry
	poetry install

install-dev: ## Install development dependencies
	pip install --upgrade pip
	pip install poetry
	poetry install
	poetry install --with dev

uninstall:
	pip freeze | xargs pip uninstall -y

# Environment
env-example: ## Create .env from .env.example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ .env file created from .env.example"; \
		echo "⚠️  Please update .env with your actual values"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi

env-check: ## Check environment setup
	@echo "Checking environment..."
	@python --version
	@pip --version
	@poetry --version || echo "Poetry not installed"
	@pre-commit --version || echo "Pre-commit not installed"
	@echo "✅ Environment check complete!"

setup-pre-commit: ## Setup pre-commit hooks
	pip install pre-commit
	pre-commit install
	pre-commit install --hook-type pre-push

setup-gitlab-ci: ## Setup GitLab CI configuration
	chmod +x ./docs/setup-gitlab-ci.sh
	./docs/setup-gitlab-ci.sh

# Testing
test: ## Run all tests
	python -m pytest tests/ -v

test-coverage: ## Run tests with coverage
	python -m pytest tests/ --cov=app --cov-report=html --cov-report=term-missing

test-fast: ## Run tests without coverage (faster)
	python -m pytest tests/ -v --tb=short

test-organizations: ## Run organizations feature tests
	python -m pytest tests/features/organizations/ -v

test-auth: ## Run auth feature tests
	python -m pytest tests/features/auth/ -v

test-unit: ## Run unit tests only
	python -m pytest tests/ -v -m "unit"

test-integration: ## Run integration tests only
	python -m pytest tests/ -v -m "integration"

# Code quality
lint: ## Run linting checks
	flake8 app/ tests/ --max-line-length=119 --extend-ignore=E203,W503,E501
	mypy app/ --ignore-missing-imports --no-strict-optional

format: ## Format code with black and isort
	black app/ tests/ --line-length=119
	isort app/ tests/ --profile=black --line-length=119

format-check: ## Check code formatting without changing files
	black --check app/ tests/ --line-length=119
	isort --check-only app/ tests/ --profile=black --line-length=119

# Security
security: ## Run security checks
	bandit -r app/ -f json -o bandit-report.json
	bandit -r app/ -f txt

security-html: ## Run security checks with HTML report
	bandit -r app/ -f html -o bandit-report.html

# Pre-commit
pre-commit: ## Run pre-commit hooks on all files
	pre-commit run --all-files

pre-commit-update: ## Update pre-commit hooks to latest versions
	pre-commit autoupdate

pre-commit-clean: ## Clean pre-commit cache
	pre-commit clean

# Database
migrate: ## Run database migrations
	alembic upgrade head

migrate-create: ## Create new migration (usage: make migrate-create MESSAGE="your message")
	alembic revision --autogenerate -m "$(MESSAGE)"

migrate-check: ## Check migration status
	alembic check

migrate-reset: ## Reset database (drop and recreate)
	@echo "⚠️  This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		alembic downgrade base; \
		alembic upgrade head; \
		echo "✅ Database reset complete"; \
	else \
		echo "❌ Database reset cancelled"; \
	fi

# pg_search
pg-search-generate: ## Generate pg_search migration code from model metadata
	python scripts/generate_pg_search_migration.py

pg-search-verify: ## Verify pg_search model configurations
	python -c "from app.utils.pg_search_helper import verify_searchable_models; is_valid, errors = verify_searchable_models(); print('✅ All models valid!' if is_valid else '\n'.join(errors))"

# Docker
docker-build: ## Build Docker image
	docker build -t portal-data:latest .

docker-run: ## Run Docker container
	docker run -p 8000:8000 portal-data:latest

docker-test: ## Run tests in Docker container
	docker run --rm portal-data:latest python -m pytest tests/ -v

docker-up: ## Start all Docker services
	docker-compose up -d

docker-down: ## Stop all Docker services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-restart: ## Restart all Docker services
	docker-compose restart

# Development
dev: ## Start development server
	uvicorn app:app --reload --host 0.0.0.0 --port 8000

dev-docker: ## Start development with Docker Compose
	docker-compose up --build

# Production
prod: ## Start production server with Gunicorn
	gunicorn -c gunicorn.conf.py app:app

prod-daemon: ## Start production server in daemon mode
	gunicorn -c gunicorn.conf.py --daemon app:app

prod-stop: ## Stop production server
	pkill -TERM gunicorn

prod-reload: ## Reload production server (zero downtime)
	pkill -HUP gunicorn

prod-restart: ## Restart production server
	$(MAKE) prod-stop
	sleep 2
	$(MAKE) prod

prod-workers-increase: ## Increase number of workers
	pkill -TTIN gunicorn

prod-workers-decrease: ## Decrease number of workers
	pkill -TTOU gunicorn

setup: ## Complete project setup
	$(MAKE) install-dev
	$(MAKE) env-example
	$(MAKE) setup-pre-commit
	$(MAKE) migrate
	@echo "✅ Project setup complete!"

# CI/CD simulation
ci-lint: ## Simulate CI linting stage
	@echo "Running CI linting stage..."
	pip install black isort flake8 mypy
	black --check app/ tests/
	isort --check-only app/ tests/
	flake8 app/ tests/ --max-line-length=88
	mypy app/ --ignore-missing-imports

ci-test: ## Simulate CI testing stage
	@echo "Running CI testing stage..."
	pip install pytest pytest-asyncio pytest-cov
	python -m pytest tests/ -v --cov=app --cov-report=xml

ci-security: ## Simulate CI security stage
	@echo "Running CI security stage..."
	pip install bandit
	bandit -r app/ -f json -o bandit-report.json

ci-all: ## Simulate complete CI pipeline
	@echo "Running complete CI pipeline..."
	$(MAKE) ci-lint
	$(MAKE) ci-test
	$(MAKE) ci-security

# GitLab CI specific
gitlab-ci-test: ## Test GitLab CI configuration locally
	gitlab-runner exec docker pre-commit
	gitlab-runner exec docker test

# Cleanup
clean: ## Clean up generated files
	find . -type f -name "*.pyc" -delete
	find . -type d -name "__pycache__" -delete
	find . -type d -name "*.egg-info" -exec rm -rf {} +
	rm -rf .pytest_cache/
	rm -rf .coverage
	rm -rf htmlcov/
	rm -rf bandit-report.*

clean-docker: ## Clean up Docker resources
	docker system prune -f
	docker volume prune -f

clean-all: ## Clean up everything
	$(MAKE) clean
	$(MAKE) clean-docker
	rm -rf venv/
	rm -rf .env

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@echo "Documentation files:"
	@echo "- README.md"
	@echo "- docs/README.md"
	@echo "- tests/features/organizations/README.md"
	@echo "- tests/features/organizations/SUMMARY.md"

# Quality gates
quality-check: ## Run all quality checks
	@echo "Running quality checks..."
	$(MAKE) format-check
	$(MAKE) lint
	$(MAKE) test-fast
	$(MAKE) security
	@echo "✅ All quality checks passed!"

quality-fix: ## Fix code quality issues
	@echo "Fixing code quality issues..."
	$(MAKE) format
	$(MAKE) clean
	@echo "✅ Code quality issues fixed!"

# Release
release-check: ## Check if ready for release
	@echo "Checking release readiness..."
	$(MAKE) quality-check
	$(MAKE) test-coverage
	@echo "✅ Ready for release!"

# Help for specific targets
test-help: ## Show test-related commands
	@echo "Test commands:"
	@echo "  make test                 - Run all tests"
	@echo "  make test-organizations   - Run organizations tests"
	@echo "  make test-auth            - Run auth tests"
	@echo "  make test-unit            - Run unit tests"
	@echo "  make test-integration     - Run integration tests"
	@echo "  make test-coverage        - Run tests with coverage"
	@echo "  make test-fast            - Run tests without coverage"

lint-help: ## Show linting-related commands
	@echo "Linting commands:"
	@echo "  make lint          - Run linting checks"
	@echo "  make format        - Format code"
	@echo "  make format-check  - Check formatting"
	@echo "  make pre-commit    - Run pre-commit hooks"

ci-help: ## Show CI-related commands
	@echo "CI commands:"
	@echo "  make ci-lint       - Simulate CI linting"
	@echo "  make ci-test       - Simulate CI testing"
	@echo "  make ci-security   - Simulate CI security"
	@echo "  make ci-all        - Simulate complete CI pipeline"

docker-help: ## Show Docker-related commands
	@echo "Docker commands:"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-run      - Run Docker container"
	@echo "  make docker-up       - Start all services"
	@echo "  make docker-down     - Stop all services"
	@echo "  make docker-logs     - View logs"
	@echo "  make docker-restart  - Restart services"

db-help: ## Show database-related commands
	@echo "Database commands:"
	@echo "  make migrate              - Run migrations"
	@echo "  make migrate-create       - Create new migration"
	@echo "  make migrate-check        - Check migration status"
	@echo "  make migrate-reset        - Reset database (⚠️  destructive)"
	@echo "  make pg-search-generate   - Generate pg_search migration from models"
	@echo "  make pg-search-verify     - Verify pg_search configurations"

prod-help: ## Show production server commands
	@echo "Production server commands:"
	@echo "  make prod                    - Start production server with Gunicorn"
	@echo "  make prod-daemon             - Start production server in daemon mode"
	@echo "  make prod-stop               - Stop production server"
	@echo "  make prod-reload             - Reload server (zero downtime)"
	@echo "  make prod-restart            - Restart production server"
	@echo "  make prod-workers-increase   - Increase number of workers"
	@echo "  make prod-workers-decrease   - Decrease number of workers"
