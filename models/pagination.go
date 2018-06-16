package models

import (
	"math"
)

// Pagination helper strcut to create paginated response
type Pagination struct {
	From        uint64
	To          uint64
	PerPage     uint64
	Total       uint64
	CurrentPage uint64
	LastPage    uint64
	Data        interface{}
}

// Make creates a new pagination model using given data
func (pgn *Pagination) Make(data interface{}, total, currentPage, perPage uint64) {
	pgn.Data = data
	pgn.From = (currentPage - 1) * perPage
	pgn.To = pgn.From + perPage
	pgn.CurrentPage = currentPage
	pgn.PerPage = perPage
	pgn.Total = total
	pgn.LastPage = func() uint64 {
		p := math.Floor(float64(total / perPage))
		if math.Mod(float64(total), float64(perPage)) == 0.0 {
			return uint64(p)
		}

		return uint64(p) + 1
	}()

	return

}
