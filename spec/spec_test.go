package spec

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
)

const testSpecDir = "../fixtures/"

func TestLoadSpecifications(t *testing.T) {
	config.Restore()
	viper.Set(config.SpecDir, testSpecDir)

	tests := []struct {
		name    string
		specLoc string
		wantErr bool
	}{
		{
			name:    "success - load common specifications",
			specLoc: "common_api.json",
			wantErr: false,
		},
		{
			name:    "success - load specifications with use of allof",
			specLoc: "allof_api.json",
			wantErr: false,
		},
		{
			name:    "success - load specifications with use of path depth",
			specLoc: "depth_api.json",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(config.SpecFilename, tt.specLoc)

			if err := LoadSpecifications(); (err != nil) != tt.wantErr {
				t.Errorf("LoadSpecifications() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
