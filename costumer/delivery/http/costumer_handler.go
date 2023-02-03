package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/bxcodec/go-clean-arch/domain"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// CostumerHandler  represent the httphandler for costumer
type CostumerHandler struct {
	AUsecase domain.CostumerUsecase
}

// NewCostumerHandler will initialize the costumer/ resources endpoint
func NewCostumerHandler(e *echo.Echo, us domain.CostumerUsecase) {
	handler := &CostumerHandler{
		AUsecase: us,
	}
	e.GET("/costumer", handler.FetchCostumer)
	e.POST("/costumer", handler.Store)
	e.GET("/costumer/:id", handler.GetByID)
	e.DELETE("/costumer/:id", handler.Delete)
}

// FetchCostumer will fetch the costumer based on given params
func (a *CostumerHandler) FetchCostumer(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

// GetByID will get costumer by given id
func (a *CostumerHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.AUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.Costumer) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the costumer by given request body
func (a *CostumerHandler) Store(c echo.Context) (err error) {
	var costumer domain.Costumer
	err = c.Bind(&costumer)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&costumer); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &costumer)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, costumer)
}

// Delete will delete costumer by given param
func (a *CostumerHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = a.AUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
