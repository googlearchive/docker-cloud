//
// Copyright (C) 2013 The Docker Cloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/googlecloudplatform/docker-cloud/dockercloud"
	"github.com/rakyll/command"
)

// Try to connect to a tunnel to the docker dameon if it exists.
// url is the URL to test.
// returns true, if the connection was successful, false otherwise
type Tunnel struct {
	url.URL
}

func (t Tunnel) isActive() bool {
	_, err := http.Get(t.String())
	return err == nil
}

type ProxyServer struct {
	// TODO(jbd): instance name and zone belongs to GCECloud.
	instanceName string
	zone         string

	tunnelPort int
	dockerPort int

	cloud dockercloud.Cloud
}

func (server ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := server.doServe(w, r)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "{'error': '%s'}", err)
	}
}

func (server ProxyServer) doServe(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ip string
	path := r.URL.Path
	query := r.URL.RawQuery
	host := fmt.Sprintf("localhost:%d", server.tunnelPort)
	targetUrl := fmt.Sprintf("http://%s%s?%s", host, path, query)

	w.Header().Add("Content-Type", "application/json")

	// Try to find a VM instance.
	ip, err = server.cloud.GetPublicIPAddress(server.instanceName, server.zone)
	instanceRunning := len(ip) > 0
	// err is 404 if the instance doesn't exist, so we only error out when
	// instanceRunning is true.
	if err != nil && instanceRunning {
		return err
	}

	// If there's no VM instance, and the request is 'ps' just return []
	if r.Method == "GET" && strings.HasSuffix(path, "/containers/json") && !instanceRunning {
		w.WriteHeader(200)
		fmt.Fprintf(w, "[]")
		return nil
	}

	// Otherwise create a new VM.
	if !instanceRunning {
		ip, err = server.cloud.CreateInstance(server.instanceName, server.zone)
		if err != nil {
			return err
		}
	}

	// Test for the SSH tunnel, create if it doesn't exist.
	tunnelUrl, err := url.Parse("http://" + host + "/v1.6/containers/json")
	if err != nil {
		return err
	}
	tunnel := Tunnel{*tunnelUrl}

	if !tunnel.isActive() {
		fmt.Printf("Creating tunnel")
		_, err = server.cloud.OpenSecureTunnel(
			server.instanceName, server.zone, server.tunnelPort, server.dockerPort)
		if err != nil {
			return err
		}
	}

	err = proxyRequest(targetUrl, r, w)
	if err != nil {
		return err
	}
	if strings.HasSuffix(path, "/stop") {
		server.maybeDelete(host, server.instanceName, server.zone)
	}
	return nil
}

func proxyRequest(url string, r *http.Request, w http.ResponseWriter) error {
	var res *http.Response
	var err error

	// Proxy the request.
	if r.Method == "GET" {
		res, err = http.Get(url)
	}
	if r.Method == "POST" {
		res, err = http.Post(url, "application/json", r.Body)
	}
	if err != nil {
		return err
	}
	w.WriteHeader(res.StatusCode)
	defer res.Body.Close()
	// TODO(bburns) : Intercept 'ps' here and substitute in the ip address.
	_, err = io.Copy(w, res.Body)
	return err
}

// TODO(bburns) : clone this from docker somehow?
type ContainerPort struct {
	PrivatePort float64
	PublicPort  float64
	Type        string
}

type ContainerStatus struct {
	Id         string
	Image      string
	Command    string
	Created    float64
	Status     string
	Ports      []ContainerPort
	SizeRW     float64
	SizeRootFs float64
}

func (server ProxyServer) maybeDelete(host string, instanceName string, zone string) error {
	res, err := http.Get(fmt.Sprintf("http://%s/v1.6/containers/json", host))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf(string(body))
	var containers []ContainerStatus
	err = json.Unmarshal(body, &containers)
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		err = server.cloud.DeleteInstance(instanceName, zone)
		if err != nil {
			return err
		}
	}
	return nil
}

// Starts the authorization wizard to retrieve
// and OAuth 2.0 access and refresh token.
type authCmd struct {
	clientId, clientSecret, scope, projectId *string
}

// Defines auth subcommand related flags.
func (cmd *authCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
	// Authorization flags.
	cmd.clientId = fs.String("id", "", "Client id")
	cmd.clientSecret = fs.String("secret", "", "Client Secret")
	cmd.scope = fs.String("scope", "", "OAuth 2 Scope")
	cmd.projectId = fs.String("project", "", "Google Cloud Project Name")
	return fs
}

// Handles the auth command.
func (cmd *authCmd) Run(args []string) {
	if err := dockercloud.ConfigureGCE(*cmd.clientId, *cmd.clientSecret, *cmd.scope, *cmd.projectId); err != nil {
		log.Fatal(err)
	}
}

// Starts a Docker proxy server to intecept the Docker commands
// to create and delete VMs on demand.
type startCmd struct {
	proxyPort    *int
	dockerPort   *int
	tunnelPort   *int
	instanceName *string
	zone         *string
	projectId    *string
}

// Defines the flags required by start subcommand.
func (cmd *startCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
	cmd.proxyPort = fs.Int("port", 8080, "The local port to run on.")
	cmd.dockerPort = fs.Int("dockerport", 8000, "The remote port to run docker on")
	cmd.tunnelPort = fs.Int("tunnelport", 8001, "The local port open the tunnel to docker")
	cmd.instanceName = fs.String("instancename", "docker-instance", "The name of the instance")
	cmd.zone = fs.String("zone", "us-central1-a", "The zone to run in")
	cmd.projectId = fs.String("project", "", "Google Cloud Project Name")
	return fs
}

// Handles the start command.
func (cmd *startCmd) Run(args []string) {
	gce, err := dockercloud.NewCloudGCE(*cmd.projectId)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", ProxyServer{
		instanceName: *cmd.instanceName,
		zone:         *cmd.zone,
		dockerPort:   *cmd.dockerPort,
		tunnelPort:   *cmd.tunnelPort,
		cloud:        gce,
	})
	addr := fmt.Sprintf(":%d", *cmd.proxyPort)
	log.Printf("Server started, now you can use docker -H tcp://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	// Registers authCmd to handle `program [args...] auth [subcommand-args...]`
	command.On("auth", "Allow you to authorize and configure project settings.", &authCmd{}, []string{"project"})
	// Registers startCmd to handle `program [args...] start [subcommand-args...]`
	command.On("start", "Starts the proxy server.", &startCmd{}, []string{})
	command.ParseAndRun()
}
