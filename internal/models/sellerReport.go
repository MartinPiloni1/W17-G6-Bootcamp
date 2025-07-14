package models

// SellerReport representa el reporte de cantidad de sellers por localidad
type SellerReport struct {
	LocalityID   string `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	SellersCount int    `json:"sellers_count"`
}
