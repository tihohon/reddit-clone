package helpers

import (
	"fmt"
	"net/http"
)

func WriteBadRequest(w http.ResponseWriter, text string) {

	http.Error(w, fmt.Sprintf(`{"error": %v}`, text), http.StatusBadRequest)

}
