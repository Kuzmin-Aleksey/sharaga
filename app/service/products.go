package service

import (
	"app/models"
	"strconv"
)

func (s *Service) NewProduct(product *models.Product) error {
	id := new(models.Id)

	if err := s.request(id, post, urlProducts, product); err != nil {
		return err
	}

	product.Id = id.Id

	return nil
}

func (s *Service) GetProducts() ([]models.Product, error) {
	product := make([]models.Product, 0)

	if err := s.request(&product, get, urlProducts, nil); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) UpdateProduct(product *models.Product) error {
	if err := s.request(nil, put, urlProducts, product); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProduct(productId int) error {
	if err := s.request(nil, del, urlProducts, nil,
		[2]string{"product_id", strconv.Itoa(productId)},
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) NewProductType(t *models.ProductType) error {
	id := new(models.Id)

	if err := s.request(id, post, urlProductsTypes, t); err != nil {
		return err
	}

	t.Id = id.Id
	return nil
}

func (s *Service) GetProductTypes() ([]models.ProductType, error) {
	productTypes := make([]models.ProductType, 0)
	if err := s.request(&productTypes, get, urlProductsTypes, nil); err != nil {
		return nil, err
	}

	return productTypes, nil
}

func (s *Service) UpdateProductType(productType *models.ProductType) error {
	if err := s.request(nil, put, urlProductsTypes, productType); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProductType(productTypeId int) error {
	if err := s.request(nil, del, urlProductsTypes, nil,
		[2]string{"product_type_id", strconv.Itoa(productTypeId)},
	); err != nil {
		return err
	}

	return nil
}
