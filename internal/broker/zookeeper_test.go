
package broker

import (
	"testing"
	"time"

	"github.com/go-zookeeper/zk"
)

func TestCreatePersistentNode(t *testing.T) {
	conn, _, err := zk.Connect(
		[]string{"localhost:2181"},
		5*time.Second,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := &ZookeeperClient{conn: conn}

	path := "/test-kafka-zero/config"
	_, err = conn.Create(
		"/test-kafka-zero",
		[]byte{},
		0,
		zk.WorldACL(zk.PermAll),
	)
	if err != nil && err != zk.ErrNodeExists {
		t.Fatal(err)
	}

	if err := client.CreatePersistentNode(path, []byte("v1")); err != nil {
		t.Fatal(err)
	}
	if err := client.CreatePersistentNode(path, []byte("v2")); err != nil {
		t.Fatal(err)
	}

	data, _, err := conn.Get(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "v2" {
		t.Fatalf("got %q, want %q", data, "v2")
	}

	err = conn.Delete(path, -1)
	if err != nil {
		t.Fatal(err)
	}
	err = conn.Delete("/test-kafka-zero", -1)
	if err != nil {
		t.Fatal(err)
	}
}
