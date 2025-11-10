package controller

import (
	"net/http"
	"strconv"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/adapter/controller/dto"
	"github.com/labstack/echo/v4"
)

type paymentSettingController struct {
	paymentSettingsService paymentsettings.IPaymentSettingsService
}

func NewPaymentSettingController(e *echo.Group, paymentSettingsService paymentsettings.IPaymentSettingsService) (controller *paymentSettingController) {
	controller = &paymentSettingController{paymentSettingsService: paymentSettingsService}
	e.GET("/payment-settings", controller.FetchPaymentSettings)
	e.POST("/payment-settings", controller.CreatePaymentSetting)
	e.PUT("/payment-settings/:id", controller.UpdatePaymentSetting)
	e.DELETE("/payment-settings/:id", controller.DeletePaymentSetting)
	return controller
}

func (c *paymentSettingController) FetchPaymentSettings(ctx echo.Context) (err error) {
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	params := paymentsettings.PaymentSettingFetchParams{
		Currency: ctx.QueryParam("currency"),
		Limit:    limit,
		Cursor:   ctx.QueryParam("cursor"),
	}
	settings, nextCursor, err := c.paymentSettingsService.FetchPaymentSettings(params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.Response().Header().Set("X-Next-Cursor", nextCursor)
	return ctx.JSON(http.StatusOK, dto.FromPaymentSettingListToResponse(settings))
}

func (c *paymentSettingController) CreatePaymentSetting(ctx echo.Context) (err error) {
	var paymentSettingRequest *dto.CreatePaymentSettingRequest
	if err = ctx.Bind(&paymentSettingRequest); err != nil {
		return err
	}
	paymentSetting := paymentSettingRequest.ToPaymentSetting()
	err = c.paymentSettingsService.CreatePaymentSetting(&paymentSetting)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentSettingToResponse(paymentSetting))
}

func (c *paymentSettingController) UpdatePaymentSetting(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	var paymentSettingRequest *dto.UpdatePaymentSettingRequest
	if err = ctx.Bind(&paymentSettingRequest); err != nil {
		return err
	}
	paymentSetting := paymentSettingRequest.ToPaymentSetting(id)
	err = c.paymentSettingsService.UpdatePaymentSetting(&paymentSetting)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentSettingToResponse(paymentSetting))
}

func (c *paymentSettingController) DeletePaymentSetting(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	err = c.paymentSettingsService.DeletePaymentSetting(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusNoContent)
}
