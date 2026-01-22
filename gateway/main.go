package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	userSvc := getEnv("USER_SERVICE_URL", "http://localhost:8080")
	movieSvc := getEnv("MOVIE_SERVICE_URL", "http://localhost:8083")
	cinemaSvc := getEnv("CINEMA_SERVICE_URL", "http://localhost:8081")
	bookingSvc := getEnv("BOOKING_SERVICE_URL", "http://localhost:8082")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	router := gin.Default()

	router.POST("/api/auth/register", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		req, err := http.NewRequest("POST", strings.TrimRight(userSvc, "/")+"/auth/register", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/auth/login", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		req, err := http.NewRequest("POST", strings.TrimRight(userSvc, "/")+"/auth/login", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/movies", func(c *gin.Context) {
		req, err := http.NewRequest("GET", strings.TrimRight(movieSvc, "/")+"/movies", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.URL.RawQuery = c.Request.URL.RawQuery

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "movie service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/movies/:id", func(c *gin.Context) {
		id := c.Param("id")
		req, err := http.NewRequest("GET", strings.TrimRight(movieSvc, "/")+"/movies/"+id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "movie service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/movies", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		req, err := http.NewRequest("POST", strings.TrimRight(movieSvc, "/")+"/movies", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "movie service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/sessions", func(c *gin.Context) {
		req, err := http.NewRequest("GET", strings.TrimRight(cinemaSvc, "/")+"/sessions", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.URL.RawQuery = c.Request.URL.RawQuery

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/sessions/:id", func(c *gin.Context) {
		id := c.Param("id")
		req, err := http.NewRequest("GET", strings.TrimRight(cinemaSvc, "/")+"/sessions/"+id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/sessions", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		req, err := http.NewRequest("POST", strings.TrimRight(cinemaSvc, "/")+"/sessions", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/bookings", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}

		req, err := http.NewRequest("GET", strings.TrimRight(bookingSvc, "/")+"/bookings", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.URL.RawQuery = c.Request.URL.RawQuery

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/bookings/:id", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		id := c.Param("id")

		req, err := http.NewRequest("GET", strings.TrimRight(bookingSvc, "/")+"/bookings/"+id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/bookings", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		req, err := http.NewRequest("POST", strings.TrimRight(bookingSvc, "/")+"/bookings", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.PATCH("/api/bookings/:id", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		id := c.Param("id")

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		req, err := http.NewRequest("PATCH", strings.TrimRight(bookingSvc, "/")+"/bookings/"+id, bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.DELETE("/api/bookings/:id", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		id := c.Param("id")

		req, err := http.NewRequest("DELETE", strings.TrimRight(bookingSvc, "/")+"/bookings/"+id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/bookings/:id/confirm", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		id := c.Param("id")

		req, err := http.NewRequest("POST", strings.TrimRight(bookingSvc, "/")+"/bookings/"+id+"/confirm", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.POST("/api/bookings/:id/cancel", func(c *gin.Context) {
		if !validateJWT(c) {
			return
		}
		id := c.Param("id")

		req, err := http.NewRequest("POST", strings.TrimRight(bookingSvc, "/")+"/bookings/"+id+"/cancel", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "booking service unavailable"})
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
			return
		}
		c.Data(resp.StatusCode, "application/json", b)
	})

	router.GET("/api/sessions/:id/aggregate", func(c *gin.Context) {
		id := c.Param("id")

		req, err := http.NewRequest("GET", strings.TrimRight(cinemaSvc, "/")+"/sessions/"+id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "cinema service unavailable"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
				return
			}
			c.Data(resp.StatusCode, "application/json", b)
			return
		}

		var session map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to decode session"})
			return
		}

		movieID := toIDString(session["movie_id"])
		var movie map[string]interface{}
		if movieID != "" {
			req2, err := http.NewRequest("GET", strings.TrimRight(movieSvc, "/")+"/movies/"+movieID, nil)
			if err == nil {
				r2, err := httpClient.Do(req2)
				if err == nil && r2 != nil {
					defer r2.Body.Close()
					if r2.StatusCode < 400 {
						_ = json.NewDecoder(r2.Body).Decode(&movie)
					}
				}
			}
		}

		hallID := toIDString(session["hall_id"])
		var hall map[string]interface{}
		if hallID != "" {
			req3, err := http.NewRequest("GET", strings.TrimRight(cinemaSvc, "/")+"/halls/"+hallID, nil)
			if err == nil {
				r3, err := httpClient.Do(req3)
				if err == nil && r3 != nil {
					defer r3.Body.Close()
					if r3.StatusCode < 400 {
						_ = json.NewDecoder(r3.Body).Decode(&hall)
					}
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
	return fmt.Sprintf("%.0f", f)
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
