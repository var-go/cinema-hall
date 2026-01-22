package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	userSvc := getEnv("USER_SERVICE_URL", "http://localhost:8080")
	movieSvc := getEnv("MOVIE_SERVICE_URL", "http://localhost:8083")
	cinemaSvc := getEnv("CINEMA_SERVICE_URL", "http://localhost:8081")
	bookingSvc := getEnv("BOOKING_SERVICE_URL", "http://localhost:8082")

	router := gin.Default()

	router.POST("/api/auth/register", func(c *gin.Context) {

		body, _ := io.ReadAll(c.Request.Body)
		resp, err := http.Post(strings.TrimRight(userSvc, "/")+"/auth/register", "application/json", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/auth/login", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		resp, err := http.Post(strings.TrimRight(userSvc, "/")+"/auth/login", "application/json", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/movies", func(c *gin.Context) {
		resp, err := http.Get(strings.TrimRight(movieSvc, "/") + "/movies")
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "movie service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/movies/:id", func(c *gin.Context) {
		id := c.Param("id")
		resp, err := http.Get(strings.TrimRight(movieSvc, "/") + "/movies/" + id)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "movie service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/movies", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		resp, err := http.Post(strings.TrimRight(movieSvc, "/")+"/movies", "application/json", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "movie service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/sessions", func(c *gin.Context) {
		resp, err := http.Get(strings.TrimRight(cinemaSvc, "/") + "/sessions")
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/sessions/:id", func(c *gin.Context) {
		id := c.Param("id")
		resp, err := http.Get(strings.TrimRight(cinemaSvc, "/") + "/sessions/" + id)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/sessions", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		resp, err := http.Post(strings.TrimRight(cinemaSvc, "/")+"/sessions", "application/json", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/bookings", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		resp, err := http.Get(strings.TrimRight(bookingSvc, "/") + "/booking")
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/bookings/:id", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		id := c.Param("id")
		resp, err := http.Get(strings.TrimRight(bookingSvc, "/") + "/booking/" + id)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/bookings", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		body, _ := io.ReadAll(c.Request.Body)
		resp, err := http.Post(strings.TrimRight(bookingSvc, "/")+"/booking", "application/json", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/sessions/:id/aggregate", func(c *gin.Context) {
		id := c.Param("id")

		resp, err := http.Get(strings.TrimRight(cinemaSvc, "/") + "/sessions/" + id)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			c.Data(resp.StatusCode, "application/json", b)
			return
		}
		var session map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&session)

		movieID := toIDString(session["movie_id"])
		var movie map[string]interface{}
		if movieID != "" {
			r2, err := http.Get(strings.TrimRight(movieSvc, "/") + "/movies/" + movieID)
			if err == nil && r2 != nil {
				defer r2.Body.Close()
				if r2.StatusCode < 400 {
					_ = json.NewDecoder(r2.Body).Decode(&movie)
				}
			}
		}

		hallID := toIDString(session["hall_id"])
		var hall map[string]interface{}
		if hallID != "" {
			r3, err := http.Get(strings.TrimRight(cinemaSvc, "/") + "/halls/" + hallID)
			if err == nil && r3 != nil {
				defer r3.Body.Close()
				if r3.StatusCode < 400 {
					_ = json.NewDecoder(r3.Body).Decode(&hall)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"session": session, "movie": movie, "hall": hall})
	})

	port := getEnv("PORT", "8085")
	router.Run(":" + port)
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func toIDString(v interface{}) string {
	switch t := v.(type) {
	case float64:
		return strings.TrimRight(strings.TrimRight(fmtFloat(t), ".0"), ".")
	case string:
		return t
	default:
		return ""
	}
}

func fmtFloat(f float64) string {

	return strings.TrimSuffix(strings.TrimSuffix(strings.TrimSpace(strings.ReplaceAll(strings.TrimSpace(""), "", "")), "."), ".0")
}

func validateJWT(c *gin.Context) bool {
	secret := []byte(os.Getenv("JWT_SECRET"))
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return false
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return false
	}
	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return false
	}
	return true
}
