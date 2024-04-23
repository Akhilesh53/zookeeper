package utils

import "sync"

var clusterInfo *ClusterInfo

type ClusterInfo struct {
	liveNodes []string
	allNodes  []string
	master    string
}

var clusterInfoOnce sync.Once

func GetClusterInfo() *ClusterInfo {
	clusterInfoOnce.Do(func() {
		clusterInfo = &ClusterInfo{
			liveNodes: make([]string, 0),
			allNodes:  make([]string, 0),
			master:    "",
		}
	})
	return clusterInfo
}

func (c *ClusterInfo) GetLiveNodes() []string {
	return c.liveNodes
}

func (c *ClusterInfo) SetLiveNodes(liveNodes []string) {
	c.liveNodes = liveNodes
}

func (c *ClusterInfo) GetAllNodes() []string {
	return c.allNodes
}

func (c *ClusterInfo) SetAllNodes(allNodes []string) {
	c.allNodes = allNodes
}

func (c *ClusterInfo) GetMaster() string {
	return c.master
}

func (c *ClusterInfo) SetMaster(master string) {
	c.master = master
}

func (c *ClusterInfo) AddLiveNode(node string) {
	c.liveNodes = append(c.liveNodes, node)
}

func (c *ClusterInfo) AddAllNode(node string) {
	c.allNodes = append(c.allNodes, node)
}

func (c *ClusterInfo) RemoveLiveNode(node string) {
	for i, n := range c.liveNodes {
		if n == node {
			c.liveNodes = append(c.liveNodes[:i], c.liveNodes[i+1:]...)
			break
		}
	}
}

func (c *ClusterInfo) RemoveAllNode(node string) {
	for i, n := range c.allNodes {
		if n == node {
			c.allNodes = append(c.allNodes[:i], c.allNodes[i+1:]...)
			break
		}
	}
}

func (c *ClusterInfo) Reset() {
	c.liveNodes = make([]string, 0)
	c.allNodes = make([]string, 0)
	c.master = ""
}
