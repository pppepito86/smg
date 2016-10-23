package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"request/util"
	"strconv"
)

type WeekPointsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *WeekPointsHandler) Execute() error {
	h.W.Header().Set("Content-Type", "application/json")

	h.R.ParseForm()
	var userId int64
	var err error
	if len(h.R.Form["id"]) == 1 {
		userId, err = strconv.ParseInt(h.R.Form["id"][0], 10, 64)
		if err != nil {
			userId = h.User.Id
		}
	} else {
		userId = h.User.Id
	}
	if h.User.RoleName != "admin" && userId != h.User.Id {
		return nil
	}

	Response := db.GetPointsPerWeek(userId)

	json, err := json.Marshal(Response)
	fmt.Println("json", json)
	if err != nil {
		fmt.Println("err", err.Error())
		http.Error(h.W, err.Error(), http.StatusInternalServerError)
		return nil
	}
	h.W.Write(json)

	return nil
}
