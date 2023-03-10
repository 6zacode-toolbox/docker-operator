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
	Volume        []string `json:"Volume,omitempty"`
	Network       []string `json:"Network,omitempty"`
	Authorization []string `json:"Authorization,omitempty"`
	Log           []string `json:"Log,omitempty"`
}
type DockerIo struct {
	Name     string   `json:"Name,omitempty"`
	Mirrors  []string `json:"Mirrors,omitempty"`
	Secure   bool     `json:"Secure,omitempty"`
	Official bool     `json:"Official,omitempty"`
}
type HubproxyDockerInternal5000 struct {
	Name     string   `json:"Name,omitempty"`
	Mirrors  []string `json:"Mirrors,omitempty"`
	Secure   bool     `json:"Secure,omitempty"`
	Official bool     `json:"Official,omitempty"`
}
type IndexConfigs struct {
	DockerIo                   DockerIo                   `json:"docker.io,omitempty"`
	HubproxyDockerInternal5000 HubproxyDockerInternal5000 `json:"hubproxy.docker.internal:5000,omitempty"`
}

type RegistryConfig struct {
	AllowNondistributableArtifactsCIDRs     []string     `json:"AllowNondistributableArtifactsCIDRs,omitempty"`
	AllowNondistributableArtifactsHostnames []string     `json:"AllowNondistributableArtifactsHostnames,omitempty"`
	InsecureRegistryCIDRs                   []string     `json:"InsecureRegistryCIDRs,omitempty"`
	IndexConfigs                            IndexConfigs `json:"IndexConfigs,omitempty"`
	Mirrors                                 []string     `json:"Mirrors,omitempty"`
}

type IoContainerdRuncV2 struct {
	Path string `json:"path,omitempty"`
}
type IoContainerdRuntimeV1Linux struct {
	Path string `json:"path,omitempty"`
}
type Runc struct {
	Path string `json:"path,omitempty"`
}
type Runtimes struct {
	IoContainerdRuncV2         IoContainerdRuncV2         `json:"io.containerd.runc.v2,omitempty"`
	IoContainerdRuntimeV1Linux IoContainerdRuntimeV1Linux `json:"io.containerd.runtime.v1.linux,omitempty"`
	Runc                       Runc                       `json:"runc,omitempty"`
}

type Swarm struct {
	NodeID           string   `json:"NodeID,omitempty"`
	NodeAddr         string   `json:"NodeAddr,omitempty"`
	LocalNodeState   string   `json:"LocalNodeState,omitempty"`
	ControlAvailable bool     `json:"ControlAvailable,omitempty"`
	Error            string   `json:"Error,omitempty"`
	RemoteManagers   []string `json:"RemoteManagers,omitempty"`
}

type ContainerdCommit struct {
	ID       string `json:"ID,omitempty"`
	Expected string `json:"Expected,omitempty"`
}
type RuncCommit struct {
	ID       string `json:"ID,omitempty"`
	Expected string `json:"Expected,omitempty"`
}
type InitCommit struct {
	ID       string `json:"ID,omitempty"`
	Expected string `json:"Expected,omitempty"`
}

type Plugins struct {
	SchemaVersion    string `json:"SchemaVersion,omitempty"`
	Vendor           string `json:"Vendor,omitempty"`
	Version          string `json:"Version,omitempty"`
	ShortDescription string `json:"ShortDescription,omitempty"`
	Name             string `json:"Name,omitempty"`
	Path             string `json:"Path,omitempty"`
	URL              string `json:"URL,omitempty"`
}
type ClientInfo struct {
	Debug   bool      `json:"Debug,omitempty"`
	Context string    `json:"Context,omitempty"`
	Plugins []Plugins `json:"Plugins,omitempty"`
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
	Command      string `json:"Command,omitempty"`
	CreatedAt    string `json:"CreatedAt,omitempty"`
	ID           string `json:"ID,omitempty"`
	Image        string `json:"Image,omitempty"`
	Labels       string `json:"Labels,omitempty"`
	LocalVolumes string `json:"LocalVolumes,omitempty"`
	Mounts       string `json:"Mounts,omitempty"`
	Names        string `json:"Names,omitempty"`
	Networks     string `json:"Networks,omitempty"`
	Ports        string `json:"Ports,omitempty"`
	RunningFor   string `json:"RunningFor,omitempty"`
	Size         string `json:"Size,omitempty"`
	State        string `json:"State,omitempty"`
	Status       string `json:"Status,omitempty"`
}

const ConnectionModeSSH = "ssh"
const ConnectionModeTLS = "tls"

// DockerHostSpec defines the desired state of DockerHost
type DockerHostSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	HostIP string `json:"hostip,omitempty"`

	SSHConnection SSHConnection `json:"sshConnection,omitempty"`

	//Possible Improvements
	// - Frequency of check in seconds
	// - Secret Name
	// - Secret Namespace
	// - ...
}

type SSHConnection struct {
	SSHUser string `json:"sshUser"`
}

// DockerHostStatus defines the observed state of DockerHost
type DockerHostStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Containers        []DockerContainerSummary `json:"containers,omitempty"`
	DockerHost        DockerInfo               `json:"host,omitempty"`
	Instanced         bool                     `json:"instanced,omitempty"`
	SuccessValidation bool                     `json:"successValidation,omitempty"`
	Validated         bool                     `json:"validated,omitempty"`
	Error             string                   `json:"error,omitempty"`
	AppConfigMap      map[string]string        `json:"appConfigMap,omitempty"`
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
