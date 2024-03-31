.PHONY: all clean

all: build run

build:
	@echo "Building Go application..."
	go build -o notify-api.exe

run:
	@echo "Starting Go application..."
	@echo "Use Ctrl+C to gracefully shutdown..."
	.\notify-api.exe &

stop:
	@echo "Forcefully stopping Go application..."
	taskkill /F /IM notify-api.exe

start-nginx:
	@echo "Starting NGINX..."
	nginx -c path/to/nginx.conf

stop-nginx:
	@echo "Stopping NGINX..."
	taskkill /F /IM nginx.exe

restart-nginx:
	@echo "Restarting NGINX..."
	nginx -s reload

clean:
	@echo "Cleaning up..."
	del .\notify-api.exe