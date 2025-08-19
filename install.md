# Quick install and run guide

This is the procedure to install this project in a Linux (Ubuntu) system.

1. Install Go (GoLang)
```
sudo apt update
sudo apt install golang-go
```

Verify the installation by running:
```
go version
```

2. Set up the project
Clone the repository and download and sync dependencies
```
git clone https://github.com/josedgm/orders_go.git
cd orders_go
go mod tidy
```

3. Run the project
```
go run main.go  
```
