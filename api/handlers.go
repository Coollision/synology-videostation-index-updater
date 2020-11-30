package api

import "net/http"

func (s *Server) bla(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("hallo het werkt yeay"))
	if err != nil{
		panic(err)
	}
}

func (s *Server) reIndex(w http.ResponseWriter, r *http.Request)  {
	err := s.videoAPI.ListLibraries()
	if err != nil{
		panic(err)
	}
	_, err = w.Write([]byte("started reindexing"))
	if err != nil{
		panic(err)
	}
}
