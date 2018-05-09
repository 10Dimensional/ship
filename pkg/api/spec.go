package api

// Spec is the top level Ship document that defines an application
type Spec struct {
	Assets    Assets    `json:"assets" yaml:"assets" hcl:"asset"`
	Lifecycle Lifecycle `json:"lifecycle" yaml:"lifecycle" hcl:"lifecycle"`
	Config    Config    `json:"config" yaml:"config" hcl:"config"`
}

// Image
type Image struct {
	URL      string `json:"url" yaml:"url" hcl:"url" meta:"url"`
	Source   string `json:"source" yaml:"source" hcl:"source" meta:"source"`
	AppSlug  string `json:"appSlug" yaml:"appSlug" hcl:"appSlug" meta:"appSlug"`
	ImageKey string `json:"imageKey" yaml:"imageKey" hcl:"imageKey" meta:"imageKey"`
}

// ReleaseMetadata
type ReleaseMetadata struct {
	CustomerID     string  `json:"customerId" yaml:"customerId" hcl:"customerId" meta:"customer-id"`
	ChannelID      string  `json:"channelId" yaml:"channelId" hcl:"channelId" meta:"channel-id"`
	ChannelName    string  `json:"channelName" yaml:"channelName" hcl:"channelName" meta:"channel-name"`
	ChannelIcon    string  `json:"channelIcon" yaml:"channelIcon" hcl:"channelIcon" meta:"channel-icon"`
	Semver         string  `json:"semver" yaml:"semver" hcl:"semver" meta:"release-version"`
	ReleaseNotes   string  `json:"releaseNotes" yaml:"releaseNotes" hcl:"releaseNotes" meta:"release-notes"`
	Created        string  `json:"created" yaml:"created" hcl:"created" meta:"release-date"`
	RegistrySecret string  `json:"registrySecret" yaml:"registrySecret" hcl:"registrySecret" meta:"registry-secret"`
	Images         []Image `json:"images" yaml:"images" hcl:"images" meta:"images"`
}

// Release
type Release struct {
	Metadata ReleaseMetadata
	Spec     Spec
}
