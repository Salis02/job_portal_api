package healthcheck

type Service struct{}

func NewService() *Service{
	return &Service{}
}

func (s *Service) Status() map[string]string {
	return map[string]string{
		"status": "Ok",
	}
}