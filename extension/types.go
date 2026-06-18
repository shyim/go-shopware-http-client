package extension

import (
	"encoding/json"
	"fmt"
)

// ExtensionList is the response of ListAvailable.
type ExtensionList []*ExtensionDetail

// GetByName returns the extension with the given technical name, or nil.
func (l ExtensionList) GetByName(name string) *ExtensionDetail {
	for _, detail := range l {
		if detail.Name == name {
			return detail
		}
	}
	return nil
}

// FilterByUpdatable returns only the extensions that have a newer version
// available.
func (l ExtensionList) FilterByUpdatable() ExtensionList {
	out := make(ExtensionList, 0)
	for _, detail := range l {
		if detail.IsUpdatable() {
			out = append(out, detail)
		}
	}
	return out
}

// ExtensionDetail describes a single extension. Fields whose shape is not
// stable across Shopware versions are kept as json.RawMessage; decode them
// yourself when needed.
type ExtensionDetail struct {
	Name             string `json:"name"`
	Label            string `json:"label"`
	Description      string `json:"description"`
	ProducerName     string `json:"producerName"`
	License          string `json:"license"`
	Version          string `json:"version"`
	LatestVersion    string `json:"latestVersion"`
	NumberOfRatings  int    `json:"numberOfRatings"`
	LocalID          string `json:"localId"`
	Active           bool   `json:"active"`
	Type             string `json:"type"`
	IsTheme          bool   `json:"isTheme"`
	Configurable     bool   `json:"configurable"`
	Source           string `json:"source"`
	UpdateSource     string `json:"updateSource"`
	IconRaw          *string `json:"iconRaw"`

	InstalledAt *ExtensionDate `json:"installedAt"`
	UpdatedAt   *ExtensionDate `json:"updatedAt"`

	// Permissions, Images, Categories, etc. vary by version/source.
	Permissions json.RawMessage `json:"permissions,omitempty"`
	Images      json.RawMessage `json:"images,omitempty"`
	Categories  json.RawMessage `json:"categories,omitempty"`
}

// ExtensionDate is the PHP DateTime envelope Shopware serializes timestamps as.
type ExtensionDate struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

// IsPlugin reports whether the extension is a plugin (vs an app).
func (e ExtensionDetail) IsPlugin() bool {
	return e.Type == "plugin"
}

// IsUpdatable reports whether a newer version than the installed one exists.
func (e ExtensionDetail) IsUpdatable() bool {
	return e.LatestVersion != "" && e.LatestVersion != e.Version
}

// Status returns a human-readable status line for the extension.
func (e ExtensionDetail) Status() string {
	var text string
	switch {
	case e.Source == "store":
		text = "can be downloaded from store"
	case e.Active:
		text = "installed, activated"
	case e.InstalledAt != nil:
		text = "installed, not activated"
	default:
		text = "not installed, not activated"
	}

	if e.IsUpdatable() {
		text = fmt.Sprintf("%s, update available to %s", text, e.LatestVersion)
	}
	return text
}
