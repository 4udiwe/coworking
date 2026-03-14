package entity

// HeatmapCell представляет одну ячейку тепловой карты: день недели, час и количество бронирований.
type HeatmapCell struct {
	Weekday uint8 // День недели (1-7, понедельник=0)
	Hour    uint8 // Час суток (0-23)
	Count   uint64 // Количество бронирований в данную ячейку
}
