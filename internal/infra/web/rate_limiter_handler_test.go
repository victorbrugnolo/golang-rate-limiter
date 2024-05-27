package web

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorbrugnolo/golang-rate-limiter/internal/entity"
	"github.com/victorbrugnolo/golang-rate-limiter/internal/infra/database"
)

func TestHandleRateLimitAndGetIsAllowedResponse(t *testing.T) {
	ctx := context.Background()
	clientID := "abc"
	rateLimiterConfig := &entity.RateLimiterConfig{Limit: 10, Window: 10, BlockWindow: 10}
	rateLimiterData := &entity.RateLimiterData{Requests: 0, LastSeen: time.Now(), Blocked: false}
	repository := &database.RateLimiterDataRepositoryMock{}

	repository.On("Find", ctx, clientID).Return(rateLimiterData, nil)
	repository.On("Save", ctx, clientID, rateLimiterData).Return(nil)

	rateLimiterHandler := NewRateLimiterHandler(rateLimiterConfig, repository)
	result := rateLimiterHandler.HandleRateLimit(ctx, clientID)

	assert.True(t, result)
	repository.AssertExpectations(t)
}

func TestHandleRateLimitAndGetIsNotAllowedResponseWhenClientIsBloked(t *testing.T) {
	ctx := context.Background()
	clientID := "abc"
	rateLimiterConfig := &entity.RateLimiterConfig{Limit: 10, Window: time.Duration(5) * time.Second, BlockWindow: time.Duration(10) * time.Second}
	rateLimiterData := &entity.RateLimiterData{Requests: 0, LastSeen: time.Now(), Blocked: true}
	repository := &database.RateLimiterDataRepositoryMock{}

	repository.On("Find", ctx, clientID).Return(rateLimiterData, nil)

	rateLimiterHandler := NewRateLimiterHandler(rateLimiterConfig, repository)
	result := rateLimiterHandler.HandleRateLimit(ctx, clientID)

	assert.False(t, result)
	repository.AssertExpectations(t)
	repository.AssertNotCalled(t, "Save", ctx, clientID, rateLimiterData)
}

func TestHandleRateLimitAndGetIsNotAllowedWhenGetErrorOnFindRateLimiterData(t *testing.T) {
	ctx := context.Background()
	clientID := "abc"
	rateLimiterConfig := &entity.RateLimiterConfig{Limit: 10, Window: time.Duration(5) * time.Second, BlockWindow: time.Duration(10) * time.Second}
	repository := &database.RateLimiterDataRepositoryMock{}

	repository.On("Find", ctx, clientID).Return(&entity.RateLimiterData{}, errors.New("error"))

	rateLimiterHandler := NewRateLimiterHandler(rateLimiterConfig, repository)
	result := rateLimiterHandler.HandleRateLimit(ctx, clientID)

	assert.False(t, result)
	repository.AssertExpectations(t)
	repository.AssertNotCalled(t, "Save", ctx, clientID, nil)
}

func TestHandleRateLimitAndGetIsAllowedResponseWhenClientIsBlokedOutOfBlockWindow(t *testing.T) {
	ctx := context.Background()
	clientID := "abc"
	rateLimiterConfig := &entity.RateLimiterConfig{Limit: 10, Window: time.Duration(5) * time.Second, BlockWindow: time.Duration(10) * time.Microsecond}
	rateLimiterData := &entity.RateLimiterData{Requests: 0, LastSeen: time.Now(), Blocked: true}
	repository := &database.RateLimiterDataRepositoryMock{}

	repository.On("Find", ctx, clientID).Return(rateLimiterData, nil)
	repository.On("Save", ctx, clientID, rateLimiterData).Return(nil)

	rateLimiterHandler := NewRateLimiterHandler(rateLimiterConfig, repository)
	result := rateLimiterHandler.HandleRateLimit(ctx, clientID)

	assert.True(t, result)
	repository.AssertExpectations(t)
}

func TestHandleRateLimitAndGetIsNotAllowedResponseWhenClientRequestsGreatherThanLimit(t *testing.T) {
	ctx := context.Background()
	clientID := "abc"
	rateLimiterConfig := &entity.RateLimiterConfig{Limit: 10, Window: time.Duration(5) * time.Second, BlockWindow: time.Duration(10) * time.Microsecond}
	rateLimiterData := &entity.RateLimiterData{Requests: 10, LastSeen: time.Now(), Blocked: false}
	repository := &database.RateLimiterDataRepositoryMock{}

	repository.On("Find", ctx, clientID).Return(rateLimiterData, nil)
	repository.On("Save", ctx, clientID, rateLimiterData).Return(nil)

	rateLimiterHandler := NewRateLimiterHandler(rateLimiterConfig, repository)
	result := rateLimiterHandler.HandleRateLimit(ctx, clientID)

	assert.False(t, result)
	repository.AssertExpectations(t)
}
