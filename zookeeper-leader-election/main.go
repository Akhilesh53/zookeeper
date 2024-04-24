package main

import (
	"fmt"
	"time"
	"zk/constants"
	"zk/services"
	"zk/utils"
)

// on start of the api server, create all parent nodes, live nodes and election node
func main() {

	zookeeperService := services.NewZookeeperService([]string{"localhost:2181"})
	//defer zookeeperService.Close()

	clusterInfo := utils.GetClusterInfo()

	// create all parent nodes
	if err := zookeeperService.CreateAllParentNodes(); err != nil {
		fmt.Println("Error while creating all parent nodes: ", err)
		return
	}

	// add a node to all_nodes
	// we will create a persistent node and add it to all nodes
	if err := zookeeperService.AddToAllNode(constants.SAMPLE_NODE_1, constants.SAMPLE_NODE_1); err != nil {
		fmt.Println("Error while adding node to all nodes: ", err)
		return
	}

	// add all the nodes to cluster info
	clusterInfo.Reset()
	clusterInfo.SetAllNodes(zookeeperService.GetAllNodes())
	fmt.Println("All nodes: ", clusterInfo.GetAllNodes())

	// create ephermeral + sequential node for /election
	if err := zookeeperService.AddToElectionNodes(constants.SAMPLE_NODE_1, constants.SAMPLE_NODE_1); err != nil {
		fmt.Println("Error while creating election node: ", err)
		return
	}

	//todo: to check wther ephemeral node is created or not, we need to halt our program cz ephermeral node will be deleted once the connection is closed
	// so we will add time.Sleep(100 Seconds) to check the ephemeral node
	//time.Sleep(30 * time.Second)

	//if no node is attached to election, the very first node will the the leader
	// so we will add the first node to the leader
	childs, _, err := zookeeperService.Children(constants.ELECTION_NODES)
	if err != nil {
		fmt.Println("Error while getting children of election node: ", err)
		return
	}

	fmt.Println("Children of election node: ", childs)

	if len(childs) == 1 {
		if err := zookeeperService.SetMaster(constants.SAMPLE_NODE_1); err != nil {
			fmt.Println("Error while setting master: ", err)
			return
		}
	}

	// set the current node to live nodes
	if err := zookeeperService.AddToLiveNode(constants.SAMPLE_NODE_1, constants.SAMPLE_NODE_1); err != nil {
		fmt.Println("Error while adding to live node: ", err)
		return
	}

	clusterInfo.AddLiveNode(constants.SAMPLE_NODE_1)
	fmt.Println("Live nodes: ", clusterInfo.GetLiveNodes())

	// if number of nodes under live nodes is greater than 1, then we will check the leader
	if len(clusterInfo.GetLiveNodes()) > 1 {
		leader, err := zookeeperService.ElectLeader()
		if err != nil {
			fmt.Println("Error while getting leader: ", err)
			return
		}
		clusterInfo.SetMaster(leader)
		fmt.Println("Leader: ", leader)
	}

	time.Sleep(100 * time.Second)
}
