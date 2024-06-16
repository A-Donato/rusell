package structures

type Item struct {
	Id                         string
	Name                       string
	Highest_actual_price       string
	Highest_actual_price_store string
	Lowest_actual_price        string
	Lowest_actual_price_store  string
	Average_price              string
	Search_terms               []string
	Tags                       []string
	Is_in_sale                 bool
}

type Scrap_tagets struct {
	Id                 string
	Name               string
	Base_url           string
	Logo_url           string
	Search_bar_element string
}

type Price_analysis struct {
	Id                        string
	Item_id                   string
	Analysis_interval_in_days int
	Last_analysis             string
	Measurements              map[string][]int
}

type Items_in_target struct {
	Id         string
	Item_id    string
	Is_enabled string
	SKU        string
	GTIN       string
	UPC        string
	Targets    map[string]string
}
