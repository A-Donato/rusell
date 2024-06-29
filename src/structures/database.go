package structures

type Item struct {
	Id                         string   `firestore:"id,omitempty"`
	Name                       string   `firestore:"name,omitempty"`
	Highest_actual_price       string   `firestore:"highest_actual_price,omitempty"`
	Highest_actual_price_store string   `firestore:"highest_actual_price_store,omitempty"`
	Lowest_actual_price        string   `firestore:"lowest_actual_price,omitempty"`
	Lowest_actual_price_store  string   `firestore:"lowest_actual_price_store,omitempty"`
	Average_price              string   `firestore:"average_price,omitempty"`
	Search_terms               []string `firestore:"search_terms,omitempty"`
	Tags                       []string `firestore:"tags,omitempty"`
	Is_in_sale                 bool     `firestore:"is_in_sale,omitempty"`
	Images_urls                []string `firestore:"images_urls,omitempty"`
}

type Scrap_targets struct {
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

type Target struct {
	Url        string `firestore:"url,omitempty"`
	HtmlTarget string `firestore:"html_target,omitempty"`
}

type Items_in_target struct {
	Id         string            `firestore:"id,omitempty"`
	Item_id    string            `firestore:"item_id,omitempty"`
	Is_enabled bool              `firestore:"is_enabled,omitempty"`
	SKU        string            `firestore:"sku,omitempty"`
	GTIN       string            `firestore:"gtin,omitempty"`
	UPC        string            `firestore:"upc,omitempty"`
	Targets    map[string]Target `firestore:"targets,omitempty"`
}
