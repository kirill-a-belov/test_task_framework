package reservation

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	orderModel "github.com/kirill-a-belov/test_task_framework/internal/pkg/order/model"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/reservation/model"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

type createReservationForm struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
	UserID  int       `json:"user_id"`
	// Deprecated (KB)
	UserEmail  string `json:"email"`
	CreditCard string `json:"credit_card"`
}

// @Title createReservationHandler
// @Description reservation creation
// @Param   form    body     reservation.createReservationForm     true        "Reservation details"
// @Success 200 {object} model.Reservation "reservation"
// @Failure XXX {object} object    "error message"
// @Tags /api/reservation/
// @Router /api/reservation/create [post]
func (m *Module) createReservationHandler(c *gin.Context) {
	ctx := c.Request.Context()
	_, span := tracer.Start(ctx, "internal.pkg.reservation.Module.createReservationHandler")
	defer span.End()

	var form createReservationForm
	if err := c.BindJSON(&form); err != nil {
		_ = c.Error(errors.Wrap(err, "parsing form"))
		c.Status(http.StatusBadRequest)

		return
	}

	reservation, err := m.CreateReservation(ctx, &CreateReservationRequest{
		UserID: form.UserID,
		ReservationDetails: model.ReservationDetails{
			OrderDetails: orderModel.OrderDetails{
				HotelID: form.HotelID,
				RoomID:  form.RoomID,
				From:    form.From,
				To:      form.To,
			},
		},
	})
	if err != nil {
		_ = c.Error(errors.Wrap(err, "creating reservation"))
		c.Status(http.StatusBadRequest)

		return
	}

	logger.New("createReservationHandler").Info("Order successfully created:",
		form, reservation,
	)

	c.JSON(http.StatusOK, reservation.Reservation)
}

func (m *Module) RegisterRoutes(e *gin.Engine) {
	e.POST("/api/reservation/create", m.createReservationHandler)
}
