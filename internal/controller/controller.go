package controller

import (
	"context"
	"net/http"
)

type PingController interface {
	GetPing(ctx context.Context) http.HandlerFunc
}

type TenderController interface {
	GetTenders(ctx context.Context) http.HandlerFunc
	PostNewTender(ctx context.Context) http.HandlerFunc
	GetUserTenders(ctx context.Context) http.HandlerFunc
	GetTenderStatus(ctx context.Context) http.HandlerFunc
	PutTenderStatus(ctx context.Context) http.HandlerFunc
	PatchTender(ctx context.Context) http.HandlerFunc
	PutTenderRollback(ctx context.Context) http.HandlerFunc
}

type BidController interface {
	PostNewBid(ctx context.Context) http.HandlerFunc
	GetUserBids(ctx context.Context) http.HandlerFunc
	GetTenderBids(ctx context.Context) http.HandlerFunc
	GetBidStatus(ctx context.Context) http.HandlerFunc
	PutBidStatus(ctx context.Context) http.HandlerFunc
	PatchBid(ctx context.Context) http.HandlerFunc
	PutBidSubmitDecision(ctx context.Context) http.HandlerFunc
	PutBidFeedback(ctx context.Context) http.HandlerFunc
	PutBidRollback(ctx context.Context) http.HandlerFunc
	GetBidReviews(ctx context.Context) http.HandlerFunc
}
