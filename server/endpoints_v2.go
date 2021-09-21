// Copyright 2020, 2021 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"net/http"
	"path/filepath"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// ReportEndpointV2 https://issues.redhat.com/browse/CCXDEV-5097
	ReportEndpointV2 = "cluster/{cluster}/reports"

	// ClustersDetail https://issues.redhat.com/browse/CCXDEV-5088
	ClustersDetail = "rule/{rule_id}/clusters_detail/"

	// RuleContentV2 https://issues.redhat.com/browse/CCXDEV-5094
	// additionally group info is added too
	// https://github.com/RedHatInsights/insights-results-smart-proxy/pull/604
	RuleContentV2 = "rule/{rule_id}"

	// ContentV2 returns all the static content avaiable for the user
	ContentV2 = "content"
)

// addV2EndpointsToRouter adds API V2 specific endpoints to the router
func (server *HTTPServer) addV2EndpointsToRouter(router *mux.Router) {
	apiV2Prefix := server.Config.APIv2Prefix
	openAPIv2URL := apiV2Prefix + filepath.Base(server.Config.APIv2SpecFile)
	aggregatorBaseEndpoint := server.ServicesConfig.AggregatorBaseEndpoint

	// Common REST API endpoints
	router.HandleFunc(apiV2Prefix+MainEndpoint, server.mainEndpoint).Methods(http.MethodGet)

	// Reports endpoints
	server.addV2ReportsEndpointsToRouter(router, apiV2Prefix, aggregatorBaseEndpoint)

	// Content related endpoints
	server.addV2ContentEndpointsToRouter(router, apiV2Prefix)

	// Rules related endpoints
	server.addV2RuleEndpointsToRouter(router, apiV2Prefix, aggregatorBaseEndpoint)

	// Prometheus metrics
	router.Handle(apiV2Prefix+MetricsEndpoint, promhttp.Handler()).Methods(http.MethodGet)

	// OpenAPI specs
	router.HandleFunc(
		openAPIv2URL,
		httputils.CreateOpenAPIHandler(server.Config.APIv2SpecFile, server.Config.Debug, true),
	).Methods(http.MethodGet)
}

// addV2ReportsEndpointsToRouter method registers handlers for endpoints that
// return cluster report or reports to client
func (server *HTTPServer) addV2ReportsEndpointsToRouter(router *mux.Router, apiPrefix string, aggregatorBaseURL string) {
	router.HandleFunc(apiPrefix+ReportEndpointV2, server.reportEndpoint).Methods(http.MethodGet, http.MethodOptions)
}

// addV2RuleEndpointsToRouter method registers handlers for endpoints that handle
// rule-related operations (voting etc.)
func (server *HTTPServer) addV2RuleEndpointsToRouter(router *mux.Router, apiPrefix string, aggregatorBaseEndpoint string) {
	return
}

// addV2ContentEndpointsToRouter method registers handlers for endpoints that
// returns content to clients
func (server HTTPServer) addV2ContentEndpointsToRouter(router *mux.Router, apiPrefix string) {
	router.HandleFunc(apiPrefix+RuleContentV2, server.getContentWithGroupsForRule).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+ContentV2, server.getContentWithGroups).Methods(http.MethodGet)
}