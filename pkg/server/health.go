package server

func (s *Server) testConnectionToQtumd() error {
	_, err := s.qtumRPCClient.GetNetworkInfo()
	return err
}
