/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudinit

import (
	"fmt"

	"github.com/zawachte-msft/cluster-api-k3s/pkg/secret"
)

const (
	controlPlaneCloudInit = `{{.Header}}
{{template "files" .WriteFiles}}
runcmd:
{{- template "commands" .PreK3sCommands }}
  - 'curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=%s sh -s - server'
{{- template "commands" .PostK3sCommands }}
users:
  - name: capv
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: false
    ssh_authorized_keys:
    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCFi51uN6SZvK0dI86wGNaZzOrYaJlFCXlXMiPM5p63kvuRM1KpdPctfxpIkeV8o9F6sOMaof0pD0dnf7WatYM0ztxKa+ZEm6zm1Zlz7/xPUv/NOn7s6AoG3Fm9+ywOJI8SUB/E0KzQzhIengds7IIC/ySRVzlI9xCVa1N2/Rd/C8W/hEvtPS4h6BgM/8b17QDhp0uHVxPYicXRGEIDFKCQxZrSLVu4Q56EETUkSX0tQTBsMwzOJ5PYoJvTkzw/p6LcaysXwMunkFW5jVxOLcm5vM/IqnKbTfrlkuR2va3SBn7/LP/H3790hz88D+TnPyIWHMDKw79wjwd3k5muyfz root@hcm
    groups: users
    hashed_passwd: $1$xyz$4jfsyCqeX7YjcodykY1rB1
`
)

// ControlPlaneInput defines the context to generate a controlplane instance user data.
type ControlPlaneInput struct {
	BaseUserData
	secret.Certificates
}

// NewInitControlPlane returns the user data string to be used on a controlplane instance.
func NewInitControlPlane(input *ControlPlaneInput) ([]byte, error) {
	input.Header = cloudConfigHeader
	input.WriteFiles = input.Certificates.AsFiles()
	input.WriteFiles = append(input.WriteFiles, input.AdditionalFiles...)
	input.WriteFiles = append(input.WriteFiles, input.ConfigFile)

	controlPlaneCloudJoinWithVersion := fmt.Sprintf(controlPlaneCloudInit, input.K3sVersion)
	userData, err := generate("InitControlplane", controlPlaneCloudJoinWithVersion, input)
	if err != nil {
		return nil, err
	}

	return userData, nil
}
