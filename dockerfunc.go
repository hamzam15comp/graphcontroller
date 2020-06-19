package main

import (
	"context"
	"fmt"
	"path"
	"strings"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)


var images []string
type Networks struct {
	networkid   string
	networkname string
}
var graphnet = Networks{}
type stateID struct {
	state string
	id    string
}
var conState map[string]stateID


func CreateContainer(img string) error {
	conid, err := getConID(img)
	if err != nil {
		return err
	}
	if conid != "" {
		return nil
	}
	ctx := context.Background()
	cli, cerr := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if cerr != nil {
		return fmt.Errorf(
			"Failed to connect to docker engine",
			cerr,
		)
	}
	var portbinding = nat.PortMap{}
	vname, ports := parseVertex(img)
	portn, _ := strconv.Atoi(ports)
	if vname == "vertex"{
		portbinding = nat.PortMap {
			nat.Port("7000/tcp"): []nat.PortBinding{{
				HostIP: "0.0.0.0",
				HostPort: strconv.Itoa(portn+7000),
			}},
		}
	} else if vname == "edge" {
		portbinding = nat.PortMap {
			nat.Port("15672/tcp"): []nat.PortBinding{{
				HostIP: "0.0.0.0",
				HostPort: strconv.Itoa(portn+15672),
			}},
			nat.Port("5672/tcp"): []nat.PortBinding{{
				HostIP: "0.0.0.0",
				HostPort: strconv.Itoa(portn+5672),
			}},
		}
	} else {
		portbinding = nat.PortMap {
			nat.Port("15672/tcp"): []nat.PortBinding{{
				HostIP: "0.0.0.0",
				HostPort: strconv.Itoa(portn+15672),
			}},
			nat.Port("5672/tcp"): []nat.PortBinding{{
				HostIP: "0.0.0.0",
				HostPort: strconv.Itoa(portn+5672),
			}},
			nat.Port("7000/tcp"): []nat.PortBinding{{
				HostIP: "0.0.0.0",
				HostPort: strconv.Itoa(portn+7000),
			}},
		}
	}

	resp, terr := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:    img,
			Hostname: img,
			ExposedPorts: nat.PortSet{
				nat.Port("7000/tcp"): {},
				nat.Port("15672/tcp"): {},
				nat.Port("5672/tcp"): {},
			},
		},
		&container.HostConfig{
			PortBindings: portbinding,
			RestartPolicy: container.RestartPolicy{
				Name: "on-failure",
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				graphnet.networkname: &network.EndpointSettings{
					NetworkID: graphnet.networkid,
				},
			},
		},
		img,
	)
	if terr != nil {
		return fmt.Errorf(
			"Failed to create Container",
			terr,
		)
	}

	serr := cli.ContainerStart(
		ctx,
		resp.ID,
		types.ContainerStartOptions{},
	)
	if serr != nil {
		return fmt.Errorf(
			"Failed to start Container",
			serr,
		)
	}
	fmt.Println(img, "created!")
	return nil
}


func CreateNetwork(name string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return fmt.Errorf(
			"Failed to connect to docker engine",
			err,
		)
	}
	netres, nerr := cli.NetworkCreate(
		ctx,
		name,
		types.NetworkCreate{
			CheckDuplicate: true,
		})
	if nerr != nil {
		return fmt.Errorf(
			"Failed to create network",
			nerr,
		)
	}
	graphnet.networkname = name
	graphnet.networkid = netres.ID
	return nil
}

// GetImages 
func GetImages(iname string) (bool, error){
	images = nil
	exist := false
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return false, fmt.Errorf(
			"Failed to connect to docker engine",
			err,
		)
	}

	imagesStruct, err := cli.ImageList(
		ctx,
		types.ImageListOptions{},
	)
	if err != nil {
		return false, fmt.Errorf(
			"Failed to fetch Image list",
			err,
		)
	}
	for _, image := range imagesStruct {
		if len(image.RepoTags) == 0 {
			continue
		}
		imagename := strings.Split(image.RepoTags[0], ":")[0]
		if imagename == iname {
			exist = true
		}
		images = append(images, imagename)
	}
	return exist, nil
}

// Build the image using the native docker api
func BuildImage(name string) error {
	exist, err := GetImages(name)
	if err != nil {
		return fmt.Errorf(
			"Failed to fetch Image list",
			err,
		)
	}
	if exist {
		return nil
	}
	pathname := pwd + name
	ctx := context.Background()
	cli, cerr := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if cerr != nil {
		return fmt.Errorf(
			"failed to build image %s",
			name,
			cerr,
		)
	}
	options := types.ImageBuildOptions{
		SuppressOutput: true,
		Remove:         true,
		PullParent:	true,
		Tags:       []string{(path.Base(pathname) + ":latest")},
		Dockerfile: "Dockerfile",
	}

	buildCtx, terr := archive.TarWithOptions(
		pathname,
		&archive.TarOptions{},
	)
	if terr != nil {
		return fmt.Errorf("failed to build image %s", pathname, terr)
	}
	defer buildCtx.Close()
	buildResponse, berr := cli.ImageBuild(
		ctx,
		buildCtx,
		options,
	)
	if berr != nil {
		return fmt.Errorf(
			"failed to build image %s",
			name,
			berr,
		)
	}
	defer buildResponse.Body.Close()
	return nil
}

func removeContainer(vname string) (error) {
        ctx := context.Background()
        cli, err := client.NewClientWithOpts(
                client.FromEnv,
                client.WithAPIVersionNegotiation(),
        )
        if err != nil {
                return fmt.Errorf("Failed to connect to docker engine", err)
        }
        conid, err := getConID(vname)
        if err != nil {
                return fmt.Errorf("Failed to fetch Container list", err)
        }
	if conid == "" {
		return nil
	}
        err = cli.ContainerStop(ctx, conid, nil)
        if err != nil {
                return fmt.Errorf("Failed to stop container", err)
        }
        err = cli.ContainerRemove(ctx,
                conid,
                types.ContainerRemoveOptions{
                        Force: true,
                },
        )
        if err != nil {
                return fmt.Errorf("Failed to remove container", err)
        }
	fmt.Println(vname, "removed.")
	return nil
}

func getConID(conName string) (string, error) {
        var delid string = ""
        ctx := context.Background()
        cli, err := client.NewClientWithOpts(
                client.FromEnv,
                client.WithAPIVersionNegotiation(),
        )
        if err != nil {
                return delid, fmt.Errorf("Failed to connect to docker engine", err)
        }
        conlist, cerr := cli.ContainerList(
                ctx,
                types.ContainerListOptions{
                        All: false,
                },
        )
        if cerr != nil {
                return delid, fmt.Errorf("Failed to get Container list", err)
        }
        for _, c := range(conlist) {
                if c.Image == conName {
                        delid = c.ID
                        break
                }
        }
        return delid, nil
}

/*
func main(){
	fmt.Println("hello")
}
*/
