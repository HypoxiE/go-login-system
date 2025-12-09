BINARY_NAME = gologin

INSTALL_DIR = /bin

build:
	go build -o $(BINARY_NAME) ./cmd/go-login

install: build
	mkdir -p $(INSTALL_DIR)
	mv $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"

test: build
	-@./gologin
	@$(MAKE) clean > /dev/null 2>&1

clean:
	rm -f $(BINARY_NAME)
