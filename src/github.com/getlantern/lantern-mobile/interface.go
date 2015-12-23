package client

import (
	"net"

	"github.com/getlantern/appdir"
	"github.com/getlantern/flashlight/config"
	"github.com/getlantern/flashlight/lantern"
	"github.com/getlantern/flashlight/logging"
	"github.com/getlantern/flashlight/settings"
	"github.com/getlantern/golog"
)

var (
	log         = golog.LoggerFor("lantern-android.client")
	appSettings *settings.Settings
)

type Provider interface {
    Verbose() bool
	Model() string
	Device() string
	Version() string
	AppName() string
	SettingsDir() string
	Protect(fileDescriptor int)
}

func Configure(provider Provider) error {

	log.Debugf("Configuring Lantern version: %s", lantern.GetVersion())

	settingsDir := provider.SettingsDir()
	log.Debugf("settings directory is %s", settingsDir)

	appdir.AndroidDir = settingsDir
	settings.SetAndroidPath(settingsDir)
	appSettings = settings.Load(lantern.GetVersion(), lantern.GetRevisionDate(), "")

    net.Callback = provider.Protect

    golog.Verbose = provider.Verbose()

	return nil
}

// Start creates a new client at the given address.
func Start(provider Provider) error {

	go func() {

		androidProps := map[string]string{
			"androidDevice":     provider.Device(),
			"androidModel":      provider.Model(),
			"androidSdkVersion": provider.Version(),
		}
		logging.ConfigureAndroid(androidProps)

		cfgFn := func(cfg *config.Config) {

		}

		l, err := lantern.Start(false, true, false,
			true, cfgFn)

		if l != nil && err != nil {
			log.Fatalf("Could not start Lantern")
            
		}
	}()
	return nil
}

func Restart(provider Provider) {
	log.Debugf("Restarting Lantern..")
	Stop()
	Configure(provider)
	Start(provider)
}

func Stop() {
	go lantern.Exit(nil)
}
