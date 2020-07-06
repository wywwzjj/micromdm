/*
	Package version provides utilities for displaying version information about a Go application.
	To use this package, a program would set the package variables at build time, using the
	-ldflags go build flag.
	Example:
		go build -ldflags "-X micromdm.io/v2/pkg/version.version=1.0.0"
	Available values and defaults to use with ldflags:
		version   = "unknown"
		branch    = "unknown"
		revision  = "unknown"
		buildDate = "unknown"
		buildUser = "unknown"
		appName   = "unknown"
*/
package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	goversion "rsc.io/goversion/version"
)

var (
	version   = "unknown"
	branch    = "unknown"
	revision  = "unknown"
	buildDate = "unknown"
	buildUser = "unknown"
	appName   = "micromdm"
)

// Info holds version and build info about the program.
type Info struct {
	Version   string `json:"version"`
	Branch    string `json:"branch"`
	Revision  string `json:"revision"`
	BuildDate string `json:"build_date"`
	BuildUser string `json:"build_user"`
}

// Version returns a struct with the current version information.
func Version() Info {
	return Info{
		Version:   version,
		Branch:    branch,
		Revision:  revision,
		BuildDate: buildDate,
		BuildUser: buildUser,
	}
}

// Print outputs the app name and version string.
func Print() {
	v := Version()
	fmt.Printf("%s version %s\n", appName, v.Version)
}

// PrintFull outputs the app name and detailed version information.
func PrintFull() error {
	v := Version()
	fmt.Printf("%s - version %s\n", appName, v.Version)
	fmt.Printf("branch: \t%s\n", v.Branch)
	fmt.Printf("revision: \t%s\n", v.Revision)
	fmt.Printf("build date: \t%s\n", v.BuildDate)
	fmt.Printf("build user: \t%s\n", v.BuildUser)

	binary, err := os.Executable()
	if err != nil {
		return err
	}

	binVersion, err := goversion.ReadExe(binary)
	if err != nil {
		return err
	}
	fmt.Printf("go release: \t%s\n", binVersion.Release)
	fmt.Println()
	fmt.Println(binVersion.ModuleInfo)
	return nil
}

// Handler provides an HTTP Handler which returns JSON formatted version info.
func Handler() http.Handler {
	v := Version()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(v)
	})
}
