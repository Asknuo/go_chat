package api

type ApiGroup struct {
	LoginApi         LoginApi
	RegisterApi      RegisterApi
	PersonalInfoApi  PersonalInfoApi
	ForgetPsApi      ForgetPsApi
	LogoutApi        LogoutApi
	GethistorymsgApi GethistorymsgApi
}
