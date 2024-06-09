package controllers

import (
	"testing"
)

func TestAuthController_SignInReq(t *testing.T) {
	//testCases := []struct {
	//	name        string
	//	expectedErr error
	//	request     models.SignInReq
	//}{
	//	{
	//		name: "valid request",
	//		request: models.SignInReq{
	//			Email:    "test@example.com",
	//			Password: "ValidPassword123",
	//		},
	//		expectedErr: nil,
	//	},
	//	{
	//		name: "invalid email",
	//		request: models.SignInReq{
	//			Email:    "invalid_email",
	//			Password: "ValidPassword123",
	//		},
	//		expectedErr: models.ErrInvalidEmail,
	//	},
	//	{
	//		name: "invalid password",
	//		request: models.SignInReq{
	//			Email:    "test@example.com",
	//			Password: "invalid",
	//		},
	//		expectedErr: models.ErrInvalidPassword,
	//	},
	//	{
	//		name: "invalid request",
	//		request: models.SignInReq{
	//			Email:    "invalid_email",
	//			Password: "invalid",
	//		},
	//		expectedErr: models.ErrInvalidEmail,
	//	},
	//}
	//
	//for _, testCase := range testCases {
	//	t.Run(testCase.name, func(t *testing.T) {
	//		err := testCase.request.Validate()
	//		assert.Equal(t, testCase.expectedErr, err, testCase.name)
	//	})
	//}
}

func TestAuthController_SignUpReq(t *testing.T) {
	//testCases := []struct {
	//	name        string
	//	expectedErr error
	//	request     models.SignUpReq
	//}{
	//	{
	//		name: "valid request",
	//		request: models.SignUpReq{
	//			Email:       "test@example.com",
	//			Password:    "ValidPassword123",
	//			PhoneNumber: "87082260629",
	//		},
	//		expectedErr: nil,
	//	},
	//	{
	//		name: "invalid email",
	//		request: models.SignUpReq{
	//			Email:       "invalid_email",
	//			Password:    "ValidPassword123",
	//			PhoneNumber: "87082260629",
	//		},
	//		expectedErr: models.ErrInvalidEmail,
	//	},
	//	{
	//		name: "invalid password",
	//		request: models.SignUpReq{
	//			Email:       "test@example.com",
	//			Password:    "invalid",
	//			PhoneNumber: "87082260629",
	//		},
	//		expectedErr: models.ErrInvalidPassword,
	//	},
	//	{
	//		name: "invalid request",
	//		request: models.SignUpReq{
	//			Email:       "invalid_email",
	//			Password:    "invalid",
	//			PhoneNumber: "invalid",
	//		},
	//		expectedErr: models.ErrInvalidEmail,
	//	},
	//	{
	//		name: "invalid phone number format",
	//		request: models.SignUpReq{
	//			Email:       "test@example.com",
	//			Password:    "ValidPassword123",
	//			PhoneNumber: "invalid_number",
	//		},
	//		expectedErr: models.ErrInvalidPhoneNumber,
	//	},
	//	{
	//		name: "valid phone number format",
	//		request: models.SignUpReq{
	//			Email:       "test@example.com",
	//			Password:    "ValidPassword123",
	//			PhoneNumber: "87082260629",
	//		},
	//		expectedErr: nil,
	//	},
	//}
	//
	//for _, testCase := range testCases {
	//	t.Run(testCase.name, func(t *testing.T) {
	//		err := testCase.request.Validate()
	//		assert.Equal(t, testCase.expectedErr, err, testCase.name)
	//	})
	//}
}

func TestOAuthController_StateValidation(t *testing.T) {
	//testCases := []struct {
	//	name  string
	//	valid bool
	//	state string
	//}{
	//	{
	//		name:  "invalid state",
	//		valid: false,
	//		state: "invalid_state",
	//	},
	//	{
	//		name:  "nil state",
	//		valid: false,
	//		state: "",
	//	},
	//	{
	//		name:  "invalid state",
	//		valid: false,
	//		state: "invalid_state",
	//	},
	//	{
	//		name:  "valid state: new uuid",
	//		valid: true,
	//		state: uuid.New().String(),
	//	},
	//	{
	//		name:  "valid state: copied uuid",
	//		valid: true,
	//		state: "ac575f74-bed1-4181-9ed2-1df726774044",
	//	},
	//}
	//
	//for _, testCase := range testCases {
	//	t.Run(testCase.name, func(t *testing.T) {
	//		assert.Equal(t, testCase.valid, tools.IsUUID(testCase.state), testCase.name)
	//	})
	//}
}
