package proxiedsites

import (
	"fmt"
	"sync"

	"github.com/getlantern/detour"
	"github.com/getlantern/golog"
	"github.com/getlantern/proxiedsites"

	"github.com/getlantern/flashlight/config"
	"github.com/getlantern/flashlight/ui"
)

const (
	messageType = `ProxiedSites`
)

var (
	log = golog.LoggerFor("flashlight.proxiedsites")

	service    *ui.Service
	PACURL     string
	startMutex sync.Mutex
)

func Configure(cfg *proxiedsites.Config) {
	delta := proxiedsites.Configure(cfg)
	startMutex.Lock()

	if delta != nil {
		updateDetour(delta)
	}

    startMutex.Unlock()
    return
}

func updateDetour(delta *proxiedsites.Delta) {
	// TODO: subscribe changes of geolookup and set country accordingly
	// safe to hardcode here as IR has all detection rules
	detour.SetCountry("IR")

	// for simplicity, detour matches whitelist using host:port string
	// so we add ports to each proxiedsites
	for _, v := range delta.Deletions {
		detour.RemoveFromWl(v + ":80")
		detour.RemoveFromWl(v + ":443")
	}
	for _, v := range delta.Additions {
		detour.AddToWl(v+":80", true)
		detour.AddToWl(v+":443", true)
	}
}

func start() (err error) {
	newMessage := func() interface{} {
		return &proxiedsites.Delta{}
	}

	// Registering a websocket service.
	helloFn := func(write func(interface{}) error) error {
		return write(proxiedsites.ActiveDelta())
	}

	if service, err = ui.Register(messageType, newMessage, helloFn); err != nil {
		return fmt.Errorf("Unable to register channel: %q", err)
	}

	// Initializing reader.
	go read()

	return nil
}

func read() {
	for msg := range service.In {
		err := config.Update(func(updated *config.Config) error {
			log.Debugf("Applying update from UI")
			updated.ProxiedSites.Delta.Merge(msg.(*proxiedsites.Delta))
			return nil
		})
		if err != nil {
			log.Debugf("Error applying update from UI: %v", err)
		}
	}
}
