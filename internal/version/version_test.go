package version

import (
	"runtime/debug"
	"testing"
)

func TestFromBuildInfo(t *testing.T) {
	tests := []struct {
		name string
		info *debug.BuildInfo
		ok   bool
		want string
	}{
		{name: "tagged", info: &debug.BuildInfo{Main: debug.Module{Version: "v0.0.3"}}, ok: true, want: "0.0.3"},
		{name: "pseudo version", info: &debug.BuildInfo{Main: debug.Module{Version: "v0.0.3-0.20261201000000-abcdef123456+dirty"}}, ok: true, want: "0.0.3-0.20261201000000-abcdef123456+dirty"},
		{name: "devel", info: &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, ok: true, want: "dev"},
		{name: "empty", info: &debug.BuildInfo{}, ok: true, want: "dev"},
		{name: "no build info", info: nil, ok: false, want: "dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromBuildInfo(tt.info, tt.ok); got != tt.want {
				t.Errorf("fromBuildInfo = %q, want %q", got, tt.want)
			}
		})
	}
}

// Under "go test", the build records no module version.
func TestGet(t *testing.T) {
	if got := Get(); got != devVersion {
		t.Errorf("Get() = %q, want %q", got, devVersion)
	}
}
