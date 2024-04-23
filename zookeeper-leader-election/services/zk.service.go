package services

import (
	"errors"
	"fmt"
	"time"
	ct "zk/constants"

	"github.com/go-zookeeper/zk"
)

type ZookeeperService interface {
	Exists(path string) (bool, *zk.Stat, error)
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	Get(path string) ([]byte, *zk.Stat, error)
	Children(path string) ([]string, *zk.Stat, error)
	Close()

	// for our reference
	CreateAllParentNodes() error
	AddToAllNode(node, data string) error
	GetAllNodes() []string
	AddToElectionNodes(node, data string) error
	GetElectionNodes() []string
	AddToElectionNode(node, data string) error
	GetElectionNode() []string
	AddToLiveNode(node, data string) error
	GetLiveNode() []string
	SetMaster(node string) error
}

type zkService struct {
	zkClient *zk.Conn
}

func NewZookeeperService(connStrings []string) ZookeeperService {
	client, _, err := zk.Connect(connStrings, 10*time.Second)
	if err != nil {
		panic(err)
	}
	return &zkService{zkClient: client}
}

func (z *zkService) Exists(path string) (bool, *zk.Stat, error) {
	return z.zkClient.Exists(path)
}

func (z *zkService) Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return z.zkClient.Create(path, data, flags, acl)
}

func (z *zkService) Get(path string) ([]byte, *zk.Stat, error) {
	return z.zkClient.Get(path)
}

func (z *zkService) Children(path string) ([]string, *zk.Stat, error) {
	return z.zkClient.Children(path)
}

func (z *zkService) Close() {
	z.zkClient.Close()
}

func (z *zkService) CreateAllParentNodes() error {
	// if all nodes doesnot exists, create all nodes
	if exists, _, _ := z.Exists(ct.ALL_NODES); !exists {
		if _, err := z.Create(ct.ALL_NODES, []byte("all nodes can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			if !errors.Is(err, zk.ErrNodeExists) {
				return err
			}
		}
	}

	// if live nodes are present, create live nodes
	if exists, _, _ := z.Exists(ct.LIVE_NODES); !exists {
		if _, err := z.Create(ct.LIVE_NODES, []byte("all live can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			if !errors.Is(err, zk.ErrNodeExists) {
				return err
			}
		}
	}

	// if election nodes doesnot exists, create persistent leader node
	if exists, _, _ := z.Exists(ct.ELECTION_NODES); !exists {
		if _, err := z.Create(ct.ELECTION_NODES, []byte("election node can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			if !errors.Is(err, zk.ErrNodeExists) {
				return err
			}
		}
	}

	if exists, _, _ := z.Exists(ct.ELECTION_NODE); !exists {
		if _, err := z.Create(ct.ELECTION_NODE, []byte("election node can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			if !errors.Is(err, zk.ErrNodeExists) {
				return err
			}
		}
	}

	if exists, _, _ := z.Exists(ct.ELECTION_LEADER_NODE); !exists {
		if _, err := z.Create(ct.ELECTION_LEADER_NODE, []byte("leader node can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			if !errors.Is(err, zk.ErrNodeExists) {
				return err
			}
		}

	}

	fmt.Println("All parent nodes created successfully...")
	return nil
}

func (z *zkService) AddToAllNode(node, data string) error {
	// if all nodes doesnot exists, create all nodes
	if exists, _, _ := z.Exists(ct.ALL_NODES); !exists {
		if _, err := z.Create(ct.ALL_NODES, []byte("all nodes can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}

	path := fmt.Sprintf("%s/%s", ct.ALL_NODES, node)
	fmt.Println(path)
	/// create persistent node
	if _, err := z.zkClient.Create(path, []byte(data), 0, zk.WorldACL(zk.PermAll)); err != nil {
		if !errors.Is(err, zk.ErrNodeExists) {
			return err
		}
	}

	fmt.Println(node, " added to all nodes successfully...")
	return nil
}

func (z *zkService) GetAllNodes() []string {
	nodes, _, _ := z.Children(ct.ALL_NODES)
	return nodes
}

func (z *zkService) AddToElectionNodes(node, data string) error {
	// if election node doesnot exists, create election node
	if exists, _, _ := z.Exists(ct.ELECTION_NODES); !exists {
		if _, err := z.Create(ct.ELECTION_NODES, []byte("election node can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}

	// create ephemeral + sequential node
	path := fmt.Sprintf("%s/%s-", ct.ELECTION_NODES, node)
	if _, err := z.Create(path, []byte(data), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll)); err != nil {
		return err
	}
	fmt.Println(node, " added to election node successfully...")
	return nil
}

func (z *zkService) GetElectionNodes() []string {
	nodes, _, _ := z.Children(ct.ELECTION_NODES)
	return nodes
}

func (z *zkService) AddToElectionNode(node, data string) error {
	// if election node doesnot exists, create election node
	if exists, _, _ := z.Exists(ct.ELECTION_NODE); !exists {
		if _, err := z.Create(ct.ELECTION_NODE, []byte("election node can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}

	// create ephemeral + sequential node
	path := fmt.Sprintf("%s/%s-", ct.ELECTION_NODE, node)
	if _, err := z.Create(path, []byte(data), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll)); err != nil {
		return err
	}
	fmt.Println(node, " added to election node successfully...")
	return nil
}

func (z *zkService) GetElectionNode() []string {
	nodes, _, _ := z.Children(ct.ELECTION_NODE)
	return nodes
}

func (z *zkService) SetMaster(node string) error {
	// if leader node doesnot exists, create persistent leader node
	exists, status, _ := z.Exists(ct.ELECTION_LEADER_NODE)
	if !exists {
		if _, err := z.Create(ct.ELECTION_LEADER_NODE, []byte("leader node can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}

	bytesData, _, err := z.Get(fmt.Sprintf("%s/%s", ct.ALL_NODES, node))
	if err != nil {
		return err
	}
	// set the master
	if _, err := z.zkClient.Set(ct.ELECTION_LEADER_NODE, bytesData, status.Version); err != nil {
		return err
	}
	fmt.Println(node, " set as master successfully...")
	return nil
}

func (z *zkService) AddToLiveNode(node, data string) error {
	// if election node doesnot exists, create election node
	if exists, _, _ := z.Exists(ct.LIVE_NODES); !exists {
		if _, err := z.Create(ct.LIVE_NODES, []byte("live node can be displayed here"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}

	// create ephemeral + sequential node
	path := fmt.Sprintf("%s/%s-", ct.LIVE_NODES, node)
	if _, err := z.Create(path, []byte(data), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll)); err != nil {
		return err
	}
	fmt.Println(node, " added to election node successfully...")
	return nil
}

func (z *zkService) GetLiveNode() []string {
	nodes, _, _ := z.Children(ct.LIVE_NODES)
	return nodes
}
