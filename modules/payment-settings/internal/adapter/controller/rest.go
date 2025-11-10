package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/adapter/controller/dto"
)

type paymentSettingController struct {
	paymentSettingsService paymentsettings.IPaymentSettingsService
}

func NewPaymentSettingController(e *echo.Group, paymentSettingsService paymentsettings.IPaymentSettingsService) (controller *paymentSettingController) {
	controller = &paymentSettingController{paymentSettingsService: paymentSettingsService}
	e.GET("/payment-settings", controller.FetchPaymentSettings)
	e.POST("/payment-settings", controller.CreatePaymentSetting)
	e.GET("/payment-settings/:id", controller.GetPaymentSetting)
	e.PUT("/payment-settings/:id", controller.UpdatePaymentSetting)
	e.DELETE("/payment-settings/:id", controller.DeletePaymentSetting)
	return controller
}

func (c *paymentSettingController) FetchPaymentSettings(ctx echo.Context) (err error) {
	cursor := ctx.QueryParam("cursor")
	limit := 10
	if limitParam := ctx.QueryParam("limit"); limitParam != "" {
		if val, convErr := strconv.Atoi(limitParam); convErr == nil && val > 0 {
			limit = val
		}
	}

	result, nextCursor, err := c.paymentSettingsService.FetchPaymentSettings(paymentsettings.PaymentSettingFetchParams{
		Cursor:     cursor,
		Limit:      limit,
		Currency:   ctx.QueryParam("currency"),
		SettingKey: ctx.QueryParam("settingKey"),
		Status:     ctx.QueryParam("status"),
	})
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("X-Next-Cursor", nextCursor)
	return ctx.JSON(http.StatusOK, dto.FromPaymentSettingListToResponse(result))
}

func (c *paymentSettingController) CreatePaymentSetting(ctx echo.Context) (err error) {
	var paymentSettingRequest *dto.CreatePaymentSettingRequest
	if err = ctx.Bind(&paymentSettingRequest); err != nil {
		return err
	}
	paymentSetting := paymentSettingRequest.ToPaymentSetting()
	err = c.paymentSettingsService.CreatePaymentSetting(&paymentSetting)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, dto.FromPaymentSettingToResponse(paymentSetting))
}

func (c *paymentSettingController) GetPaymentSetting(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	paymentSetting, err := c.paymentSettingsService.GetPaymentSetting(id)
	if err != nil {
		return err
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
		return err
	}
	return ctx.JSON(http.StatusOK, dto.FromPaymentSettingToResponse(paymentSetting))
}

func (c *paymentSettingController) DeletePaymentSetting(ctx echo.Context) (err error) {
	id := ctx.Param("id")
	err = c.paymentSettingsService.DeletePaymentSetting(id)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
