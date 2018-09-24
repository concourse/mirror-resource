package resource

import "os"

const DefaultInitialVersion = "mock"

type Source struct {
	// fetch the resource itself as an image
	MirrorSelf bool `json:"mirror_self"`

	// initial version that the resource should emit from /check
	// (default: 'mock')
	RawInitialVersion string `json:"initial_version"`

	// don't emit an initial version; useful for testing pipeline triggering
	NoInitialVersion bool `json:"no_initial_version"`

	// version to emit regardless of any version specified during check
	ForceVersion string `json:"force_version"`

	// a map of file paths to create with the associated contents
	//
	// contents can either be a string or an arbitrary object (which will be
	// JSON-marshalled)
	CreateFiles map[string]interface{} `json:"create_files"`

	// an amount of time (in Go duration format) to sleep before the check
	// returns versions
	CheckDelay string `json:"check_delay"`

	// force checking to fail with this message on stderr
	CheckFailure string `json:"check_failure"`
}

func (s Source) InitialVersion() string {
	if s.RawInitialVersion != "" {
		return s.RawInitialVersion
	}

	fromEnv := os.Getenv("MIRRORED_VERSION")
	if fromEnv != "" {
		return fromEnv
	}

	return DefaultInitialVersion
}

type Version struct {
	Version string `json:"version"`
}

type GetParams struct {
	// same as configuring mirror_self in source, but in params
	MirrorSelfViaParams bool `json:"mirror_self_via_params"`

	// similar to create_files in source; merged in so that additional (or
	// replaced) files can be specified
	CreateFiles map[string]interface{} `json:"create_files_via_params"`
}

type PutParams struct {
	// version to "create"
	Version string `json:"version"`
}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
