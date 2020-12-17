package spec

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
)

var specReplacer *strings.Replacer

func loadReplacer() {
	// Build a replacer to search/replace specification URLs
	if specReplacer == nil {
		var replacements []string

		// Configure the replacer with key=value pairs
		for k, v := range viper.GetStringMapString(config.SpecRewriteURL) {
			if v != "" {
				// Map between configured to=from URL pair
				replacements = append(replacements, k, v)
			} else {
				// Map between configured URL and site URL
				replacements = append(replacements, k, viper.GetString(config.SiteURL))
			}
		}

		specReplacer = strings.NewReplacer(replacements...)
	}
}

func replace(data []byte) []byte {
	return []byte(specReplacer.Replace(string(data)))
}
