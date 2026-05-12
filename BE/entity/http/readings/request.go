package readings

import "backend/util/pagination"

type GetReadings struct {
	Pgn pagination.Pagination `form:"pgn"`
}
