package server

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func echoMonitoring() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			if err := next(c); err != nil {
				c.Error(err)
			}

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
			// db := database.DB()
			// if err := db.Table(l.TableName()).Create(&l).Error; err != nil {
			// 	logger.Log.Errorln(err)
			// }
			logger.Printf("%+v", l)
			return nil
		}
	}
}
