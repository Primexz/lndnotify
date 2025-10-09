package lnd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	channelmanager "github.com/Primexz/lndnotify/internal/channel_manager"
	"github.com/Primexz/lndnotify/internal/events"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
)

// ClientConfig holds the configuration for the LND client
type ClientConfig struct {
	Host         string
	Port         int
	TLSCertPath  string
	MacaroonPath string
}

// Client represents an LND node client
type Client struct {
	cfg            *ClientConfig
	conn           *grpc.ClientConn
	client         lnrpc.LightningClient
	router         routerrpc.RouterClient
	channelManager *channelmanager.ChannelManager
	mu             sync.Mutex
	eventSub       chan events.Event
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewClient creates a new LND client
func NewClient(cfg *ClientConfig) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		cfg:      cfg,
		eventSub: make(chan events.Event, 100),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Connect establishes a connection to the LND node
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil // Already connected
	}

	// Read TLS certificate
	tlsCert, err := credentials.NewClientTLSFromFile(c.cfg.TLSCertPath, "")
	if err != nil {
		return fmt.Errorf("reading TLS cert: %w", err)
	}

	// Read macaroon
	macBytes, err := os.ReadFile(c.cfg.MacaroonPath)
	if err != nil {
		return fmt.Errorf("reading macaroon: %w", err)
	}

	// Create gRPC connection
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port),
		grpc.WithTransportCredentials(tlsCert),
		grpc.WithPerRPCCredentials(&MacaroonCredential{
			MacaroonHex: hex.EncodeToString(macBytes),
		}),
	)
	if err != nil {
		return fmt.Errorf("connecting to LND: %w", err)
	}

	c.conn = conn
	c.client = lnrpc.NewLightningClient(conn)
	c.router = routerrpc.NewRouterClient(conn)
	c.channelManager = channelmanager.NewChannelManager(c.client)

	return nil
}

// Disconnect closes the connection to the LND node
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()
	c.wg.Wait()

	if c.channelManager != nil {
		c.channelManager.Stop()
		c.channelManager = nil
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("closing connection: %w", err)
		}
		c.conn = nil
		c.client = nil
	}
	return nil
}

// IsConnected returns true if connected to LND
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn != nil
}

// SubscribeEvents subscribes to LND events
func (c *Client) SubscribeEvents() (<-chan events.Event, error) {
	if !c.IsConnected() {
		if err := c.Connect(); err != nil {
			return nil, fmt.Errorf("connecting to LND: %w", err)
		}
	}

	if err := c.channelManager.Start(); err != nil {
		return nil, fmt.Errorf("starting channel manager: %w", err)
	}

	// Start subscription handlers
	handlers := []func(){
		c.handleForwards,
		c.handlePeerEvents,
		c.handleChannelEvents,
		c.handleInvoiceEvents,
		c.handleFailedHtlcEvents,
		c.handleKeysendEvents,
	}
	c.wg.Add(len(handlers))
	for _, h := range handlers {
		go h()
	}

	return c.eventSub, nil
}
