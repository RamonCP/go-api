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

func (pr *ProductRepository) CreateProduct(product domain.Product) (domain.Product, error) {
	query, err := pr.connection.Prepare("INSERT INTO product (product_name, price) VALUES($1, $2) RETURNING id")

	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}

	defer query.Close()

	var id int
	err = query.QueryRow(product.Name, product.Price).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}

	return domain.Product{
		ID: id,
		Name: product.Name,
		Price: product.Price,
	}, nil
}

func (pr *ProductRepository) DeleteProduct(id int) error {
	query, err := pr.connection.Prepare("DELETE FROM product WHERE id = $1")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer query.Close()

	result, err := query.Exec(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if rowsAffected == 0 {
		fmt.Println(err)
		return fmt.Errorf("Produto com id não encontrado")
	}

	return nil
}

func (pr *ProductRepository) UpdateProduct(product domain.Product, id int) (domain.Product, error) {
	query, err := pr.connection.Prepare("UPDATE product SET product_name=$1, price=$2 WHERE id = $3")
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}

	result, err := query.Exec(product.Name, product.Price, id)
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}

	if rowsAffected == 0 {
		fmt.Println(err)
		return domain.Product{}, err
	}

	return domain.Product{
		ID: id,
		Name: product.Name,
		Price: product.Price,
	}, nil
}