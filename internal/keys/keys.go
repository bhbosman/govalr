package keys

import (
	"encoding/json"
	"fmt"
	"go.uber.org/fx"
	"os"
	"os/user"
	"path"
)

type ValrConnectionSettings struct {
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

func ProvideValrConnectionSettings() fx.Option {

	return fx.Provide(
		func() (*ValrConnectionSettings, error) {
			data := &ValrConnectionSettings{}
			current, err := user.Current()
			if err != nil {
				return nil,  err
			}
			f, err := os.Open(fmt.Sprintf(path.Join(current.HomeDir, ".valr", "keys.json")))
			if err != nil {
				return nil, err
			}

			defer func() {
				_ = f.Close()
			}()
			decoder := json.NewDecoder(f)
			err = decoder.Decode(data)
			if err != nil {
				return nil, err
			}
			return data, nil

			//return &ValrConnectionSettings{
			//	ApiKey:    "42299dc9cddcffbed10db28aa153027899ecfb7bf849fd026af1c42e5aba01d1",
			//	SecretKey: "1dbd3183642fdacd180d3c383a98b45de110620fb30a20efb521fdcddc08bc86",
			//}, nil
		})
}
