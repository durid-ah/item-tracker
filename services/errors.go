package services

type ServiceError struct {
	Message string
}

func (err *ServiceError) Error() string {
	return err.Message
}
