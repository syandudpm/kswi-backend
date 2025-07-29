package menu

type MenuResponse struct {
	ID       uint           `json:"id"`
	Sort     int            `json:"sequence"`
	Name     string         `json:"name"`
	Route    *string        `json:"route,omitempty"`
	Icon     *string        `json:"icon"`
	Children []MenuResponse `json:"children" gorm:"-"`
}
