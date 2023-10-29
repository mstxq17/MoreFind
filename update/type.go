package update

type AssetFormat uint

const (
	Zip AssetFormat = iota
	Tar
	Unknown
)

// FileExtension of this asset format
func (a AssetFormat) FileExtension() string {
	if a == Zip {
		return ".zip"
	} else if a == Tar {
		return ".tar.gz"
	}
	return ""
}
