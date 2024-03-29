# Makefile Usage Guide

This Makefile provides commands for building, running, and managing a Go application and NGINX server on Windows.

## Commands

- `make`: Builds and runs the Go application by default.
- `make build`: Compiles the Go application into an executable.
- `make run`: Starts the Go application.
- `make stop`: Stops the Go application.
- `make start-nginx`: Starts NGINX with a specified configuration file.
- `make stop-nginx`: Stops NGINX.
- `make restart-nginx`: Restarts NGINX.
- `make clean`: Deletes the compiled executable.

## How to Modify for Use

1. **Specify Go Executable Name**: Replace `notify-api.exe` with the name of your Go executable in the `build`, `run`, and `stop` targets.

2. **Specify NGINX Configuration File**: Replace `path/to/nginx.conf` in the `start-nginx` target with the path to your NGINX configuration file.

3. **Customize Commands**: Adjust the commands in each target to match your specific setup. For example, you may need to modify file paths, add additional build steps, or include environment variables.

4. **Save as Makefile**: Save the content as `Makefile` in the same directory as your Go application and NGINX configuration file.

5. **Run Commands**: Open a command prompt or PowerShell in the directory containing the Makefile and run the desired commands using `make`. For example, `make build` to compile the Go application or `make start-nginx` to start NGINX.

6. **Adjust as Needed**: Feel free to modify the Makefile further to suit your project's requirements. You can add additional targets, include conditional logic, or integrate with other tools as needed.

### Note
The Makefile was written for use on the Windows platform, modifications not defined in this documentation may be required for use on UNIX type platforms.