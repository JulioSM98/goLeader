package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{Filters: filters.NewArgs(filters.Arg("role", "manager"))})
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		spec_node := node.Spec
		if node.ManagerStatus.Leader {
			spec_node.Labels["isLeader"] = "true"
			fmt.Println("Node Leader: ", node.Description.Hostname)
		} else {
			spec_node.Labels["isLeader"] = "false"
		}

		err := cli.NodeUpdate(context.Background(), node.ID, node.Version, spec_node)
		if err != nil {
			panic(err)
		}

	}

}
