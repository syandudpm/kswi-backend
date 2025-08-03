package menu

type MenuResponse struct {
	ID       uint           `json:"id" gorm:"column:id"`
	ParentID uint           `json:"parent_id" gorm:"column:parent_id"`
	Sort     int            `json:"sort" gorm:"column:sort"`
	Name     string         `json:"name" gorm:"column:name"`
	Route    *string        `json:"route" gorm:"column:route"`
	Icon     *string        `json:"icon" gorm:"column:icon"`
	Children []MenuResponse `json:"children" gorm:"-"`
}

type MenuDetailResponse struct {
	ID       uint    `json:"id" gorm:"column:id"`
	ParentID uint    `json:"parent_id" gorm:"column:parent_id"`
	Sort     int     `json:"sort" gorm:"column:sort"`
	Name     string  `json:"name" gorm:"column:name"`
	Route    *string `json:"route" gorm:"column:route"`
	Icon     *string `json:"icon" gorm:"column:icon"`
	IsActive bool    `json:"is_active" gorm:"column:is_active"`
}
