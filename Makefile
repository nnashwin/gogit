ifeq ($(BIN_NAME),)
BIN_NAME=$(notdir $(CURDIR))
endif

all: build move

build:
	go build -o $(BIN_NAME)

move:
	if [ "$(GOPATH)" = "" ] ; \
	then \
		echo "Your GOPATH is not set.  You must set it to continue"; \
	else \
		mv $(BIN_NAME) $(GOPATH)/bin/; \
	fi;
