// Deals with services related to the ROBLOX catalog.
package catalog

import (
	"encoding/json"
	"github.com/anaminus/rbxweb"
	"net/url"
)

// Used with the Category field of a Query.
const (
	CatFeatured     byte = 0
	CatAll          byte = 1
	CatCollectibles byte = 2
	CatClothing     byte = 3
	CatBodyParts    byte = 4
	CatGear         byte = 5
	CatModels       byte = 6
	CatPlugins      byte = 7
	CatDecals       byte = 8
	CatAudio        byte = 9
)

// Used with the Subcategory field of a Query.
const (
	SubcatFeatured      byte = 0
	SubcatAll           byte = 1
	SubcatCollectibles  byte = 2
	SubcatClothing      byte = 3
	SubcatBodyParts     byte = 4
	SubcatGear          byte = 5
	SubcatModels        byte = 6
	SubcatDecals        byte = 7
	SubcatHats          byte = 8
	SubcatFaces         byte = 9
	SubcatPackages      byte = 10
	SubcatShirts        byte = 11
	SubcatTshirts       byte = 12
	SubcatPants         byte = 13
	SubcatHeads         byte = 14
	SubcatAudio         byte = 15
	SubcatRobloxCreated byte = 16
)

// Used with the Gears field of a Query.
const (
	GearMeleeWeapon        byte = 1
	GearRangedWeapon       byte = 2
	GearExplosive          byte = 3
	GearPowerUp            byte = 4
	GearNavigationEnhancer byte = 5
	GearMusicalInstrument  byte = 6
	GearSocialItem         byte = 7
	GearBuildingTool       byte = 8
	GearPersonalTransport  byte = 9
)

// Used with the Genres field of a Query.
const (
	GenreTownandCity byte = 1
	GenreMedieval    byte = 2
	GenreSciFi       byte = 3
	GenreFighting    byte = 4
	GenreHorror      byte = 5
	GenreNaval       byte = 6
	GenreAdventure   byte = 7
	GenreSports      byte = 8
	GenreComedy      byte = 9
	GenreWestern     byte = 10
	GenreMilitary    byte = 11
	GenreBuilding    byte = 13
	GenreFPS         byte = 14
	GenreRPG         byte = 15
	//GenreSkatePark byte = 12
)

// Used with the CurrencyType field of a Query.
const (
	CurrencyAll           byte = 0
	CurrencyRobux         byte = 1
	CurrencyTickets       byte = 2
	CurrencyCustomRobux   byte = 3
	CurrencyCustomTickets byte = 4
	CurrencyFree          byte = 5
)

// Used wuth SortType field of a Query.
const (
	SortRelevance       byte = 0
	SortMostFavorited   byte = 1
	SortBestselling     byte = 2
	SortRecentlyUpdated byte = 3
	SortPriceLowToHigh  byte = 4
	SortPriceHighToLow  byte = 5
)

// Used with the AggregationFrequency field of a Query.
const (
	AggrPastDay   byte = 0
	AggrPastWeek  byte = 1
	AggrPastMonth byte = 2
	AggrAllTime   byte = 3
)

// Used with the SortCurrency field of a Query.
const (
	SortCurrencyRobux   byte = 0
	SortCurrencyTickets byte = 1
)

// This is the maximum number of results that can be returned by the catalog
// search.
const MaxResults = 42

// Query is used with Search and SearchAll to query assets.
type Query struct {
	Gears                []byte
	Genres               []byte
	Subcategory          byte
	Category             byte
	CurrencyType         byte
	SortType             byte
	AggregationFrequency byte
	SortCurrency         byte
	Keyword              string
	CreatorID            int
	PxMin                int
	PxMax                int
	IncludeNotForSale    bool
	PageNumber           int
	ResultsPerPage       int
}

// Result represents a single asset returned from a call to Search or
// SearchAll.
type Result struct {
	AssetId                int64
	Name                   string
	Url                    string
	PriceInRobux           string
	PriceInTickets         string
	Updated                string
	Favorited              string
	Sales                  string
	Remaining              string
	Creator                string
	CreatorUrl             string
	PrivateSales           string
	PriceView              int32
	BestPrice              string
	ContentRatingTypeID    int32
	AssetTypeID            int32
	CreatorID              int32
	CreatedDate            string
	UpdatedDate            string
	IsForSale              bool
	IsPublicDomain         bool
	IsLimited              bool
	IsLimitedUnique        bool
	MinimumMembershipLevel int32
}

// Converts a Query to URL values.
func convertQuery(query Query) (values url.Values) {
	values = url.Values{}

	if v, ok := query["Gears"].([]byte); ok {
		for _, b := range v {
			values.Add(k, strconv.Itoa(int(b)))
		}
	}
	if v, ok := query["Genres"].([]byte); ok {
		for _, b := range v {
			values.Add(k, strconv.Itoa(int(b)))
		}
	}
	if v, ok := query["Subcategory"].(byte); ok {
		values.Set(k, strconv.Itoa(int(v)))
	}
	if v, ok := query["Category"].(byte); ok {
		values.Set(k, strconv.Itoa(int(v)))
	}
	if v, ok := query["CurrencyType"].(byte); ok {
		values.Set(k, strconv.Itoa(int(v)))
	}
	if v, ok := query["SortType"].(byte); ok {
		values.Set(k, strconv.Itoa(int(v)))
	}
	if v, ok := query["AggregationFrequency"].(byte); ok {
		values.Set(k, strconv.Itoa(int(v)))
	}
	if v, ok := query["SortCurrency"].(byte); ok {
		values.Set(k, strconv.Itoa(int(v)))
	}
	if v, ok := query["Keyword"].(string); ok {
		values.Set(k, v)
	}
	if v, ok := query["CreatorID"].(int); ok {
		values.Set(k, strconv.Itoa(v))
	}
	if v, ok := query["PxMin"].(int); ok {
		values.Set(k, strconv.Itoa(v))
	}
	if v, ok := query["PxMax"].(int); ok {
		values.Set(k, strconv.Itoa(v))
	}
	if v, ok := query["IncludeNotForSale"].(bool); ok {
		values.Set(k, strconv.FormatBool(v))
	}
	if v, ok := query["PageNumber"].(int); ok {
		values.Set(k, strconv.Itoa(v))
	}
	if v, ok := query["ResultsPerPage"].(int); ok {
		values.Set(k, strconv.Itoa(v))
	}
	return
}

// Search is used to perform a search query for Roblox assets.
func Search(client *rbxweb.Client, query Query) (result []Result, err error) {
	values := convertQuery(query)
	resp, err := client.Get(client.GetURL(`www`, `/catalog/json`, values))
	if err = client.AssertResp(resp, err); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&result)
	return
}

// SearchAll is similar to Search, but issues multiple requests until n
// results are found. If n is less than 0, then every found result will be
// returned. If PageNumber is specified in query, then requests will start
// from that page.
func SearchAll(client *rbxweb.Client, n int, query Query) (result []Result, err error) {
	if n == 0 {
		return
	}

	i := 0
	p := query.PageNumber
loop:
	for {
		query.Set("PageNumber", p)
		rs, err := Catalog(client, query)
		if err != nil {
			return nil, err
		}
		if len(rs) == 0 {
			break loop
		}
		for _, r := range rs {
			result = append(result, r...)
			i = i + 1
			if n > 0 && i >= n {
				break loop
			}
		}
		p = p + 1
	}
	return
}
