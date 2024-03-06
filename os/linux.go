package os

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/k0sproject/rig/cmd"
	"github.com/k0sproject/rig/kv"
	"github.com/k0sproject/rig/log"
)

// ResolveLinux resolves the OS release information for a linux host.
func ResolveLinux(conn cmd.SimpleRunner) (*Release, bool) {
	if conn.IsWindows() {
		return nil, false
	}

	if err := conn.Exec("uname | grep -q Linux"); err != nil {
		log.Trace(context.Background(), "linux os resolver: host is not linux", log.HostAttr(conn), log.ErrorAttr(err))
		return nil, false
	}

	reader := conn.ExecReader("cat /etc/os-release || cat /usr/lib/os-release")
	decoder := kv.NewDecoder(reader)

	version := &Release{}
	if err := decoder.Decode(version); err != nil {
		log.Trace(context.Background(), "linux os resolver: execreader returned an error", log.HostAttr(conn), log.ErrorAttr(err))
		return nil, false
	}

	str, _ := json.Marshal(version)
	log.Trace(context.Background(), "linux os resolver: resolved os release", log.HostAttr(conn), slog.String("data", string(str)))

	return version, true
}

// RegisterLinux registers the linux OS release resolver to a provider.
func RegisterLinux(provider *Provider) {
	provider.Register(ResolveLinux)
}
