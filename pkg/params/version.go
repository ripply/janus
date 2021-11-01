package params

import "fmt"

const (
	unstable     = "unstable"
	stable       = "stable"
	VersionMajor = 0        // Major version component of the current release
	VersionMinor = 1        // Minor version component of the current release
	VersionPatch = 0        // Patch version component of the current release
	VersionMeta  = unstable // Version metadata to append to the version string
)

var GitSha = "development"

var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

var VersionWithGitSha = func() string {
	return fmt.Sprintf("%s-%s", Version, GitSha)
}()
