# core

This lib provides methods to build maas clients with defaults required parameters: logger, namespace, maasAgentUrl and authSupplier.

<!-- TOC -->
* [core](#core)
  * [Kafka](#kafka)
    * [Default usage:](#default-usage)
    * [You can override any of default parameters like shown in the code snippet below:](#you-can-override-any-of-default-parameters-like-shown-in-the-code-snippet-below)
    * [Watching tenant topics:](#watching-tenant-topics)
  * [Rabbit](#rabbit)
    * [Default usage:](#default-usage-1)
<!-- TOC -->

## Kafka

### Default usage:
~~~ go 
import (
	"context"
	"fmt"
	"github.com/netcracker/qubership-core-lib-go-maas-client/v3/classifier"
	maas "github.com/netcracker/qubership-core-lib-go-maas-core/v3"
)

func kafkaClientWithDefaults(ctx context.Context) error {
	maasKafkaClient := maas.NewKafkaClient()
	topicAddr, err := maasKafkaClient.GetTopic(ctx, classifier.New("demo").WithNamespace("namespace"))
	if err != nil {
		return err
	}
	fmt.Printf("topic = %s", topicAddr.TopicName)
	return nil
}
~~~

### Override default parameters
~~~ go 
import (
	"github.com/netcracker/qubership-core-lib-go-maas-client/v3/logging"
	maas "github.com/netcracker/qubership-core-lib-go-maas-core/v3"
)

var myNamespace string
var myMaaSAgentUrl string
var myAuthSupplier func(ctx context.Context) (string, error)

func kafkaClientWithCustomParams() {
    maasKafkaClient := maas.NewKafkaClient(
        maas.WithNamespace(myNamespace), 
        maas.WithMaaSAgentUrl(myMaaSAgentUrl), 
        maas.WithAuthSupplier(myAuthSupplier))
}
~~~

### Watching tenant topics:
See example [tenant-topics-watch.go](examples/tenant-topics-watch.go)

## Rabbit

### Default usage:
~~~ go 
import (
	"context"
	"fmt"
	"github.com/netcracker/qubership-core-lib-go-maas-client/v3/classifier"
	maas "github.com/netcracker/qubership-core-lib-go-maas-core/v3"
)

func kafkaClientWithDefaults(ctx context.Context) error {
	maasRabbitClient := maas.NewRabbitClient()
	vhost, err := maasRabbitClient.GetVhost(ctx, classifier.New("demo").WithNamespace("namespace"))
	if err != nil {
		return err
	}
	fmt.Printf("vhost user = %s", vhost.Username)
	return nil
}
~~~
