package club

import "regexp"

// TODO: move this into some config file
var simpleClubMappings = map[string]string{
	"braunschweig":             "eintracht braunschweig",
	"bochum":                   "vfl bochum",
	"hertha berlin":            "hertha bsc",
	"hannover":                 "hannover 96",
	"freiburg":                 "sc freiburg",
	"bayer leverkusen":         "bayer 04 leverkusen",
	"b. monchengladbach":       "borussia moenchengladbach",
	"b. mönchengladbach":       "borussia moenchengladbach",
	"borussia mönchengladbach": "borussia moenchengladbach",
	"málaga cf":                "malaga cf",
	"málaga":                   "malaga cf",
	"malaga":                   "malaga cf",
	"getafe":                   "getafe cf",
	"levante":                  "levante ud",
	"almeria":                  "ud almeria",
	"sevilla":                  "sevilla fc",
	"valencia":                 "valencia cf",
	"benfica":                  "sl benfica",
	"desportivo aves":          "cd aves",
	"feirense":                 "cd feirense",
	"hoffenheim":               "tsg 1899 hoffenheim",
	"wolfsburg":                "vfl wolfsburg",
	"augsburg":                 "fc augsburg",
	"milan":                    "ac milan",
	"eibar":                    "sd eibar",
	"atlético madrid":          "atletico madrid",
	"atalanta":                 "atalanta bc",
	"sampdoria":                "uc sampdoria",
	"barcelona":                "fc barcelona",
	"inter milano":             "inter mailand",
	"bayern münchen":           "bayern munchen",
	"beşiktaş":                 "besiktas",
	"bayern munich":            "bayern munchen",
	"brighton & hove albion":   "brighton and hove albion",
	"udinese":                  "udinese calcio",
	"sassuolo":                 "us sassuolo calcio",
	"juventus":                 "juventus fc",
}

var regexClubMappings = map[string]func([]string) string{
	// Maps Under-X teams to their parent clubs
	`^(.*) u[0-9]{2}$`: func(m []string) string { return m[1] },
}

// NormaliseName attempts to map club names to a common name
func NormaliseName(club string) string {
	if mapping, ok := simpleClubMappings[club]; ok {
		return mapping
	}
	for reStr, res := range regexClubMappings {
		re := regexp.MustCompile(reStr)
		if groups := re.FindStringSubmatch(club); len(groups) > 0 {
			return res(groups)
		}
	}
	return club
}
