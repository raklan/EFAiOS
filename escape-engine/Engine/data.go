package Engine

import (
	"encoding/json"
	"escape-engine/Models"
	"fmt"
	"os"
)

// Yes, I know. I just REALLY didn't want to bring in an entire database JUST for this and Redis shouldn't be used for it
func SaveMapToDB(m Models.GameMap) (Models.GameMap, error) {
	if m.Id == "" {
		m.Id = GenerateId()
	}

	asJson, err := json.Marshal(m)
	if err != nil {
		return m, err
	}

	filename := "map_" + m.Id + ".json"
	f, err := os.Create(fmt.Sprintf("./maps/%s", filename))
	if err != nil {
		fmt.Println("FML")
		f.Close()
		return m, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		return m, err
	}

	return m, nil
}

func GetMapFromDB(mapId string) (Models.GameMap, error) {
	data, err := os.ReadFile(fmt.Sprintf("./maps/map_%s.json", mapId))
	if err != nil {
		return Models.GameMap{}, err
	}

	parsed := Models.GameMap{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return parsed, err
	}

	return parsed, nil
}
