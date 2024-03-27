package smartphones

import (
	"context"
	"encoding/json"
	"fmt"
	"pycoala/pycoala/config"
	"pycoala/pycoala/database"
	bothttp "pycoala/pycoala/utils/helpers"
)

type DeviceInfo struct {
	Link        string `json:"link"`
	Img         string `json:"img"`
	Description string `json:"description"`
}

// getDataFromUrl fetches data from a given URL, checks the database and updates if needed
func getDataFromUrl(ctx context.Context, url string) (map[string]interface{}, error) {
	// Check if data exists in the database first
	var data map[string]interface{}
	err := findURLData(ctx, url, &data) // Call the new findURLData function

	// If data exists, return it directly
	if data != nil {
		return data, nil
	}

	// If data doesn't exist, fetch from URL and add to database
	resp := bothttp.RequestGET(url, bothttp.RequestGETParams{})

	body := resp.Body()

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	err = addURLData(ctx, url, data) // Call the new addURLData function
	if err != nil {
		return data, fmt.Errorf("failed to add data to database: %w", err)
	}

	return data, nil
}

// findURLData checks if data for a specific URL exists in the database
func findURLData(ctx context.Context, url string, data *map[string]interface{}) error {
	query := `SELECT data FROM url_data WHERE url = ?;`
	row := database.DB.QueryRowContext(ctx, query, url)
	err := row.Scan(data)
	return err
}

// addURLData adds data for a specific URL to the database
func addURLData(ctx context.Context, url string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to json: %w", err)
	}
	query := `INSERT INTO url_data (url, data) VALUES (?, ?);`
	_, err = database.DB.ExecContext(ctx, query, url, string(jsonData))
	return err
}

func searchDevice(ctx context.Context, searchValue string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/search?value=%s", config.GSMArenaAPI, searchValue)
	return getDataFromUrl(ctx, url)
}

func getDevice(ctx context.Context, device string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/getdevice?value=%s", config.GSMArenaAPI, device)
	return getDataFromUrl(ctx, url)
}
