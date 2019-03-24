package artnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/StageAutoControl/controller/pkg/internal/logging"
	"github.com/jsimonetti/go-artnet"
	"github.com/sirupsen/logrus"
)

// Controller is a transport for the ArtNet protocol (DMX over UDP/IP)
type controller struct {
	logger      logging.Logger
	sender      *artnet.Controller
	state       State
	sendTrigger chan struct{}
}

// NewController returns a artnet Controller as an anonymous interface
func NewController(logger logging.Logger) (Controller, error) {
	ip, err := FindArtNetIP()
	if err != nil {
		return nil, fmt.Errorf("failed to find the art-net IP: %v", err)
	}

	if len(ip) == 0 {
		return nil, errors.New("failed to find the art-net IP: No interface found")
	}

	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	host = strings.ToLower(strings.Split(host, ".")[0])

	logger.Infof("Using ArtNet IP %s and hostname %s", ip.String(), host)

	c := artnet.NewController(host, ip, artnet.NewLogger(logger.(*logrus.Entry).WithField("module", "artnet")))
	if err := c.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Controller: %v", err)
	}

	cntl := &controller{
		logger:      logger,
		sender:      c,
		state:       NewState(),
		sendTrigger: make(chan struct{}, 1),
	}

	go cntl.sendBackground()
	go cntl.debugDevices()

	return cntl, nil
}

// Start the controller
func (c *controller) Start() error {
	return c.sender.Start()
}

// Stop the controller
func (c *controller) Stop() {
	close(c.sendTrigger)
	c.sender.Stop()
}

func (c *controller) SetDMXChannelValue(value ChannelValue) {
	c.state.Set(value.Universe, value.Channel, value.Value)
	c.triggerSend()
}

func (c *controller) SetDMXChannelValues(values []ChannelValue) {
	for _, value := range values {
		c.state.Set(value.Universe, value.Channel, value.Value)
	}

	c.triggerSend()
}

func (c *controller) triggerSend() {
	c.sendTrigger <- struct{}{}
}

func (c *controller) sendBackground() {
	for range c.sendTrigger {
		c.logger.Debug("Sending DMX Values")
		c.send()
	}
}

func (c *controller) send() {
	for universe, dmx := range c.state {
		go c.sender.SendDMXToAddress(dmx, c.universeToAddress(universe))
	}
}

// universeToAddress converts a dmx universe to a artnet address
// Code stolen from https://play.golang.org/p/pdQPC5u7JX
func (c *controller) universeToAddress(universe uint16) artnet.Address {
	v := make([]uint8, 2)
	binary.BigEndian.PutUint16(v, universe)

	return artnet.Address{
		Net:    v[0],
		SubUni: v[1],
	}
}

func (c *controller) debugDevices() {
	t := time.NewTicker(30 * time.Second)
	for range t.C {
		c.logger.Debugf("Currently %d devices are registered: %+s", len(c.sender.Nodes), ips(c.sender.Nodes))
	}
}

func ips(nodes []*artnet.ControlledNode) (ips []string) {
	for _, n := range nodes {
		ips = append(ips, NodeToString(n))
	}
	return
}