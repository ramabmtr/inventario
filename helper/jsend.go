// follow this convention https://github.com/omniti-labs/jsend
// to standardize json output

package helper

import "github.com/ramabmtr/inventario/config"

type Jsend struct {
	Status  string      `json:"status" binding:"required"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func buildResponse(status string, msg string, data interface{}) Jsend {
	return Jsend{Status: status, Message: msg, Data: data}
}

func FailResponse(msg string) Jsend {
	return buildResponse("failed", msg, nil)
}

func ErrorResponse(msg string) Jsend {
	if !config.Env.App.Debug {
		return buildResponse("error", config.ErrDefault.Error(), nil)
	}

	return buildResponse("error", msg, nil)
}

func SuccessResponse() Jsend {
	return buildResponse("success", "", nil)
}

func ObjectResponse(data interface{}, envelope string) Jsend {
	if envelope != "" {
		data = map[string]interface{}{
			envelope: data,
		}
	}
	return buildResponse("success", "", data)
}
