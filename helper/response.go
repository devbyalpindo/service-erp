package helper

import "erp-service/model/dto"

func ResponseError(status string, err interface{}, code int) dto.Response {
	return dto.Response{
		StatusCode: code,
		Status:     status,
		Error:      err,
		Data:       nil,
	}
}

func ResponseSuccess(status string, err interface{}, data map[string]interface{}, code int) dto.Response {
	return dto.Response{
		StatusCode: code,
		Status:     status,
		Error:      err,
		Data:       data,
	}
}
