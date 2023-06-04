package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/williamluisan/vrd_pdfgenerator/repository"
)

type Vrd_mailer repository.Vrd_mailer
type Vrd_mailer_return repository.Vrd_mailer_return

func (vrd *Vrd_mailer) Send() (string, error) {
	url := os.Getenv("VRD_MAILER_URL") + "/send"
	reqBody, err := json.Marshal(vrd)
	if err != nil {
		log.Println(err.Error())
	}
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    if err != nil {
        return "", err
    }

	res, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }

	defer res.Body.Close()

	statusCode := res.StatusCode
	if statusCode != http.StatusOK {
		var vrd_mailer_return Vrd_mailer_return
		err = json.NewDecoder(res.Body).Decode(&vrd_mailer_return)
		if err != nil {
			log.Println(err.Error())
		}

		return vrd_mailer_return.Message, errors.New(strconv.Itoa(statusCode))
	}

	return "", nil
}