package main

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	defer i_a_Leader(cli)

	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{Filters: filters.NewArgs(filters.Arg("role", "manager"))})
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		spec_node := node.Spec
		if node.ManagerStatus.Leader {
			spec_node.Labels["isLeader"] = "true"
		} else {
			spec_node.Labels["isLeader"] = "false"
		}

		err := cli.NodeUpdate(context.Background(), node.ID, node.Version, spec_node)
		if err != nil {
			panic(err)
		}
	}

}

func i_a_Leader(client *client.Client) {
	info, err := client.Info(context.Background())
	if err != nil {
		panic(err)
	}
	nodeId := info.Swarm.NodeID
	node, _, err := client.NodeInspectWithRaw(context.Background(), nodeId)
	if err != nil {
		panic(err)
	}

	if node.ManagerStatus.Leader {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
