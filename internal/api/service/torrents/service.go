package torrents

var _ Service = (*service)(nil)

type Service interface {
	i()

	Start()

	Shutdown()
}

type service struct {
}

func (s *service) i() {}
