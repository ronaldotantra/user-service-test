package handler

import "github.com/SawitProRecruitment/UserService/service"

type userData struct {
	Id    int64  `json:"id"`
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type userDataLogin struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type baseResponse struct {
	Message string `json:"message"`
}

type responseWithData struct {
	baseResponse
	Data any `json:"data"`
}

func newBaseResponse(message string) *baseResponse {
	return &baseResponse{
		Message: message,
	}
}

func newSuccessGetByIDResponse(u *service.User) *responseWithData {
	return &responseWithData{
		baseResponse: baseResponse{
			Message: "Successfully get user data!",
		},
		Data: userData{
			Id:    u.Id,
			Name:  u.Name,
			Phone: u.Phone,
		},
	}
}

func newSuccessRegisterResponse(id *int64) *responseWithData {
	return &responseWithData{
		baseResponse: baseResponse{
			Message: "Successfully Register!",
		},
		Data: userData{
			Id: *id,
		},
	}
}

func newSuccessLogin(u *service.ResponseLogin) *responseWithData {
	return &responseWithData{
		baseResponse: baseResponse{
			Message: "Successfully login!",
		},
		Data: userDataLogin{
			Id:    u.UserId,
			Token: u.Token,
		},
	}
}
