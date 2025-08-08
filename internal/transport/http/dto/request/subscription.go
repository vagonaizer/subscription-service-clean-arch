package request

import (
	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required" example:"Yandex Plus" minLength:"1" maxLength:"255"`
	Price       int    `json:"price" binding:"required,min=1,max=1000000" example:"400"`
	UserID      string `json:"user_id" binding:"required,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" binding:"required" example:"07-2025" pattern:"^(0[1-9]|1[0-2])-[0-9]{4}$"`
	EndDate     string `json:"end_date,omitempty" example:"12-2025" pattern:"^(0[1-9]|1[0-2])-[0-9]{4}$"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Netflix Premium" minLength:"1" maxLength:"255"`
	Price       *int    `json:"price,omitempty" minimum:"1" maximum:"1000000" example:"799"`
	StartDate   *string `json:"start_date,omitempty" example:"08-2025" pattern:"^(0[1-9]|1[0-2])-[0-9]{4}$"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025" pattern:"^(0[1-9]|1[0-2])-[0-9]{4}$"`
}

type GetSubscriptionRequest struct {
	ID string `json:"id" path:"id"`
}

type DeleteSubscriptionRequest struct {
	ID string `json:"id" path:"id"`
}

type GetUserSubscriptionsRequest struct {
	UserID string `json:"user_id" path:"user_id"`
	Limit  int    `json:"limit" query:"limit"`
	Offset int    `json:"offset" query:"offset"`
}

type GetSubscriptionsRequest struct {
	UserID      *string `json:"user_id" query:"user_id"`
	ServiceName *string `json:"service_name" query:"service_name"`
	StartDate   *string `json:"start_date" query:"start_date"`
	EndDate     *string `json:"end_date" query:"end_date"`
	Limit       int     `json:"limit" query:"limit"`
	Offset      int     `json:"offset" query:"offset"`
}

type CalculateCostRequest struct {
	UserID      *string `json:"user_id" query:"user_id"`
	ServiceName *string `json:"service_name" query:"service_name"`
	StartDate   string  `json:"start_date" query:"start_date"`
	EndDate     string  `json:"end_date" query:"end_date"`
}

func (r *CreateSubscriptionRequest) GetUserID() (uuid.UUID, error) {
	return uuid.Parse(r.UserID)
}

func (r *GetSubscriptionRequest) GetID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

func (r *DeleteSubscriptionRequest) GetID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

func (r *GetUserSubscriptionsRequest) GetUserID() (uuid.UUID, error) {
	return uuid.Parse(r.UserID)
}

func (r *GetSubscriptionsRequest) GetUserID() (*uuid.UUID, error) {
	if r.UserID == nil || *r.UserID == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*r.UserID)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
