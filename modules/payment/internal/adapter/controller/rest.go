package controller

import (
	"net/http"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/controller/dto"
	"github.com/labstack/echo/v4"
)

type paymentController struct {
	paymentService payment.IPaymentService
}

func NewPaymentController(e *echo.Group, paymentService payment.IPaymentService) *paymentController {
	controller := &paymentController{paymentService: paymentService}
	e.POST("/payments", controller.CreatePayment)
	e.GET("/payments/:id", controller.GetPayment)
	e.GET("/payments", controller.GetPayments)
	e.PUT("/payments/:id", controller.UpdatePayment)
	e.DELETE("/payments/:id", controller.DeletePayment)
	return controller
}

func (c *paymentController) CreatePayment(ctx echo.Context) error {
	var paymentRequest *dto.CreatePaymentRequest
	if err := ctx.Bind(&paymentRequest); err != nil {
		return err
	}
	paymentData := paymentRequest.ToPayment()
	err := c.paymentService.CreatePayment(&paymentData)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentToResponse(paymentData))
}

func (c *paymentController) GetPayment(ctx echo.Context) error {
	id := ctx.Param("id")
	payment, err := c.paymentService.GetPayment(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentToResponse(payment))
}

func (c *paymentController) GetPayments(ctx echo.Context) error {
	payments, nextCursor, err := c.paymentService.GetPayments()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.Response().Header().Set("X-Next-Cursor", nextCursor)
	return ctx.JSON(http.StatusOK, dto.FromPaymentListToResponse(payments))
}

func (c *paymentController) UpdatePayment(ctx echo.Context) error {
	id := ctx.Param("id")
	var paymentRequest *dto.UpdatePaymentRequest
	if err := ctx.Bind(&paymentRequest); err != nil {
		return err
	}
	paymentData := paymentRequest.ToPayment(id)
	err := c.paymentService.UpdatePayment(&paymentData)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentToResponse(paymentData))
}

func (c *paymentController) DeletePayment(ctx echo.Context) error {
	id := ctx.Param("id")
	err := c.paymentService.DeletePayment(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}
