package service

import (
	"app/models"
	"strconv"
)

func (s *Service) NewPartner(partner *models.Partner) error {
	id := new(models.Id)

	if err := s.request(id, post, urlPartners, partner); err != nil {
		return err
	}

	partner.Id = id.Id

	return nil
}

func (s *Service) GetPartners() ([]models.Partner, error) {
	partners := make([]models.Partner, 0)

	if err := s.request(&partners, get, urlPartners, nil); err != nil {
		return nil, err
	}

	return partners, nil
}

func (s *Service) UpdatePartner(partner *models.Partner) error {
	if err := s.request(nil, put, urlPartners, partner); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeletePartner(partnerId int) error {
	if err := s.request(nil, del, urlPartners, nil,
		[2]string{"partner_id", strconv.Itoa(partnerId)},
	); err != nil {
		return err
	}

	return nil
}
