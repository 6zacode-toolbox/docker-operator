/*
Copyright 2023 6zacode-toolbox.

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

package v1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DockerHostSpec defines the desired state of DockerHost
type DockerHostSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DockerHost. Edit dockerhost_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// DockerHostStatus defines the observed state of DockerHost
type DockerHostStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Containers []DockerContainerSummary `json:"containers,omitempty"`
	HostInfo   DockerInfo               `json:"host,omitempty`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DockerHost is the Schema for the dockerhosts API
type DockerHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DockerHostSpec   `json:"spec,omitempty"`
	Status DockerHostStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DockerHostList contains a list of DockerHost
type DockerHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DockerHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DockerHost{}, &DockerHostList{})
}

type HostPlugins struct {
	Volume        []string `json:"Volume"`
	Network       []string `json:"Network"`
	Authorization []string `json:"Authorization"`
	Log           []string `json:"Log"`
}

type DockerInfo struct {
	ID                 string           `json:"ID"`
	Containers         int              `json:"Containers"`
	ContainersRunning  int              `json:"ContainersRunning"`
	ContainersPaused   int              `json:"ContainersPaused"`
	ContainersStopped  int              `json:"ContainersStopped"`
	Images             int              `json:"Images"`
	Driver             string           `json:"Driver"`
	Plugins            HostPlugins      `json:"Plugins"`
	MemoryLimit        bool             `json:"MemoryLimit"`
	SwapLimit          bool             `json:"SwapLimit"`
	KernelMemory       bool             `json:"KernelMemory"`
	KernelMemoryTCP    bool             `json:"KernelMemoryTCP"`
	CPUCfsPeriod       bool             `json:"CpuCfsPeriod"`
	CPUCfsQuota        bool             `json:"CpuCfsQuota"`
	CPUShares          bool             `json:"CPUShares"`
	CPUSet             bool             `json:"CPUSet"`
	PidsLimit          bool             `json:"PidsLimit"`
	IPv4Forwarding     bool             `json:"IPv4Forwarding"`
	BridgeNfIptables   bool             `json:"BridgeNfIptables"`
	BridgeNfIP6Tables  bool             `json:"BridgeNfIp6tables"`
	Debug              bool             `json:"Debug"`
	NFd                int              `json:"NFd"`
	OomKillDisable     bool             `json:"OomKillDisable"`
	NGoroutines        int              `json:"NGoroutines"`
	SystemTime         time.Time        `json:"SystemTime"`
	LoggingDriver      string           `json:"LoggingDriver"`
	CgroupDriver       string           `json:"CgroupDriver"`
	CgroupVersion      string           `json:"CgroupVersion"`
	NEventsListener    int              `json:"NEventsListener"`
	KernelVersion      string           `json:"KernelVersion"`
	OperatingSystem    string           `json:"OperatingSystem"`
	OSVersion          string           `json:"OSVersion"`
	OSType             string           `json:"OSType"`
	Architecture       string           `json:"Architecture"`
	IndexServerAddress string           `json:"IndexServerAddress"`
	RegistryConfig     RegistryConfig   `json:"RegistryConfig"`
	Ncpu               int              `json:"NCPU"`
	MemTotal           int64            `json:"MemTotal"`
	DockerRootDir      string           `json:"DockerRootDir"`
	HTTPProxy          string           `json:"HttpProxy"`
	HTTPSProxy         string           `json:"HttpsProxy"`
	NoProxy            string           `json:"NoProxy"`
	Name               string           `json:"Name"`
	Labels             []string         `json:"Labels"`
	ExperimentalBuild  bool             `json:"ExperimentalBuild"`
	ServerVersion      string           `json:"ServerVersion"`
	Runtimes           Runtimes         `json:"Runtimes"`
	DefaultRuntime     string           `json:"DefaultRuntime"`
	Swarm              Swarm            `json:"Swarm"`
	LiveRestoreEnabled bool             `json:"LiveRestoreEnabled"`
	Isolation          string           `json:"Isolation"`
	InitBinary         string           `json:"InitBinary"`
	ContainerdCommit   ContainerdCommit `json:"ContainerdCommit"`
	RuncCommit         RuncCommit       `json:"RuncCommit"`
	InitCommit         InitCommit       `json:"InitCommit"`
	SecurityOptions    []string         `json:"SecurityOptions"`
	ClientInfo         ClientInfo       `json:"ClientInfo"`
}

type DockerIo struct {
	Name     string   `json:"Name"`
	Mirrors  []string `json:"Mirrors"`
	Secure   bool     `json:"Secure"`
	Official bool     `json:"Official"`
}
type HubproxyDockerInternal5000 struct {
	Name     string   `json:"Name"`
	Mirrors  []string `json:"Mirrors"`
	Secure   bool     `json:"Secure"`
	Official bool     `json:"Official"`
}
type IndexConfigs struct {
	DockerIo                   DockerIo                   `json:"docker.io"`
	HubproxyDockerInternal5000 HubproxyDockerInternal5000 `json:"hubproxy.docker.internal:5000"`
}
type RegistryConfig struct {
	AllowNondistributableArtifactsCIDRs     []string     `json:"AllowNondistributableArtifactsCIDRs"`
	AllowNondistributableArtifactsHostnames []string     `json:"AllowNondistributableArtifactsHostnames"`
	InsecureRegistryCIDRs                   []string     `json:"InsecureRegistryCIDRs"`
	IndexConfigs                            IndexConfigs `json:"IndexConfigs"`
	Mirrors                                 []string     `json:"Mirrors"`
}
type IoContainerdRuncV2 struct {
	Path string `json:"path"`
}
type IoContainerdRuntimeV1Linux struct {
	Path string `json:"path"`
}
type Runc struct {
	Path string `json:"path"`
}
type Runtimes struct {
	IoContainerdRuncV2         IoContainerdRuncV2         `json:"io.containerd.runc.v2"`
	IoContainerdRuntimeV1Linux IoContainerdRuntimeV1Linux `json:"io.containerd.runtime.v1.linux"`
	Runc                       Runc                       `json:"runc"`
}
type Swarm struct {
	NodeID           string   `json:"NodeID"`
	NodeAddr         string   `json:"NodeAddr"`
	LocalNodeState   string   `json:"LocalNodeState"`
	ControlAvailable bool     `json:"ControlAvailable"`
	Error            string   `json:"Error"`
	RemoteManagers   []string `json:"RemoteManagers"`
}
type ContainerdCommit struct {
	ID       string `json:"ID"`
	Expected string `json:"Expected"`
}
type RuncCommit struct {
	ID       string `json:"ID"`
	Expected string `json:"Expected"`
}
type InitCommit struct {
	ID       string `json:"ID"`
	Expected string `json:"Expected"`
}
type Plugins struct {
	SchemaVersion    string `json:"SchemaVersion"`
	Vendor           string `json:"Vendor"`
	Version          string `json:"Version"`
	ShortDescription string `json:"ShortDescription"`
	Name             string `json:"Name"`
	Path             string `json:"Path"`
	URL              string `json:"URL,omitempty"`
}
type ClientInfo struct {
	Debug   bool      `json:"Debug"`
	Context string    `json:"Context"`
	Plugins []Plugins `json:"Plugins"`
}

type DockerContainerSummary struct {
	Command      string `json:"Command"`
	CreatedAt    string `json:"CreatedAt"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Labels       string `json:"Labels"`
	LocalVolumes string `json:"LocalVolumes"`
	Mounts       string `json:"Mounts"`
	Names        string `json:"Names"`
	Networks     string `json:"Networks"`
	Ports        string `json:"Ports"`
	RunningFor   string `json:"RunningFor"`
	Size         string `json:"Size"`
	State        string `json:"State"`
	Status       string `json:"Status"`
}
