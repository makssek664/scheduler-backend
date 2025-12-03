package main;
import (
	"time"
	"gorm.io/gorm"
	"github.com/guregu/null"
)
type Event struct {
	gorm.Model
	UserID uint 
	Date time.Time 
	Name string
	Desc null.String 
	Color null.Int
}
