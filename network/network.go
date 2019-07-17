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

// Package network creates connections with or without TLS.
package network

import (
	"crypto/tls"
	"net"

	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
)

// NewListener creates a new network Listener
func NewListener() (net.Listener, error) {
	log().Infof("listening on %s for unsecured connections", viper.GetString(config.BindAddr))
	return net.Listen("tcp", viper.GetString(config.BindAddr))
}

// NewSecuredListener creates a secure network Listener
func NewSecuredListener() (net.Listener, error) {
	crt, err := tls.LoadX509KeyPair(viper.GetString(config.TLSCert), viper.GetString(config.TLSKey))
	if err != nil {
		return nil, err
	}

	tlscfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{crt},
	}

	log().Infof("listening on %s for SECURED connections", viper.GetString(config.BindAddr))
	return tls.Listen("tcp", viper.GetString(config.BindAddr), tlscfg)
}
