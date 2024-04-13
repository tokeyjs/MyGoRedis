package cluster

import (
	"MyGoRedis/config"
	"MyGoRedis/resp/client"
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type connectionFactory struct {
	Peer string
}

func (cf *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	c, err := client.MakeClient(cf.Peer)
	if err != nil {
		return nil, err
	}
	c.Start()

	// 进行连接认证
	if len(config.Properties.RequirePass) > 0 {
		authData := make([][]byte, 0)
		authData = append(authData, []byte("auth"))
		authData = append(authData, []byte(config.Properties.RequirePass))
	}
	return pool.NewPooledObject(c), nil
}

func (cf *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	c.Close()
	return nil
}

func (cf *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (cf *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (cf *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
