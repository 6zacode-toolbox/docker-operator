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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type HostPlugins struct {
	Volume        []string `json:"Volume"`
	Network       []string `json:"Network"`
	Authorization []string `json:"Authorization"`
	Log           []string `json:"Log"`
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

type DockerInfo struct {
	ID                 string           `json:"ID"`
	Containers         int              `json:"Containers,omitempty"`
	ContainersRunning  int              `json:"ContainersRunning,omitempty"`
	ContainersPaused   int              `json:"ContainersPaused,omitempty"`
	ContainersStopped  int              `json:"ContainersStopped,omitempty"`
	Images             int              `json:"Images,omitempty"`
	Driver             string           `json:"Driver,omitempty"`
	Plugins            HostPlugins      `json:"Plugins,omitempty"`
	MemoryLimit        bool             `json:"MemoryLimit,omitempty"`
	SwapLimit          bool             `json:"SwapLimit,omitempty"`
	KernelMemory       bool             `json:"KernelMemory,omitempty"`
	KernelMemoryTCP    bool             `json:"KernelMemoryTCP,omitempty"`
	CPUCfsPeriod       bool             `json:"CpuCfsPeriod,omitempty"`
	CPUCfsQuota        bool             `json:"CpuCfsQuota,omitempty"`
	CPUShares          bool             `json:"CPUShares,omitempty"`
	CPUSet             bool             `json:"CPUSet,omitempty"`
	PidsLimit          bool             `json:"PidsLimit,omitempty"`
	IPv4Forwarding     bool             `json:"IPv4Forwarding,omitempty"`
	BridgeNfIptables   bool             `json:"BridgeNfIptables,omitempty"`
	BridgeNfIP6Tables  bool             `json:"BridgeNfIp6tables,omitempty"`
	Debug              bool             `json:"Debug,omitempty"`
	NFd                int              `json:"NFd,omitempty"`
	OomKillDisable     bool             `json:"OomKillDisable,omitempty"`
	NGoroutines        int              `json:"NGoroutines,omitempty"`
	LoggingDriver      string           `json:"LoggingDriver,omitempty"`
	CgroupDriver       string           `json:"CgroupDriver,omitempty"`
	CgroupVersion      string           `json:"CgroupVersion,omitempty"`
	NEventsListener    int              `json:"NEventsListener,omitempty"`
	KernelVersion      string           `json:"KernelVersion,omitempty"`
	OperatingSystem    string           `json:"OperatingSystem,omitempty"`
	OSVersion          string           `json:"OSVersion,omitempty"`
	OSType             string           `json:"OSType,omitempty"`
	Architecture       string           `json:"Architecture,omitempty"`
	IndexServerAddress string           `json:"IndexServerAddress,omitempty"`
	RegistryConfig     RegistryConfig   `json:"RegistryConfig,omitempty"`
	Ncpu               int              `json:"NCPU,omitempty"`
	MemTotal           int64            `json:"MemTotal,omitempty"`
	DockerRootDir      string           `json:"DockerRootDir,omitempty"`
	HTTPProxy          string           `json:"HttpProxy,omitempty"`
	HTTPSProxy         string           `json:"HttpsProxy,omitempty"`
	NoProxy            string           `json:"NoProxy,omitempty"`
	Name               string           `json:"Name,omitempty"`
	Labels             []string         `json:"Labels,omitempty"`
	ExperimentalBuild  bool             `json:"ExperimentalBuild,omitempty"`
	ServerVersion      string           `json:"ServerVersion,omitempty"`
	Runtimes           Runtimes         `json:"Runtimes,omitempty"`
	DefaultRuntime     string           `json:"DefaultRuntime,omitempty"`
	Swarm              Swarm            `json:"Swarm,omitempty"`
	LiveRestoreEnabled bool             `json:"LiveRestoreEnabled,omitempty"`
	Isolation          string           `json:"Isolation,omitempty"`
	InitBinary         string           `json:"InitBinary,omitempty"`
	ContainerdCommit   ContainerdCommit `json:"ContainerdCommit,omitempty"`
	RuncCommit         RuncCommit       `json:"RuncCommit,omitempty"`
	InitCommit         InitCommit       `json:"InitCommit,omitempty"`
	SecurityOptions    []string         `json:"SecurityOptions,omitempty"`
	ClientInfo         ClientInfo       `json:"ClientInfo,omitempty"`
	//Not supported
	//SystemTime         time.Time   `json:"SystemTime"`

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
	DockerHost DockerInfo               `json:"host,omitempty"`
	//Sample string `json:"sample,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// DockerHost is the Schema for the dockerhosts API
type DockerHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DockerHostSpec   `json:"spec,omitempty"`
	Status DockerHostStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// DockerHostList contains a list of DockerHost
type DockerHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DockerHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DockerHost{}, &DockerHostList{})
}
