/*
Copyright (C) 2016-2017 dapperdox.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

// Package spec provides API spec loading and parsing.
package spec

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
)

var (
	statusMapSplit = regexp.MustCompile(",")
	statusCodes    map[int]string
)

// loadStatusCodes loads status code mappings.
func loadStatusCodes() {
	var statusfile string

	if viper.GetString(config.AssetsDir) != "" {
		statusfile = filepath.Join(viper.GetString(config.AssetsDir), "status_codes.csv")
		log().Tracef("Looking in assets dir for %s", statusfile)

		if _, err := os.Stat(statusfile); os.IsNotExist(err) {
			statusfile = ""
		}
	}

	if statusfile == "" && viper.GetString(config.ThemeDir) != "" {
		statusfile = filepath.Join(viper.GetString(config.ThemeDir), viper.GetString(config.Theme), "status_codes.csv")
		log().Tracef("Looking in theme dir for %s", statusfile)

		if _, err := os.Stat(statusfile); os.IsNotExist(err) {
			statusfile = ""
		}
	}

	if statusfile == "" {
		statusfile = filepath.Join(viper.GetString(config.DefaultAssetsDir), "themes", viper.GetString(config.Theme), "status_codes.csv")
		log().Tracef("Looking in default theme dir for %s", statusfile)

		if _, err := os.Stat(statusfile); os.IsNotExist(err) {
			statusfile = ""
		}
	}

	if statusfile == "" {
		statusfile = filepath.Join(viper.GetString(config.DefaultAssetsDir), "themes", "default", "status_codes.csv")
		log().Tracef("Looking in default theme %s", statusfile)

		if _, err := os.Stat(statusfile); os.IsNotExist(err) {
			statusfile = ""
		}
	}

	if statusfile == "" {
		log().Trace("No status code map file found.")

		return
	}

	log().Tracef("Processing HTTP status code file: %s", statusfile)

	file, err := os.Open(statusfile)
	if err != nil {
		log().Errorf("Error: %s", err)

		return
	}
	defer file.Close()

	statusCodes = make(map[int]string)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		indexes := statusMapSplit.FindStringIndex(line)
		if indexes == nil {
			return
		}

		i, err := strconv.Atoi(line[0 : indexes[1]-1])
		if err != nil {
			log().Errorf("Invalid HTTP status code in csv file: %q", line)

			continue
		}

		statusCodes[i] = line[indexes[1]:]
	}

	if err := scanner.Err(); err != nil {
		log().Errorf("Error: %s", err)
	}
}

func httpStatusDescription(status int) string {
	return statusCodes[status]
}
