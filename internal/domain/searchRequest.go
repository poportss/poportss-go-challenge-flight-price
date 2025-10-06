package domain

import "time"

type SearchRequest struct {
	Origin      string    `form:"origin" binding:"required,len=3"`
	Destination string    `form:"destination" binding:"required,len=3"`
	StartDate   time.Time `form:"starDate" time_format:"2006-01-02" binding:"required"`
	EndDate     time.Time `form:"endDate" time_format:"2006-01-02" binding:"required"`
}
