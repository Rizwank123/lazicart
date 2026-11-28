ROOT_DIR := $(CURDIR)
SCRIPTS_DIR := $(ROOT_DIR)/scripts

SERVICE ?=
ACTION ?=

validate:
	@if [ -z "$(SERVICE)" ]; then \
		echo "SERVICE is required"; \
		exit 1; \
	fi
	@if [ -z "$(ACTION)" ]; then \
		echo "ACTION is required (up/down)"; \
		exit 1; \
	fi

migrate: validate
	@echo "Migration"
	@echo "Service: $(SERVICE)"
	@echo "Action : $(ACTION)"
	@bash "$(SCRIPTS_DIR)/$(SERVICE)/migration.sh" "$(ACTION)"

run-all:
	@bash "$(SCRIPTS_DIR)/run-all-services.sh"

kafka-topics:
	@bash "$(SCRIPTS_DIR)/create-kafka-topics.sh"

help:
	@echo ""
	@echo "Usage:"
	@echo "  make migrate SERVICE=auth-service ACTION=up"
	@echo "  make migrate SERVICE=auth-service ACTION=down"
	@echo "  make kafka-topics"
	@echo "  make run-all"
	@echo ""
