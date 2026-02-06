package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/duynhne/notification-service/internal/core/domain"
	logicv1 "github.com/duynhne/notification-service/internal/logic/v1"
	"github.com/duynhne/notification-service/middleware"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var notificationService = logicv1.NewNotificationService()

func SendEmail(c *gin.Context) {
	ctx, span := middleware.StartSpan(c.Request.Context(), "http.request", trace.WithAttributes(
		attribute.String("layer", "web"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.Request.URL.Path),
	))
	defer span.End()

	zapLogger := middleware.GetLoggerFromGinContext(c)

	var req domain.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.SetAttributes(attribute.Bool("request.valid", false))
		span.RecordError(err)
		zapLogger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.Bool("request.valid", true))
	notification, err := notificationService.SendEmail(ctx, req)
	if err != nil {
		span.RecordError(err)
		zapLogger.Error("Failed to send email", zap.Error(err))

		switch {
		case errors.Is(err, logicv1.ErrInvalidRecipient):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient"})
		case errors.Is(err, logicv1.ErrDeliveryFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Delivery failed"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	zapLogger.Info("Email sent", zap.String("notification_id", notification.ID))
	c.JSON(http.StatusOK, notification)
}

func SendSMS(c *gin.Context) {
	ctx, span := middleware.StartSpan(c.Request.Context(), "http.request", trace.WithAttributes(
		attribute.String("layer", "web"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.Request.URL.Path),
	))
	defer span.End()

	zapLogger := middleware.GetLoggerFromGinContext(c)

	var req domain.SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.SetAttributes(attribute.Bool("request.valid", false))
		span.RecordError(err)
		zapLogger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.Bool("request.valid", true))
	notification, err := notificationService.SendSMS(ctx, req)
	if err != nil {
		span.RecordError(err)
		zapLogger.Error("Failed to send SMS", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	zapLogger.Info("SMS sent", zap.String("notification_id", notification.ID))
	c.JSON(http.StatusOK, notification)
}

// ListNotifications handles GET /api/v1/notifications
func ListNotifications(c *gin.Context) {
	ctx, span := middleware.StartSpan(c.Request.Context(), "http.request", trace.WithAttributes(
		attribute.String("layer", "web"),
		attribute.String("api.version", "v1"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.Request.URL.Path),
	))
	defer span.End()

	zapLogger := middleware.GetLoggerFromGinContext(c)

	// Get user_id from auth middleware (falls back to "1" for demo)
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "1"
	}

	notifications, err := notificationService.ListNotifications(ctx, userID)
	if err != nil {
		span.RecordError(err)
		zapLogger.Error("Failed to list notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	zapLogger.Info("Notifications listed", zap.Int("count", len(notifications)))
	c.JSON(http.StatusOK, notifications)
}

// handleNotificationByID is a shared handler for operations on a single notification by ID.
// It extracts common boilerplate (span setup, ID extraction, error handling) to avoid duplication.
func handleNotificationByID(
	c *gin.Context,
	action func(ctx context.Context, id string) (*domain.Notification, error),
	successLog string,
) {
	ctx, span := middleware.StartSpan(c.Request.Context(), "http.request", trace.WithAttributes(
		attribute.String("layer", "web"),
		attribute.String("api.version", "v1"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.Request.URL.Path),
	))
	defer span.End()

	zapLogger := middleware.GetLoggerFromGinContext(c)
	id := c.Param("id")
	span.SetAttributes(attribute.String("notification.id", id))

	notification, err := action(ctx, id)
	if err != nil {
		span.RecordError(err)
		zapLogger.Error(successLog+" failed", zap.Error(err))

		switch {
		case errors.Is(err, logicv1.ErrNotificationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	zapLogger.Info(successLog, zap.String("notification_id", id))
	c.JSON(http.StatusOK, notification)
}

// GetNotification handles GET /api/v1/notifications/:id
func GetNotification(c *gin.Context) {
	handleNotificationByID(c, notificationService.GetNotification, "Notification retrieved")
}

// MarkAsRead handles PATCH /api/v1/notifications/:id
func MarkAsRead(c *gin.Context) {
	handleNotificationByID(c, notificationService.MarkAsRead, "Notification marked as read")
}
