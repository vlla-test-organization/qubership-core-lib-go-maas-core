package core

import (
	"context"
	"fmt"

	"github.com/netcracker/qubership-core-lib-go-maas-client/v3/kafka"
	"github.com/netcracker/qubership-core-lib-go-maas-client/v3/rabbit"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	"github.com/netcracker/qubership-core-lib-go/v3/utils"
	"github.com/netcracker/qubership-core-lib-go/v3/const"
	"github.com/netcracker/qubership-core-lib-go/v3/serviceloader"
)

type options struct {
	namespace        func() string
	maasAgentUrl     func() string
	tenantManagerUrl func() string
	httpClient       func() *resty.Client
	stompDialer      func() *websocket.Dialer
	authSupplier     func() func(ctx context.Context) (string, error)
}

type Option func(options *options)

func init() {
	serviceloader.Register(2, &serviceloader.Token{})
}

func NewKafkaClient(opts ...Option) kafka.MaasClient {
	config := configure(opts...)
	return kafka.NewClient(config.namespace(), config.maasAgentUrl(), config.tenantManagerUrl(), config.httpClient(),
		config.stompDialer(), config.authSupplier())
}

func NewRabbitClient(opts ...Option) rabbit.MaasClient {
	config := configure(opts...)
	return rabbit.NewClient(config.namespace(), config.maasAgentUrl(), config.httpClient())
}

func configure(opts ...Option) *options {
	config := &options{
		namespace:        getNamespace,
		maasAgentUrl:     getMaaSAgentUrl,
		tenantManagerUrl: getTenantManagerUrl,
		httpClient:       getHttpClient,
		stompDialer:      getStompDialer,
		authSupplier:     getAuthSupplier,
	}
	for _, option := range opts {
		option(config)
	}
	return config
}

func WithNamespace(namespace string) Option {
	return func(options *options) { options.namespace = func() string { return namespace } }
}

func WithMaaSAgentUrl(url string) Option {
	return func(options *options) { options.maasAgentUrl = func() string { return url } }
}

func WithHttpClient(client *resty.Client) Option {
	return func(options *options) { options.httpClient = func() *resty.Client { return client } }
}

func WithStompDialer(stompDialer *websocket.Dialer) Option {
	return func(options *options) { options.stompDialer = func() *websocket.Dialer { return stompDialer } }
}

func WithAuthSupplier(authSupplier func(ctx context.Context) (string, error)) Option {
	return func(options *options) {
		options.authSupplier = func() func(ctx context.Context) (string, error) { return authSupplier }
	}
}

func getMaaSAgentUrl() string {
	defaultUrl := constants.SelectUrl("http://maas-agent:8080", "https://maas-agent:8443")
	return configloader.GetOrDefaultString("maas.agent.url", defaultUrl)
}

func getTenantManagerUrl() string {
	defaultUrl := constants.SelectUrl("ws://tenant-manager:8080", "wss://tenant-manager:8443")
	return configloader.GetOrDefaultString("tenant.manager.url", defaultUrl)
}

func getNamespace() string {
	return configloader.GetKoanf().MustString("microservice.namespace")
}

func getHttpClient() *resty.Client {
	return resty.New().OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
	    tokenProvider := serviceloader.MustLoad[serviceloader.TokenProvider]()
		token, err := tokenProvider.GetToken(request.Context())
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}
		request.SetAuthToken(token)
		return nil
	}).SetTLSClientConfig(utils.GetTlsConfig()).SetRetryCount(10)
}

func getStompDialer() *websocket.Dialer {
	return &websocket.Dialer{TLSClientConfig: utils.GetTlsConfig()}
}

func getAuthSupplier() func(ctx context.Context) (string, error) {
    tokenProvider := serviceloader.MustLoad[serviceloader.TokenProvider]()
	return tokenProvider.GetToken
}