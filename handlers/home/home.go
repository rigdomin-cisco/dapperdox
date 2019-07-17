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

// Package home provides handler for homepage.
package home

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
	"github.com/kenjones-cisco/dapperdox/render"
	"github.com/kenjones-cisco/dapperdox/spec"
)

// Register creates routes for each home handler
func Register(r *pat.Router) {
	log().Debug("registering handlers for home page")

	// Homepages for each loaded specification
	var specification *spec.APISpecification // Ends up being populated with the last spec processed

	for _, specification = range spec.APISuite {

		log().Tracef("Build homepage route for specification %q", specification.ID)

		r.Path("/" + specification.ID + "/reference").Methods(http.MethodGet).HandlerFunc(specificationSummaryHandler(specification))

		// If missingh trailing slash, redirect to add it
		r.Path("/" + specification.ID).Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/"+specification.ID+"/", http.StatusFound)
		})
	}

	if len(spec.APISuite) == 1 && !viper.GetBool(config.ForceSpecList) {
		// If there is only one specification loaded, then hotwire '/' to redirect to the
		// specification summary page unless DapperDox is configured to show the specification list page.
		r.Path("/").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/"+specification.ID+"/reference", http.StatusFound)
		})
	} else {
		r.Path("/").Methods(http.MethodGet).HandlerFunc(specificationListHandler)
	}
}

func specificationListHandler(w http.ResponseWriter, req *http.Request) {
	log().Trace("Render HTML for top level index page")

	render.HTML(w, http.StatusOK, "specification_list",
		render.DefaultVars(req, nil, render.Vars{"Title": "Specifications list", "SpecificationList": true}))
}

func specificationSummaryHandler(s *spec.APISpecification) func(w http.ResponseWriter, req *http.Request) {

	// The default "theme" level reference index page.
	tmpl := "specification_summary"

	customTmpl := s.ID + "/specification_summary"

	log().Tracef("+ Test for template %q", customTmpl)

	if render.TemplateLookup(customTmpl) != nil {
		tmpl = customTmpl
	}
	return func(w http.ResponseWriter, req *http.Request) {
		render.HTML(w, http.StatusOK, tmpl,
			render.DefaultVars(req, s, render.Vars{"Title": "Specification summary", "SpecificationSummary": true}))
	}
}
