package version

import (
	"encoding/json"
	"fmt"
)

var (
	// BuildDate is date when binary was built
	BuildDate string
	// BuildVersion is the version of binary
	BuildVersion string
)

// providerVersion holds current provider version
type providerVersion struct {
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
}

func PrintVersion() (err error) {
	pv := providerVersion{
		Version:   BuildVersion,
		BuildDate: BuildDate,
	}

	var res []byte
	if res, err = json.Marshal(pv); err != nil {
		return
	}

	fmt.Printf(string(res) + "\n")
	return
}
