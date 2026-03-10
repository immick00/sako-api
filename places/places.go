package places

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/immick00/sako-api/logger"
)

const BASE_URL = "https://places.googleapis.com"
const RADIUS = 5000 // In meters

// NOTE: need to redo this
var restaurantAliases = map[string]string{
	"starbucks":  "starbucks",
	"mcdonalds":  "mcdonalds",
	"mcdonald's": "mcdonalds",
}

func normalizeRestaurantName(name string) string {
	lower := strings.ToLower(name)
	for keyword, canonical := range restaurantAliases {
		if strings.Contains(lower, keyword) {
			return canonical
		}
	}
	return lower
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DisplayName struct {
	Text string `json:"text"`
}

type Place struct {
	ID               string      `json:"id"`
	DisplayName      DisplayName `json:"displayName"`
	Image            string      `json:"imageUrl"`
	FormattedAddress string      `json:"formattedAddress"`
	Location         Location    `json:"location"`
	Rating           float64     `json:"rating"`
}

type PlacesResponse struct {
	Places        []Place `json:"places"`
	NextPageToken string  `json:"nextPageToken"`
}

type PlacesService struct {
	apiKey          string
	distributionUrl string
}

func New(apiKey, distributionUrl string) *PlacesService {
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	}
	if distributionUrl == "" {
		distributionUrl = os.Getenv("DISTRIBUTION_URL")
	}
	return &PlacesService{apiKey: apiKey, distributionUrl: distributionUrl}
}

func (s *PlacesService) searchText(query string, lat, lon float64, pageToken string) (*PlacesResponse, error) {
	body := map[string]any{
		"textQuery":      query,
		"maxResultCount": 1, // WARNING: might need to change this logic
		"locationBias": map[string]any{
			"circle": map[string]any{
				"center": map[string]any{"latitude": lat, "longitude": lon},
				"radius": RADIUS,
			},
		},
	}
	if pageToken != "" {
		body["pageToken"] = pageToken
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/places:searchText", BASE_URL), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", s.apiKey)
	req.Header.Set("X-Goog-FieldMask", "nextPageToken,places.id,places.displayName,places.formattedAddress,places.location,places.rating")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *PlacesService) GetRestaurantsAround(lat, lon float64) ([]Place, error) {
	// restaurants := []string{"mcdonalds"}
	restaurants := []string{"McDonald's", "Subway"}
	// restaurants := []string{"In-N-Out Burger"}

	var (
		mu        sync.Mutex
		wg        sync.WaitGroup
		allPlaces = []Place{}
	)

	for _, restaurant := range restaurants {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			result, err := s.searchText(query, lat, lon, "")
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			mu.Lock()
			allPlaces = append(allPlaces, result.Places...)
			mu.Unlock()
		}(restaurant)
	}
	wg.Wait()

	seen := make(map[string]struct{})
	unique := []Place{}
	for _, p := range allPlaces {
		// make sure its in our radius/remove dupes
		if _, ok := seen[p.ID]; !ok && haversine(lat, lon, p.Location.Latitude, p.Location.Longitude) <= RADIUS {
			seen[p.ID] = struct{}{}
			p.DisplayName.Text = normalizeRestaurantName(p.DisplayName.Text)
			p.Image = s.distributionUrl + "/images/logos/" + p.DisplayName.Text + ".png"
			unique = append(unique, p)
		}
	}

	return sortByDistance(unique, lat, lon), nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371000 // earth radius in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	dPhi := (lat2 - lat1) * math.Pi / 180
	dLambda := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dPhi/2)*math.Sin(dPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(dLambda/2)*math.Sin(dLambda/2)
	return r * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func sortByDistance(places []Place, lat, lon float64) []Place {
	sort.Slice(places, func(i, j int) bool {
		di := haversine(lat, lon, places[i].Location.Latitude, places[i].Location.Longitude)
		dj := haversine(lat, lon, places[j].Location.Latitude, places[j].Location.Longitude)
		return di < dj
	})
	return places
}
