package moyklass

import (
	"app/internal/model"
	"net/http"
)

var MoyClass *MoyClassKit

func InitNoyKlass(mux *http.ServeMux) error {
	moyClass, err := NewMoyClassKit(mux, model.ConfigFile.MoyClassConfig.WebHouckUrl, model.ConfigFile.MoyClassConfig.ApiKey)
	if err != nil {
		return err
	}
	MoyClass = moyClass

	return nil
}
