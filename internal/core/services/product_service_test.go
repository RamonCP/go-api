// Teste black-box (package services_test): exercita o pacote pela sua API
// pública, exatamente como um consumidor real faria.
package services_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go-api/internal/core/domain"
	"go-api/internal/core/services"
)

// fakeProductRepository é um test double escrito à mão. Cada método é um
// campo-função, então cada caso de teste configura só o comportamento de
// que precisa — sem framework de mock.
type fakeProductRepository struct {
	getByID func(id int) (domain.Product, error)
}

func (f *fakeProductRepository) GetProductById(id int) (domain.Product, error) {
	return f.getByID(id)
}

// Os demais métodos existem apenas para satisfazer a interface
// ports.ProductRepository; não são exercitados neste teste.
func (f *fakeProductRepository) GetProducts() ([]domain.Product, error) { panic("não usado") }
func (f *fakeProductRepository) CreateProduct(domain.Product) (domain.Product, error) {
	panic("não usado")
}
func (f *fakeProductRepository) DeleteProduct(int) error { panic("não usado") }
func (f *fakeProductRepository) UpdateProduct(domain.Product, int) (domain.Product, error) {
	panic("não usado")
}

func TestProductService_GetProductById(t *testing.T) {
	errNotFound := errors.New("product not found")

	tests := []struct {
		name    string
		give    int
		repo    func(id int) (domain.Product, error) // o que o "banco" devolve
		want    domain.Product
		wantErr error
	}{
		{
			name: "produto encontrado",
			give: 1,
			repo: func(int) (domain.Product, error) {
				return domain.Product{ID: 1, Name: "Nintendo Switch", Price: 1999.9}, nil
			},
			want: domain.Product{ID: 1, Name: "Nintendo Switch", Price: 1999.9},
		},
		{
			name: "repositório retorna erro",
			give: 99,
			repo: func(int) (domain.Product, error) {
				return domain.Product{}, errNotFound
			},
			wantErr: errNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeProductRepository{getByID: tt.repo}
			svc := services.NewProductService(repo)

			got, err := svc.GetProductById(tt.give)

			// errors.Is respeita o wrapping de erros (%w).
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("GetProductById(%d) error = %v, want %v", tt.give, err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetProductById(%d) mismatch (-want +got):\n%s", tt.give, diff)
			}
		})
	}
}
