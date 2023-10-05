ifeq (run,$(firstword $(MAKECMDGOALS)))

  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

  $(eval $(RUN_ARGS):;@:)

endif

build:
	./scripts/build.sh

build_tools:
	./scripts/build_tools.sh

scrape: build_tools
	./bin/map_scraper $(RUN_ARGS)
	./bin/init_folders

tidy:
	go mod tidy