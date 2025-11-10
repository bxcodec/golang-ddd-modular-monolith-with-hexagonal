package controller

import (
	"net/http"
	"strconv"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/controller/dto"
	"github.com/labstack/echo/v4"
)

type paymentController struct {
	paymentService payment.IPaymentService
}

func NewPaymentController(e *echo.Group, paymentService payment.IPaymentService) (controller *paymentController) {
	controller = &paymentController{paymentService: paymentService}
	e.POST("/payments", controller.CreatePayment)
	e.GET("/payments/:id", controller.GetPayment)
	e.GET("/payments", controller.FetchPayments)
	e.PUT("/payments/:id", controller.UpdatePayment)
	e.DELETE("/payments/:id", controller.DeletePayment)
	return controller
}

func (c *paymentController) CreatePayment(ctx echo.Context) (err error) {
	var paymentRequest *dto.CreatePaymentRequest
	if err = ctx.Bind(&paymentRequest); err != nil {
		return err
	}
	paymentData := paymentRequest.ToPayment()
	err = c.paymentService.CreatePayment(&paymentData)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentToResponse(paymentData))
}

func (c *paymentController) GetPayment(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	payment, err := c.paymentService.GetPayment(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentToResponse(payment))
}

func (c *paymentController) FetchPayments(ctx echo.Context) (err error) {
	cursor := ctx.QueryParam("cursor")
	limit := 10
	if limitParam := ctx.QueryParam("limit"); limitParam != "" {
		if val, convErr := strconv.Atoi(limitParam); convErr == nil && val > 0 {
			limit = val
		}
	}

	result, nextCursor, err := c.paymentService.FetchPayments(payment.FetchPaymentsParams{
		Cursor:   cursor,
		Limit:    limit,
		Currency: ctx.QueryParam("currency"),
		Status:   ctx.QueryParam("status"),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.Response().Header().Set("X-Next-Cursor", nextCursor)
	return ctx.JSON(http.StatusOK, dto.FromPaymentListToResponse(result))
}

func (c *paymentController) UpdatePayment(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	var paymentRequest *dto.UpdatePaymentRequest
	if err = ctx.Bind(&paymentRequest); err != nil {
		return err
	}
	paymentData := paymentRequest.ToPayment(id)
	err = c.paymentService.UpdatePayment(&paymentData)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentToResponse(paymentData))
}

func (c *paymentController) DeletePayment(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	err = c.paymentService.DeletePayment(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}
