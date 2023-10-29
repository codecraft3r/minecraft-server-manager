package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Minecraft Server Manager",
		Usage: "managing minecraft a server, obvs",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "start the server",
				Action: func(c *cli.Context) error {
					version := c.Args().First() // Get the first argument as the version
					loader := c.Args().Get(1)   // Get the second argument as the loader
					startServer(version, loader)
					return nil
				},
			},
			{
				Name:  "stop",
				Usage: "stop the server",
				Action: func(c *cli.Context) error {
					stopServer()
					return nil
				},
			},
			{
				Name:  "console",
				Usage: "open interactive console",
				Action: func(c *cli.Context) error {
					rconCommand := exec.Command("docker", "exec", "-i", "mc", "rcon-cli")
					rconCommand.Stdin = os.Stdin
					rconCommand.Stdout = os.Stdout
					rconCommand.Stderr = os.Stderr
					err := rconCommand.Run()
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "console-oneshot",
				Usage: "execute oneshot console command",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					rconCommand := exec.Command("docker", append([]string{"exec", "mc", "rcon-cli"}, args...)...)
					rconCommand.Stdin = os.Stdin
					rconCommand.Stdout = os.Stdout
					rconCommand.Stderr = os.Stderr
					err := rconCommand.Run()
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startServer(loader string, version string) {
	// handle vars
	loader = strings.ToLower(loader)

	switch loader {
	case "fabric":
		loader = "FABRIC"
	case "forge":
		loader = "FORGE"
	default:
		loader = "VANILLA"
	}

	// init docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// pull image
	reader, err := cli.ImagePull(ctx, "docker.io/itzg/minecraft-server", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	// eval data dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dataDir, err := filepath.Abs(os.ExpandEnv(homeDir + "/Desktop/ServerData"))
	if err != nil {
		panic(err)
	}
	// create container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "docker.io/itzg/minecraft-server",
		Env:   []string{"EULA=TRUE", "TYPE=" + loader, "VERSION=" + version},
		Tty:   false,
	}, &container.HostConfig{
		Binds: []string{dataDir + ":/data"},
	}, nil, nil, "mc")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	println("Server started successfully!")
}

func stopServer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containerName := "mc"

	err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		panic(err)
	}

	println("Server stopped successfully!")
}
