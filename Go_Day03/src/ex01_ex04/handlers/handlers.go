package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"ex01_04/db"
)

const pageSize = 2

var htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Places</title>
</head>
<body>
	<h5>Total: {{.Total}}</h5>
	<ul>
		{{range .Places}}
		<li>
			<div>{{.Name}}</div>
			<div>{{.Address}}</div>
			<div>{{.Phone}}</div>
		</li>
		{{end}}
	</ul>

	{{if ge .PrevPage 0}}<a href="/?page={{.PrevPage}}">Previous</a>{{end}}
	{{if ge .NextPage 0}}<a href="/?page={{.NextPage}}">Next</a>{{end}}
	{{if gt .Page 0}}<a href="/?page=0">First</a>{{end}}
	{{if lt .Page .LastPage}}<a href="/?page={{.LastPage}}">Last</a>{{end}}
</body>
</html>
`

var tmpl = template.Must(template.New("htmlTemplate").Parse(htmlTemplate))

func RenderPage(w http.ResponseWriter, r *http.Request, store db.Store) {
	page, places, total, err := getPlacesFromDB(r, store)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	lastPage := (total - 1) / pageSize
	prevPage := -1
	if page > 0 {
		prevPage = page - 1
	}

	nextPage := -1
	if (page+1)*pageSize < total {
		nextPage = page + 1
	}

	data := struct {
		Total    int
		Places   []db.Place
		PrevPage int
		NextPage int
		LastPage int
		Page     int
	}{
		Total:    total,
		Places:   places,
		PrevPage: prevPage,
		NextPage: nextPage,
		LastPage: lastPage,
		Page:     page,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func RenderJSON(w http.ResponseWriter, r *http.Request, store db.Store) {
	page, places, total, err := getPlacesFromDB(r, store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Total    int        `json:"total"`
		Places   []db.Place `json:"places"`
		PrevPage int        `json:"prev_page"`
		NextPage int        `json:"next_page"`
		LastPage int        `json:"last_page"`
	}{
		Total:    total,
		Places:   places,
		PrevPage: page - 1,
		NextPage: page + 1,
		LastPage: (total - 1) / pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RenderRecommendations(w http.ResponseWriter, r *http.Request, store db.Store) {
	latParam := r.URL.Query().Get("lat")
	lonParam := r.URL.Query().Get("lon")

	lat, err1 := strconv.ParseFloat(latParam, 64)
	lon, err2 := strconv.ParseFloat(lonParam, 64)

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid latitude or longitude", http.StatusBadRequest)
		return
	}

	places, err := store.GetClosestPlaces(lat, lon, 3)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":   "Recommendation",
		"places": places,
	})
}

func getPlacesFromDB(r *http.Request, store db.Store) (int, []db.Place, int, error) {
	page := 0
	elementsCount := int(store.GetElementsCount())
	pageParam := r.URL.Query().Get("page")

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err != nil || p < 0 || p*pageSize >= elementsCount {
			return 0, nil, 0, fmt.Errorf("invalid 'page' value: '%s'", pageParam)
		}
		page = p
	}

	offset := page * pageSize
	places, total, err := store.GetPlaces(pageSize, offset)
	if err != nil {
		return 0, nil, 0, err
	}

	return page, places, total, nil
}
