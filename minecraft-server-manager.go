package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
				Name:  "init",
				Usage: "create the server",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					var version, loader string
					if len(args) > 0 {
						version = args[0]
					}
					if len(args) > 1 {
						loader = args[1]
					}

					if version == "" {
						panic("Minecraft version is required! Usage: create-minecraft-server <version> <forge | fabric | vanilla>")
					}

					if loader == "" {
						panic("Server type is required! Usage: create-minecraft-server <version> <forge | fabric | vanilla>")
					}

					createServerAndStart(loader, version)
					return nil
				},
			},
			{
				Name:  "discard",
				Usage: "delete the server instance, data will be persisted",
				Action: func(c *cli.Context) error {
					discardServer()
					return nil
				},
			},
			{
				Name:  "start",
				Usage: "start the server",
				Action: func(c *cli.Context) error {
					startServer()
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
func startServer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	containerName := "mc"

	err = cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}

	println("Server started successfully!")
}

func stopServer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	containerName := "mc"

	timeout := int(time.Second * 10)
	err = cli.ContainerStop(ctx, containerName, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	println("Server stopped successfully!")
}

func createServerAndStart(loader string, version string) {
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
		Env:   []string{"EULA=TRUE", "TYPE=" + loader, "VERSION=" + version, "INIT_MEMORY=6G", "MAX_MEMORY=6G"},
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
	println("Server created and started successfully!")
}

func discardServer() {
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

	println("Server discarded successfully. Run 'minecraft-server-manager init' to create a new one")
}
