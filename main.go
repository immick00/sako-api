package main

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/immick00/sako-api/db"
	"github.com/immick00/sako-api/logger"
	"github.com/immick00/sako-api/menus"
	"github.com/immick00/sako-api/places"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("no .env file found")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Log.Error("failed to connect to database", "error", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	menusService := menus.New(queries)

	placesService := places.New("", "")
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logger.Log.Error("request error", "method", c.Request().Method, "path", c.Request().URL.Path, "error", err)
				return err
			}
			if code := c.Response().Status; code >= 400 {
				logger.Log.Warn("unsuccessful request", "method", c.Request().Method, "path", c.Request().URL.Path, "status", code)
			}
			return nil
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/nearby", func(c echo.Context) error {
		latStr := c.QueryParam("lat")
		lonStr := c.QueryParam("lon")

		if latStr == "" || lonStr == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "lat and lon are required"})
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid lat"})
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid lon"})
		}

		result, err := placesService.GetRestaurantsAround(lat, lon)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, result)
	})

	e.POST("/menus", func(c echo.Context) error {
		var body struct {
			Restaurants []string `json:"restaurants"`
		}
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}
		if len(body.Restaurants) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "restaurants is required"})
		}

		menus, err := menusService.GetMenus(c.Request().Context(), body.Restaurants)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, menus)
	})

	e.POST("/onboarding", func(c echo.Context) error {
		var body struct {
			Goal          *string  `json:"goal"`
			Weight        float64  `json:"weight"`
			HeightFeet    string   `json:"heightFeet"`
			HeightInches  string   `json:"heightInches"`
			AgeRange      *string  `json:"ageRange"`
			DaysPerWeek   *string  `json:"daysPerWeek"`
			ActivityLevel *string  `json:"activityLevel"`
			Cravings      []string `json:"cravings"`
			Dislikes      []string `json:"dislikes"`
		}
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}
		if body.Weight == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "weight is required"})
		}
		if body.HeightFeet == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "heightFeet is required"})
		}
		if body.HeightInches == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "heightInches is required"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	logger.Log.Info("starting server", "port", 8080)
	e.Logger.Fatal(e.Start(":8080"))
}
