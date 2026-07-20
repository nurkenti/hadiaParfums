package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/nurkenti/hadiaParfums/internal/repository/sqlc"
)

type ProductService struct {
	product *db.Queries
}

func NewProductService(product *db.Queries) *ProductService {
	return &ProductService{product: product}
}

func (p *ProductService) AddProduct(name string, category string, descrip pgtype.Text, ctx context.Context) error {
	_, err := p.product.CreateProduct(ctx, db.CreateProductParams{
		Name:        name,
		Category:    category,
		Description: descrip,
	})
	if err != nil {
		return errors.New("Product error")
	}

	return nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, id int32) error {
	err := p.product.DeleteProductByID(ctx, id)
	if err != nil {
		return errors.New("Product errors")
	}
	return nil
}

func (p *ProductService) ListProdByName()
