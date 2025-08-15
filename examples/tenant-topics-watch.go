package examples

import (
	"context"
	"fmt"

	"github.com/vlla-test-organization/qubership-core-lib-go-maas-client/v3/classifier"
	"github.com/vlla-test-organization/qubership-core-lib-go-maas-client/v3/kafka/model"
	"github.com/vlla-test-organization/qubership-core-lib-go-maas-core/v3"
)

type health struct {
	healthy bool
	reason  string
}

func (t *health) Set(healthy bool, reason string) {
	t.healthy = healthy
	t.reason = reason
}

func watchTenantsTopic(ctx context.Context, tenantTopicWatchHealth *health) {
	client := core.NewKafkaClient()
	classifierKey := classifier.New("kafka-topic")
	go func() {
		for {
			watchCtx, cancelWatching := context.WithCancel(ctx)
			err := client.WatchTenantKafkaTopics(watchCtx, classifierKey, func(topics []model.TopicAddress, err error) {
				if err != nil {
					cancelWatching()
					return
				}
				// process topics...
			})
			if err != nil {
				tenantTopicWatchHealth.Set(false, fmt.Errorf("failed to start watching tenant topics: %w", err).Error())
			} else {
				tenantTopicWatchHealth.Set(true, "ok")
			}
			// block until callback receives cancel err and terminates our watcher, then retry watch
			<-watchCtx.Done()
		}
	}()
}
