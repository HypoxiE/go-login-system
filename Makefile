BINARY_NAME = gologin

INSTALL_DIR = /bin

build:
	go build -o $(BINARY_NAME) ./cmd/go-login

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"


clean:
	rm -f $(BINARY_NAME)
