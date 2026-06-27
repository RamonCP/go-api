// Teste white-box (package http): precisa enxergar toProductResponse, que é
// uma função não-exportada. Por isso ficamos no mesmo pacote, e não em http_test.
package http

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"go-api/internal/core/domain"
)

func TestToProductResponse(t *testing.T) {
	tests := []struct {
		name string
		give domain.Product
		want ProductResponse
	}{
		{
			name: "mapeia todos os campos",
			give: domain.Product{ID: 1, Name: "Nintendo Switch", Price: 1999.9},
			want: ProductResponse{ID: 1, Name: "Nintendo Switch", Price: 1999.9},
		},
		{
			name: "produto zero-value vira response zero-value",
			give: domain.Product{},
			want: ProductResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toProductResponse(tt.give)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("toProductResponse(%+v) mismatch (-want +got):\n%s", tt.give, diff)
			}
		})
	}
}
