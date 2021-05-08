package server

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func echoMonitoring() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()
			timeTaken := time.Since(start).Seconds()
			size := res.Size

			label := req.Method + " "
			if c.Path() != "" {
				label += c.Path()
			} else {
				label += req.URL.Path
			}

			first := false
			for k, v := range req.URL.Query() {
				if !first {
					label += "?"
				}
				if first {
					label += "&"
				}

				label += fmt.Sprintf("%s=%s", k, v)
				first = true
			}

			l := struct {
				ID        primitive.ObjectID
				Label     string
				Path      string
				Status    int
				Size      int64
				IP        string
				TimeTaken float64
				CreatedAt time.Time
			}{
				ID:        primitive.NewObjectID(),
				Label:     label,
				Path:      req.URL.Path,
				Status:    res.Status,
				Size:      size,
				IP:        req.RemoteAddr,
				TimeTaken: timeTaken,
				CreatedAt: time.Now().UTC(),
			}
			db := database.GetDB()
			logCollection := db.Collection("logs")
			if l.Status != 204 {
				if _, err := logCollection.InsertOne(context.Background(), l); err != nil {
					logger.Log.Errorln(err)
				}
			}
			return nil
		}
	}
}
