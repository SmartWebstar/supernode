package restserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pastelnetwork/go-pastel"
)

func (s *RESTServer) getAllIDTickets(c echo.Context, idtype string) error {
	t, err := s.pslNode.ListIDTickets(idtype)
	if err != nil {
		if err.Error() == "Nothing found" {
			return c.JSON(http.StatusOK, t)
		}
		return err
	}
	return c.JSON(http.StatusOK, t)
}
func (s *RESTServer) getMyIDTickets(c echo.Context, idtype string) error {
	p, err := s.pslNode.ListPastelIDs()
	if err != nil {
		if err.Error() == "Nothing found" {
			return c.JSON(http.StatusOK, p)
		}
		return err
	}
	t := pastel.IDTickets{}
	for _, pid := range p {
		ticket, err := s.pslNode.FindIDTicket(pid.PastelID)
		if err != nil {
			if err.Error() == "Key is not found" {
				continue
			}
			return err
		}
		if ticket.Ticket.IDType == idtype {
			t = append(t, *ticket)
		}
	}
	return c.JSON(http.StatusOK, t)
}

func (s *RESTServer) getIDTicket(c echo.Context, idtype string) error {
	id := c.Param("id")
	t, err := s.pslNode.FindIDTicket(id)
	if err != nil {
		if err.Error() == "Key is not found" {
			return c.JSON(http.StatusOK, "")
		}
		if err.Error() == "Nothing found" {
			return c.JSON(http.StatusOK, "")
		}
		return err
	}
	if t.Ticket.IDType != idtype {
		return c.JSON(http.StatusOK, "")
	}
	return c.JSON(http.StatusOK, t)
}

func (s *RESTServer) GetAllIDTickets(c echo.Context) error {
	return s.getAllIDTickets(c, "personal")
}
func (s *RESTServer) GetMyIDTickets(c echo.Context) error {
	return s.getMyIDTickets(c, "personal")
}
func (s *RESTServer) GetIDTicket(c echo.Context) error {
	return s.getIDTicket(c, "personal")
}
func (s *RESTServer) GetAllMNIDTickets(c echo.Context) error {
	return s.getAllIDTickets(c, "mn")
}
func (s *RESTServer) GetMyMNIDTickets(c echo.Context) error {
	return s.getMyIDTickets(c, "mn")
}
func (s *RESTServer) GetMNIDTicket(c echo.Context) error {
	return s.getIDTicket(c, "mn")
}
func (s *RESTServer) GetPastelIDs(c echo.Context) error {
	p, err := s.pslNode.ListPastelIDs()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}
