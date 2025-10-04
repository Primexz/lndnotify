package lnd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Primexz/lnd-notify/internal/events"
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
	cfg      *ClientConfig
	conn     *grpc.ClientConn
	client   lnrpc.LightningClient
	router   routerrpc.RouterClient
	mu       sync.Mutex
	eventSub chan events.Event
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
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

	return nil
}

// Disconnect closes the connection to the LND node
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

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

	// Start subscription handlers
	c.wg.Add(2)
	go c.handleForwards()
	go c.handlePeerEvents()

	return c.eventSub, nil
}

// handleForwards polls for forwarding events
func (c *Client) handleForwards() {
	defer c.wg.Done()

	start := time.Now()
	for range time.Tick(time.Minute * 1) {

		resp, err := c.client.ForwardingHistory(c.ctx, &lnrpc.ForwardingHistoryRequest{
			StartTime:       uint64(start.Unix()),
			PeerAliasLookup: true,
		})
		if err != nil {
			fmt.Printf("Error fetching forwarding history: %v\n", err)
			continue
		}

		forwards := resp.GetForwardingEvents()
		for _, fwd := range forwards {

			c.eventSub <- events.NewForwardEvent(
				fwd.PeerAliasIn,
				fwd.PeerAliasOut,
				fwd.AmtInMsat,
				fwd.AmtOutMsat,
				fwd.FeeMsat,
			)
		}

		// push start time forward
		if len(forwards) > 0 {
			start = time.Now()
		}
	}
}

// handlePeerEvents handles peer connection and disconnection events
func (c *Client) handlePeerEvents() {
	ev, err := c.client.SubscribePeerEvents(c.ctx, &lnrpc.PeerEventSubscription{})
	if err != nil {
		fmt.Printf("Error subscribing to peer events: %v\n", err)
		return
	}

	for {
		peerEvent, err := ev.Recv()
		if err != nil {
			fmt.Printf("Error receiving peer event: %v\n", err)
			return
		}

		nodeInfo, err := c.client.GetNodeInfo(c.ctx, &lnrpc.NodeInfoRequest{
			PubKey: peerEvent.GetPubKey(),
		})
		if err != nil {
			fmt.Printf("Error fetching node info: %v\n", err)
			continue
		}

		switch peerEvent.GetType() {
		case lnrpc.PeerEvent_PEER_ONLINE:
			c.eventSub <- events.NewPeerOnlineEvent(nodeInfo.Node.Alias)
		case lnrpc.PeerEvent_PEER_OFFLINE:
			c.eventSub <- events.NewPeerOfflineEvent(nodeInfo.Node.Alias)
		}
	}
}

// MacaroonCredential implements the credentials.PerRPCCredentials interface
type MacaroonCredential struct {
	MacaroonHex string
}

func (m *MacaroonCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"macaroon": m.MacaroonHex,
	}, nil
}

func (m *MacaroonCredential) RequireTransportSecurity() bool {
	return true
}
