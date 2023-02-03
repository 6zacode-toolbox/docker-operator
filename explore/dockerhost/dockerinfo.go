package type

type DockerInfo struct {
	ID                 string           `json:"ID"`
	Containers         int              `json:"Containers"`
	ContainersRunning  int              `json:"ContainersRunning"`
	ContainersPaused   int              `json:"ContainersPaused"`
	ContainersStopped  int              `json:"ContainersStopped"`
	Images             int              `json:"Images"`
	Driver             string           `json:"Driver"`
	Plugins            Plugins          `json:"Plugins"`
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
	GenericResources   interface{}      `json:"GenericResources"`
	DockerRootDir      string           `json:"DockerRootDir"`
	HTTPProxy          string           `json:"HttpProxy"`
	HTTPSProxy         string           `json:"HttpsProxy"`
	NoProxy            string           `json:"NoProxy"`
	Name               string           `json:"Name"`
	Labels             []interface{}    `json:"Labels"`
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
	Warnings           interface{}      `json:"Warnings"`
	ClientInfo         ClientInfo       `json:"ClientInfo"`
}
type Plugins struct {
	Volume        []string    `json:"Volume"`
	Network       []string    `json:"Network"`
	Authorization interface{} `json:"Authorization"`
	Log           []string    `json:"Log"`
}
type DockerIo struct {
	Name     string        `json:"Name"`
	Mirrors  []interface{} `json:"Mirrors"`
	Secure   bool          `json:"Secure"`
	Official bool          `json:"Official"`
}
type HubproxyDockerInternal5000 struct {
	Name     string        `json:"Name"`
	Mirrors  []interface{} `json:"Mirrors"`
	Secure   bool          `json:"Secure"`
	Official bool          `json:"Official"`
}
type IndexConfigs struct {
	DockerIo                   DockerIo                   `json:"docker.io"`
	HubproxyDockerInternal5000 HubproxyDockerInternal5000 `json:"hubproxy.docker.internal:5000"`
}
type RegistryConfig struct {
	AllowNondistributableArtifactsCIDRs     []interface{} `json:"AllowNondistributableArtifactsCIDRs"`
	AllowNondistributableArtifactsHostnames []interface{} `json:"AllowNondistributableArtifactsHostnames"`
	InsecureRegistryCIDRs                   []string      `json:"InsecureRegistryCIDRs"`
	IndexConfigs                            IndexConfigs  `json:"IndexConfigs"`
	Mirrors                                 []interface{} `json:"Mirrors"`
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
	NodeID           string      `json:"NodeID"`
	NodeAddr         string      `json:"NodeAddr"`
	LocalNodeState   string      `json:"LocalNodeState"`
	ControlAvailable bool        `json:"ControlAvailable"`
	Error            string      `json:"Error"`
	RemoteManagers   interface{} `json:"RemoteManagers"`
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