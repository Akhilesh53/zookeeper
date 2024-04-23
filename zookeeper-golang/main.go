package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

func main() {
	// Connect to the ZooKeeper server
	zkConn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, 10*time.Second)
	if err != nil {
		panic(err)
	}

	//===============================================================================
	// check whether a node exists in ZooKeeper or not
	// zookeeper path should exists
	exists, status, err := zkConn.Exists("/zookeeper")
	if err != nil {
		fmt.Println("err checking path /zooeeper : ", err)
		return
	}
	fmt.Println("checking path /zooeeper : ", exists, status.Version)

	//===============================================================================
	// test path does not exists
	exists, status, err = zkConn.Exists("/test")
	if err != nil {
		fmt.Println("err checking path /test : ", err)
		return
	}
	fmt.Println("checking path /test : ", exists, status.Version)

	//===============================================================================
	// Create a path /golang_test in ZooKeeper and store a json key value object
	data := []byte(`{"key": "value"}`)
	path := "/golang_test"
	acls := zk.WorldACL(zk.PermAll)

	_, err = zkConn.Create(path, data, 0, acls)
	if err != nil && !errors.Is(err, zk.ErrNodeExists) {
		fmt.Println("err while creating path : ", err)
		return
	}

	//===============================================================================
	// retrieve data from the path
	data, status, err = zkConn.Get(path)
	if err != nil {
		fmt.Println("err while getting data : ", err)
		return
	}
	fmt.Println("fetched data : ", string(data), status.Version)

	//===============================================================================
	// get child of the node
	childs, status, err := zkConn.Children(path)
	if err != nil {
		fmt.Println("err while getting child nodes : ", err)
		return
	}
	fmt.Println("child of ", path, " are : ", childs)

	//===============================================================================
	// delete the node
	if err := zkConn.Delete(path, status.Version); err != nil {
		fmt.Println("err delete znode : ", err)
		return
	}
	fmt.Println("deleted znode : ", path)

	//===============================================================================

}
