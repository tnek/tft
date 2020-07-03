package riot

var (
	// RegionToPlatform is the mapping of Region name to Platform
	RegionToPlatform = map[string][]string{
		"americas": {"br1", "la1", "la2", "na1"},
		"asia":     {"jp1", "kr", "oc1"},
		"europe":   {"eun1", "euw1", "tr1", "ru"},
	}

	// PlatformToRegion is the reverse mapping of Platform name to corresponding Region
	// name (e.x. "na1" -> "americas")
	PlatformToRegion map[string]string
)

func init() {
	PlatformToRegion = make(map[string]string)
	for reg, plats := range RegionToPlatform {
		for _, plat := range plats {
			PlatformToRegion[plat] = reg
		}
	}
}
