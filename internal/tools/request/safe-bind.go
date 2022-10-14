package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v5"
)

func SafeBind(c echo.Context, i interface{}) error {
	bodyBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := json.Unmarshal(bodyBytes, i); err != nil {
		return err
	}

	return nil
}
