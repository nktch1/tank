package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseYandexURL = "https://yandex.ru/search/touch/" +
	"?service=www.yandex&ui=webmobileapp.yandex&numdoc=50&lr=213&p=0&text=%s"

func Search(ctx context.Context, query string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf(baseYandexURL, query))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
