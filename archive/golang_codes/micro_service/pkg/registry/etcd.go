package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Port    int      `json:"port"`
	Tags    []string `json:"tags"`
	Health  bool     `json:"health"`
}

// Client 服务注册发现客户端接口
type Client interface {
	Register(ctx context.Context, service *ServiceInfo) error
	Deregister(ctx context.Context, service *ServiceInfo) error
	Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error)
	Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInfo, error)
	Close() error
}

// EtcdClient etcd客户端实现
type EtcdClient struct {
	client  *clientv3.Client
	leaseID clientv3.LeaseID
	ttl     int64
}

// NewEtcdClient 创建etcd客户端
func NewEtcdClient(endpoints []string) (Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Status(ctx, endpoints[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	logrus.Info("Successfully connected to etcd")

	return &EtcdClient{
		client: client,
		ttl:    30, // 30秒TTL
	}, nil
}

// Register 注册服务
func (e *EtcdClient) Register(ctx context.Context, service *ServiceInfo) error {
	// 创建租约
	leaseResp, err := e.client.Grant(ctx, e.ttl)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}
	e.leaseID = leaseResp.ID

	// 序列化服务信息
	service.Health = true
	serviceData, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %w", err)
	}

	// 构建服务键
	serviceKey := fmt.Sprintf("/services/%s/%s:%d", service.Name, service.Address, service.Port)

	// 注册服务
	_, err = e.client.Put(ctx, serviceKey, string(serviceData), clientv3.WithLease(e.leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// 启动续约
	go e.keepAlive()

	logrus.WithFields(logrus.Fields{
		"service": service.Name,
		"address": service.Address,
		"port":    service.Port,
	}).Info("Service registered successfully")

	return nil
}

// Deregister 注销服务
func (e *EtcdClient) Deregister(ctx context.Context, service *ServiceInfo) error {
	serviceKey := fmt.Sprintf("/services/%s/%s:%d", service.Name, service.Address, service.Port)

	_, err := e.client.Delete(ctx, serviceKey)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	// 撤销租约
	if e.leaseID != 0 {
		_, err = e.client.Revoke(ctx, e.leaseID)
		if err != nil {
			logrus.WithError(err).Warn("Failed to revoke lease")
		}
	}

	logrus.WithFields(logrus.Fields{
		"service": service.Name,
		"address": service.Address,
		"port":    service.Port,
	}).Info("Service deregistered successfully")

	return nil
}

// Discover 发现服务
func (e *EtcdClient) Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
	servicePrefix := fmt.Sprintf("/services/%s/", serviceName)

	resp, err := e.client.Get(ctx, servicePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}

	var services []*ServiceInfo
	for _, kv := range resp.Kvs {
		var service ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			logrus.WithError(err).Warn("Failed to unmarshal service info")
			continue
		}
		services = append(services, &service)
	}

	return services, nil
}

// Watch 监听服务变化
func (e *EtcdClient) Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInfo, error) {
	servicePrefix := fmt.Sprintf("/services/%s/", serviceName)
	resultChan := make(chan []*ServiceInfo, 1)

	// 首次获取服务列表
	services, err := e.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	resultChan <- services

	// 监听变化
	go func() {
		defer close(resultChan)
		watchChan := e.client.Watch(ctx, servicePrefix, clientv3.WithPrefix())

		for watchResp := range watchChan {
			if watchResp.Err() != nil {
				logrus.WithError(watchResp.Err()).Error("Watch error")
				return
			}

			// 服务发生变化，重新获取服务列表
			services, err := e.Discover(ctx, serviceName)
			if err != nil {
				logrus.WithError(err).Error("Failed to discover services on watch")
				continue
			}

			select {
			case resultChan <- services:
			case <-ctx.Done():
				return
			}
		}
	}()

	return resultChan, nil
}

// Close 关闭客户端
func (e *EtcdClient) Close() error {
	return e.client.Close()
}

// keepAlive 保持租约活跃
func (e *EtcdClient) keepAlive() {
	ch, kaerr := e.client.KeepAlive(context.Background(), e.leaseID)
	if kaerr != nil {
		logrus.WithError(kaerr).Error("Failed to keep alive")
		return
	}

	for ka := range ch {
		logrus.WithField("lease_id", ka.ID).Debug("Lease renewed")
	}
}
