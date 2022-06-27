package usage

/*
OperatingSystem returns AWS Pricing OperatingSystem from Usage Platform.
*/
var OperatingSystem = map[string]string{
	"Amazon Linux":                "Linux",   // PreInstalled: NA
	"Linux/UNIX":                  "Linux",   // PreInstalled: NA
	"Linux/UNIX (Amazon VPC)":     "Linux",   // PreInstalled: NA
	"Linux with SQL Standard":     "Linux",   // PreInstalled: SQL Standard
	"Linux with SQL Web":          "Linux",   // PreInstalled: SQL Web
	"Linux with SQL Enterprise":   "Linux",   // PreInstalled: SQL Enterprise
	"Red Hat Enterprise Linux":    "RHEL",    // PreInstalled: NA
	"SUSE Linux":                  "SUSE",    // PreInstalled: NA
	"Windows":                     "Windows", // PreInstalled: NA
	"Windows (Amazon VPC)":        "Windows", // PreInstalled: NA
	"Windows with SQL Standard":   "Windows", // PreInstalled: SQL Standard
	"Windows with SQL Web":        "Windows", // PreInstalled: SQL Web
	"Windows with SQL Enterprise": "Windows", // PreInstalled: SQL Enterprise
	"Windows (BYOL)":              "",        // pricing not found
	"NoOperatingSystem":           "",        // pricing not found
}
