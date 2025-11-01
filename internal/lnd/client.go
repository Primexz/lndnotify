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
	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/internal/events"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"google.golang.org/protobuf/proto"
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
	cfg             *config.Config
	conn            *grpc.ClientConn
	client          lnrpc.LightningClient
	state           lnrpc.StateClient
	router          routerrpc.RouterClient
	channelManager  *channelmanager.ChannelManager
	pendChanManager *channelmanager.PendingChannelManager
	pendChanUpdates chan proto.Message
	mu              sync.Mutex
	eventSub        chan events.Event
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
}

// NewClient creates a new LND client
func NewClient(cfg *config.Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		cfg:             cfg,
		eventSub:        make(chan events.Event, 100),
		pendChanUpdates: make(chan proto.Message, 100),
		ctx:             ctx,
		cancel:          cancel,
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
	tlsCert, err := credentials.NewClientTLSFromFile(c.cfg.LND.TLSCertPath, "")
	if err != nil {
		return fmt.Errorf("reading TLS cert: %w", err)
	}

	// Read macaroon
	macBytes, err := os.ReadFile(c.cfg.LND.MacaroonPath)
	if err != nil {
		return fmt.Errorf("reading macaroon: %w", err)
	}

	// Create gRPC connection
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", c.cfg.LND.Host, c.cfg.LND.Port),
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
	c.state = lnrpc.NewStateClient(conn)
	c.router = routerrpc.NewRouterClient(conn)
	c.channelManager = channelmanager.NewChannelManager(c.client)
	c.pendChanManager = channelmanager.NewPendingChannelManager(c.client, c.pendChanUpdates)

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

	if c.pendChanManager != nil {
		c.pendChanManager.Stop()
		c.pendChanManager = nil
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

	// standalone handlers that can be started right away
	initHandlers := []func(){
		c.handleLndWalletState,
		c.handleLndHealth,
	}
	c.wg.Add(len(initHandlers))
	for _, h := range initHandlers {
		go h()
	}

	go retry(c.ctx, "main client", func() (string, error) {
		if err := c.channelManager.Start(); err != nil {
			return "", fmt.Errorf("starting channel manager: %w", err)
		}

		if err := c.pendChanManager.Start(); err != nil {
			return "", fmt.Errorf("starting pending channel manager: %w", err)
		}

		// Start subscription handlers
		// NOTE: Keep handlers in alphabetical order to prevent merge conflicts when adding new handlers
		handlers := []func(){
			c.handleBackupEvents,
			c.handleChannelEvents,
			c.handleChannelFeeChanges,
			c.handleFailedHtlcEvents,
			c.handleForwards,
			c.handleInvoiceEvents,
			c.handleKeysendEvents,
			c.handleOnChainEvents,
			c.handlePaymentEvents,
			c.handlePeerEvents,
			c.handlePendingChannels,
			c.handleChainSyncState,
			c.handleChannelStatusEvents,
			c.handleTLSCertExpiry,
			c.handeLndVersion,
		}
		c.wg.Add(len(handlers))
		for _, h := range handlers {
			go h()
		}

		return "", nil
	})

	return c.eventSub, nil
}
