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

type ProductDTO struct {
	ID          int32
	Name        string
	Description string
	Category    string
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

func (p *ProductService) ListProdByName(ctx context.Context, name string) ([]ProductDTO, error) {

	searchPattern := "%" + name + "%"

	prods, err := p.product.ListProductByName(ctx, db.ListProductByNameParams{Name: searchPattern, Limit: 10, Offset: 0})
	if err != nil {
		return nil, err
	}

	// Выделяем память под слайс DTO
	result := make([]ProductDTO, len(prods))
	for i, prod := range prods {
		// Безопасно извлекаем текст из pgtype.Text
		var desc string
		if prod.Description.Valid {
			desc = prod.Description.String
		}

		result[i] = ProductDTO{
			ID:          prod.ID,
			Name:        prod.Name,
			Category:    prod.Category,
			Description: desc,
		}

	}

	return result, nil
}
