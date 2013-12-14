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
package dockercloud

import (
	"os"
)

// The Cloud interface provides the contract that cloud providers should implement to enable
// running Docker containers in their cloud.
// TODO(bburns): Restructure this into Cloud, Instance and Tunnel interfaces
type Cloud interface {
	// GetPublicIPAddress returns the stringified address (e.g "1.2.3.4") of the runtime
	GetPublicIPAddress(name string, zone string) (string, error)

	// CreateInstance creates a virtual machine instance given a name and a zone.  Returns the
	// IP address of the instance.  Waits until Docker is up and functioning on the machine
	// before returning.
	CreateInstance(name string, zone string) (string, error)

	// DeleteInstance deletes a virtual machine instance, given the instance name and zone.
	DeleteInstance(name string, zone string) error

	// Open a secure tunnel (generally SSH) between the local host and a remote host.
	OpenSecureTunnel(name string, zone string, localPort int, remotePort int) (*os.Process, error)
}
