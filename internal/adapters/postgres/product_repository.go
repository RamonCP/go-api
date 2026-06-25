package postgres

import (
	"database/sql"
	"fmt"
	"go-api/internal/core/domain"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) *ProductRepository {
	return &ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]domain.Product, error) {
	query := "SELECT id, product_name, price FROM product"
	rows, err := pr.connection.Query(query)

	if err != nil {
		fmt.Println(err)
		return []domain.Product{}, err
	}

	var productList []domain.Product
	var productObj domain.Product

	for rows.Next() {
		err = rows.Scan(
			&productObj.ID,
			&productObj.Name,
			&productObj.Price,
		)

		if err != nil {
			fmt.Println(err)
			return []domain.Product{}, err
		}

		productList = append(productList, productObj)
	}

	rows.Close()

	return productList, nil
}

func (pr *ProductRepository) GetProductById(id int) (domain.Product, error) {
	query, err := pr.connection.Prepare("SELECT * FROM product WHERE id = $1")
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}

	defer query.Close()

	var productObj domain.Product
	err = query.QueryRow(id).Scan(&productObj.ID, &productObj.Name, &productObj.Price)

	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}

	return productObj, nil
}