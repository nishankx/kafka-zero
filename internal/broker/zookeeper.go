package broker

import (
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

type ZookeeperClient struct {
	conn *zk.Conn
}

func NewZookeeperClient(
	servers []string,
	sessionTimeout time.Duration,
) (*ZookeeperClient, error) {
	conn, _, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return nil, fmt.Errorf("connect to ZooKeeper: %w", err)
	}

	client := &ZookeeperClient{
		conn: conn,
	}

	return client, nil
}

func (c *ZookeeperClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}


func (c *ZookeeperClient) CreatePersistentNode(
	path string,
	data []byte,
) error {
	_, _, err := c.conn.Exists(path)
	if err != nil && err != zk.ErrNoNode {
		return fmt.Errorf("check node %q: %w", path, err)
	}

	if err == nil {
		_, err = c.conn.Set(path, data, -1)
		if err != nil {
			return fmt.Errorf("update node %q: %w", path, err)
		}
		return nil
	}

	_, err = c.conn.Create(
		path,
		data,
		0,
		zk.WorldACL(zk.PermAll),
	)
	if err != nil {
		return fmt.Errorf("create persistent node %q: %w", path, err)
	}

	return nil
}

func (c *ZookeeperClient) CreateEphemeralNode(
	path string,
	data []byte,
) (bool, error) {
	_, _, err := c.conn.Exists(path)
	if err == nil {
		return false, nil
	}
	if err != zk.ErrNoNode {
		return false, fmt.Errorf("check node %q: %w", path, err)
	}

	_, err = c.conn.Create(
		path,
		data,
		zk.FlagEphemeral,
		zk.WorldACL(zk.PermAll),
	)
	if err != nil {
		if err == zk.ErrNodeExists {
			return false, nil
		}
		return false, fmt.Errorf("create ephemeral node %q: %w", path, err)
	}

	return true, nil
}
