package cmds

import (
	"bytes"
	"context"
	"fmt"
	"git.poundadm.net/anachronism/xmcctl/pkg/apis/config"
	"git.poundadm.net/anachronism/xmcctl/pkg/remote"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"time"
)

var (
	RootCommand *cobra.Command

	ConfigPath        string
	DiscoverBindAddr  string
	DiscoverBroadcast bool
	DiscoverRefresh   bool
	DiscoverTimeout   time.Duration
	DiscoverWrite     bool
	LogDebug          bool
	LogJson           bool

	conf *config.Config
)

func init() {
	cobra.OnInitialize(
		setupConfig,
		setupLogger,
	)
}

func New() *cobra.Command {
	RootCommand = &cobra.Command{
		Use:   "xmcctl",
		Short: "Remote control Emotiva XMC-1 devices over IP.",
	}

	RootCommand.PersistentFlags().BoolVar(&LogDebug, "log-debug", false, "Enable debug logging")
	RootCommand.PersistentFlags().BoolVar(&LogJson, "log-json", false, "Enable JSON logging")
	RootCommand.PersistentFlags().StringVar(&ConfigPath, "conf", "~/.conf/xmcctl.yaml", "Path to the conf file.")

	discoverCommand := &cobra.Command{
		Use:   "discover [flags] [ip...]",
		Short: "Find Emotiva devices on the network.",
		Long: `Searches for devices that respond to Emotiva's transponder identification packets.

The broadcast IP address is used by default, which should find any devices on
the local network. If IPs are supplied, those IPs will be attempted instead.

This command can be used to write any discovered devices to a conf file.
`,
		RunE: discoverCmd,
	}
	discoverCommand.Flags().BoolVarP(&DiscoverWrite, "write", "w", false, "Write discovered devices to the conf file.")
	discoverCommand.Flags().BoolVarP(&DiscoverRefresh, "refresh", "r", false, "Attempt to rediscover any configured devices.")
	discoverCommand.Flags().BoolVarP(&DiscoverBroadcast, "broadcast", "B", false, "Always broadcast a discovery request.")
	discoverCommand.Flags().StringVarP(&DiscoverBindAddr, "bind", "b", "0.0.0.0", "IP address to listen for discovery responses on.")
	discoverCommand.Flags().DurationVarP(&DiscoverTimeout, "timeout", "t", 3*time.Second, "Maximum duration to wait for discovery responses.")
	RootCommand.AddCommand(discoverCommand)

	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Prints the version.",
		Run:   versionCmd,
	}
	RootCommand.AddCommand(versionCommand)

	return RootCommand
}

func setupConfig() {
	conf = &config.Config{}
}

func setupLogger() {
	if LogDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if LogJson {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func versionCmd(cmd *cobra.Command, args []string) {
	fmt.Println("no versions yet :(")
}

func discoverCmd(cmd *cobra.Command, args []string) error {
	bindIP := net.ParseIP(DiscoverBindAddr)
	if bindIP == nil {
		return errors.New("unable to parse bind address: " + DiscoverBindAddr)
	}

	// args to this command should be IP addresses
	addrs := make([]net.IP, 0)
	if len(args) > 0 {
		if DiscoverBroadcast {
			addrs = append(addrs, net.IPv4bcast)
		}
		for _, arg := range args {
			ip := net.ParseIP(arg)
			if ip == nil {
				return errors.New("unable to parse ip address: " + arg)
			}

			// don't add duplicate entries for the broadcast address
			if DiscoverBroadcast {
				if bytes.Compare(ip, net.IPv4bcast) == 0 {
					continue
				}
			}

			addrs = append(addrs, ip)
		}
	} else {
		addrs = append(addrs, net.IPv4bcast)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DiscoverTimeout)
	remotes, err := remote.DiscoverTransponders(ctx, bindIP, addrs)
	if err != nil {
		cancel()
		return errors.New("error discovering devices: " + err.Error())
	}

	for _, remote := range remotes {
		fmt.Printf("%+v\n", remote)
	}

	return nil
}
