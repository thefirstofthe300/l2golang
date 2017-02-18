package characters

import (
	"encoding/json"
	"fmt"
	"myHTTP"
	"myJSON"
	"net/http"
	"regexp"
	"strconv"
)

type character struct {
	Name  string `json: "name"`
	Race  string `json: "race"`
	Level int    `json: "level"`
}

const (
	listCharPattern      = `^/api/characters$`
	addCharPattern       = `^/api/characters/add$`
	listCharLevelPattern = `^/api/characters/\d+$`
)

type handler func(*http.Request) ([]byte, error)

var dispatch = map[string]handler{
	listCharPattern:      listChars,
	addCharPattern:       AddCharacter,
	listCharLevelPattern: listLevelChars,
}

func Dispatcher(r *http.Request) ([]byte, error) {
	p := r.URL.Path
	for k, v := range dispatch {
		ok, err := regexp.MatchString(k, p)
		if err != nil {
			return nil, err
		}
		if ok {
			return v(r)
		}
	}
	return nil, myHTTP.NotFoundErr
}

func loadChars() ([]character, error) {
	b, err := myJSON.LoadJSON("characters")
	if err != nil {
		return nil, err
	}
	var s []character
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}

func saveChars(c []character) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return myJSON.SaveJSON("chracters", b)
}

func listChars(r *http.Request) ([]byte, error) {
	return myJSON.LoadJSON("characters")
}

func getChar(r *http.Request) ([]byte, error) {
	x, err := loadChars()
	if err != nil {
		return nil, err
	}
	char := r.URL.Path[len("/api/character/"):]
	for _, e := range x {
		if e.Name == char {
			return json.Marshal(e)
		}
	}
	return []byte(fmt.Sprintf("The character %s has not been created yet.", char)), nil

}
func listLevelChars(r *http.Request) ([]byte, error) {
	t, err := loadChars()
	if err != nil {
		return nil, err
	}
	s, err := strconv.Atoi(r.URL.Path[len("/api/characters/"):])
	if err != nil {
		return nil, err
	}
	var o []character
	for _, e := range t {
		if e.Level == s {
			o = append(o, e)
		}
	}
	if len(o) != 0 {
		return json.Marshal(o)
	}
	return []byte(fmt.Sprintf("There are no characters at level %d ...yet", s)), nil

}
func AddCharacter(r *http.Request) ([]byte, error) {
	return nil, nil
}
