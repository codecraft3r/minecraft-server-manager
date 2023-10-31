# Minecraft Server Manager

The Minecraft Server Manager is a command-line interface (CLI) app built with Go that allows you to manage a single Docker container with a Minecraft server in it. This app provides various commands to create, start, stop, and interact with the Minecraft server.

## Prerequisites

Before using the Minecraft Server Manager, make sure you have the following prerequisites installed:

- Docker: The app relies on Docker to manage the Minecraft server container. If you don't have Docker you can download it here: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

## Installation

To install the Minecraft Server Manager, you can download the binary for your operating system from [Github Releases](https://github.com/codecraft3r/minecraft-server-manager/releases). Be sure to choose the appropriate version for your system.

After downloading the binary, place it in a directory on your system's PATH.

# Usage

## Data Directory

The Minecraft server data directory is located at `~/Desktop/ServerData`. This is where your world and mod files live.

## Commands:

### init

Create the server instance with the specified Minecraft version and server type.

Usage: `minecraft-server-manager init <version> <forge | fabric | vanilla>`

Example: `minecraft-server-manager init 1.16.5 forge`

The `init` command creates the Minecraft server container with the specified version and server type. It sets up the necessary configurations and downloads the required server files.

### discard

Delete the server instance. Data will not be deleted.

Usage: `minecraft-server-manager discard`

The `discard` command deletes the Minecraft server container, but preserves the server data. You may need to move or delete the server data if you want to change the world, version, or modloader.
You will need to run `init` again if you want to make a new server.

### start

Start the Minecraft server.

Usage: `minecraft-server-manager start`

The `start` command starts the Minecraft server container, if it exists. 

**NOTE: This command will fail if the container doesn't exist, or the server is already running.**

### stop

Stop the Minecraft server.

Usage: `minecraft-server-manager stop`

The `stop` command stops the Minecraft server container without deleting it. To restart the server use the `start` command

**NOTE: This command will fail if the container doesn't exist, or the server is already stopped.**

### console

Open an interactive console to interact with the Minecraft server.

Usage: `minecraft-server-manager console`

The `console` command opens an interactive console that allows you to send commands directly to the Minecraft server.

### console-oneshot

Execute a one-shot console command on the Minecraft server.

Usage: `minecraft-server-manager console-oneshot <command>`

Example: `minecraft-server-manager console-oneshot say Hello, world!`

The `console-oneshot` command allows you to execute a single console command on the Minecraft server without opening an interactive console.

If you encounter any problems, please feel free to create an issue on the GitHub repository.
