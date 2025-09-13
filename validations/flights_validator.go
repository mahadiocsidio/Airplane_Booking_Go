package validations
import(
	"time"
	"airplane_booking_go/models"
)

type CreateFlightRequest struct {
	Airline       string          `json:"airline" binding:"required"`
	FlightNumber  string          `json:"flightNumber" binding:"required"`
	Departure     models.Airport  `json:"departure" binding:"required"`
	Arrival       models.Airport  `json:"arrival" binding:"required"`
	DepartureTime time.Time       `json:"departureTime" binding:"required"`
	ArrivalTime   time.Time       `json:"arrivalTime" binding:"required"`
	Duration      int             `json:"duration" binding:"required"`
	Price         float64         `json:"price" binding:"required"`
	Seats         []models.Seat   `json:"seats" binding:"required"`
}

type SearchFlightRequest struct {
	From  string `form:"from" binding:"omitempty,len=3"`              // kode IATA 3 huruf
	To    string `form:"to" binding:"omitempty,len=3"`
	Date  string `form:"date" binding:"omitempty,datetime=2006-01-02"` // format: YYYY-MM-DD
	Page  int    `form:"page,default=1"`
	Limit int    `form:"limit,default=10"`
}
