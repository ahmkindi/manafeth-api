package models

type Product struct {
	ProductID     int64  `json:"product_id"`
	ProductDescEN string `json:"product_desc_en"`
	ProductDescAR string `json:"product_desc_ar"`
}

type Country struct {
	CountryID     int64  `json:"country_id"`
	CountryNameEN string `json:"country_name_en"`
	CountryNameAR string `json:"country_name_ar"`
}

type Port struct {
	PortID     int64  `json:"port_id"`
	PortNameEN string `json:"port_name_en"`
	PortNameAR string `json:"port_name_ar"`
	PortTypeEN string `json:"port_type_en"`
	PortTypeAR string `json:"port_type_ar"`
	ModeID     int    `json:"mode_id"`
}

type TradeSummary struct {
	Year              int   `json:"year"`
	ImportValue       int64 `json:"import_value"`
	ExportValue       int64 `json:"export_value"`
	ReExportValue     int64 `json:"reexport_value"`
	TradeBalanceValue int64 `json:"tradebalance_value"`
	TotalTradeValue   int64 `json:"total_trade_value"`
}

type TradeBalance struct {
	StartYear     int   `json:"start_year"`
	EndYear       int   `json:"end_year"`
	TotalImport   int64 `json:"total_import"`
	TotalExport   int64 `json:"total_export"`
	TotalReExport int64 `json:"total_reexport"`
	TradeBalance  int64 `json:"trade_balance"`
}

type AggregateRequest struct {
	DateRange  DateRange  `json:"date_range"`
	TradeTypes []string   `json:"trade_types,omitempty"`
	GroupBy    []string   `json:"group_by"`
	Filters    Filters    `json:"filters,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
	Sorting    Sorting    `json:"sorting,omitempty"`
}

type DateRange struct {
	StartYear int `json:"start_year"`
	EndYear   int `json:"end_year"`
}

type Filters struct {
	ProductIDs []int64  `json:"product_ids,omitempty"`
	CountryIDs []int64  `json:"country_ids,omitempty"`
	PortIDs    []int64  `json:"port_ids,omitempty"`
	PortTypes  []string `json:"port_types,omitempty"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Sorting struct {
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type AggregateResult struct {
	Year          *int    `json:"year,omitempty"`
	ProductID     *int64  `json:"product_id,omitempty"`
	ProductDescEN *string `json:"product_desc_en,omitempty"`
	ProductDescAR *string `json:"product_desc_ar,omitempty"`
	CountryID     *int64  `json:"country_id,omitempty"`
	CountryNameEN *string `json:"country_name_en,omitempty"`
	CountryNameAR *string `json:"country_name_ar,omitempty"`
	PortID        *int64  `json:"port_id,omitempty"`
	PortNameEN    *string `json:"port_name_en,omitempty"`
	PortNameAR    *string `json:"port_name_ar,omitempty"`
	TradeType     *string `json:"trade_type,omitempty"`
	TotalValue    int64   `json:"total_value"`
}

type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalCount  int64 `json:"total_count"`
	TotalPages  int   `json:"total_pages"`
}
