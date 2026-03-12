package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	square "github.com/square/square-go-sdk"
	squareclient "github.com/square/square-go-sdk/client"
	"github.com/square/square-go-sdk/checkout"
	"github.com/square/square-go-sdk/option"
	"github.com/retrosnack-clothing/retrosnack/internal/orders"
)

type Service interface {
	CreateCheckout(ctx context.Context, req CreateCheckoutRequest, redirectURL string) (*CheckoutSession, error)
	ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*PaymentResult, error)
	HandleWebhook(ctx context.Context, payload []byte, signatureHeader string) error
}

type service struct {
	orders          orders.Service
	square          *squareclient.Client
	locationID      string
	webhookSigKey   string
	webhookNotifURL string
}

func NewService(ordersSvc orders.Service, accessToken, locationID, webhookSigKey, webhookNotifURL, squareEnv string) Service {
	opts := []option.RequestOption{
		option.WithToken(accessToken),
	}
	if squareEnv == "sandbox" {
		opts = append(opts, option.WithBaseURL(square.Environments.Sandbox))
	}
	c := squareclient.NewClient(opts...)
	return &service{
		orders:          ordersSvc,
		square:          c,
		locationID:      locationID,
		webhookSigKey:   webhookSigKey,
		webhookNotifURL: webhookNotifURL,
	}
}

func (s *service) CreateCheckout(ctx context.Context, req CreateCheckoutRequest, redirectURL string) (*CheckoutSession, error) {
	order, err := s.orders.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	lineItems := make([]*square.OrderLineItem, 0, len(order.Items))
	for _, item := range order.Items {
		lineItems = append(lineItems, &square.OrderLineItem{
			Name:     square.String(item.VariantID.String()),
			Quantity: strconv.Itoa(item.Quantity),
			BasePriceMoney: &square.Money{
				Amount:   square.Int64(item.PriceCents),
				Currency: square.CurrencyCad.Ptr(),
			},
			Note: square.String("order:" + order.ID.String()),
		})
	}

	idempotencyKey := order.ID.String()

	sqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.square.Checkout.PaymentLinks.Create(sqCtx, &checkout.CreatePaymentLinkRequest{
		IdempotencyKey: &idempotencyKey,
		Order: &square.Order{
			LocationID:  s.locationID,
			ReferenceID: square.String(order.ID.String()),
			LineItems:   lineItems,
		},
		CheckoutOptions: &square.CheckoutOptions{
			RedirectURL: &redirectURL,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	link := resp.PaymentLink
	if link == nil || link.ID == nil {
		return nil, fmt.Errorf("square returned empty payment link")
	}

	if err := s.orders.SetCheckoutSession(ctx, order.ID, *link.ID); err != nil {
		return nil, err
	}

	url := ""
	if link.URL != nil {
		url = *link.URL
	}

	return &CheckoutSession{
		ID:      *link.ID,
		OrderID: order.ID,
		URL:     url,
	}, nil
}

func (s *service) ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*PaymentResult, error) {
	order, err := s.orders.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	idempotencyKey := uuid.New().String()

	payCtx, payCancel := context.WithTimeout(ctx, 10*time.Second)
	defer payCancel()

	resp, err := s.square.Payments.Create(payCtx, &square.CreatePaymentRequest{
		SourceID:       req.SourceID,
		IdempotencyKey: idempotencyKey,
		AmountMoney: &square.Money{
			Amount:   square.Int64(order.TotalCents),
			Currency: square.CurrencyCad.Ptr(),
		},
		LocationID:  &s.locationID,
		ReferenceID: square.String(order.ID.String()),
	})
	if err != nil {
		slog.Error("square payment failed", "order_id", order.ID, "error", err)
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	if resp.Payment == nil || resp.Payment.ID == nil {
		slog.Error("square returned empty payment", "order_id", order.ID)
		return nil, fmt.Errorf("square returned empty payment")
	}

	status := "pending"
	if resp.Payment.Status != nil && *resp.Payment.Status == "COMPLETED" {
		if err := s.orders.MarkPaid(ctx, order.ID); err != nil {
			return nil, fmt.Errorf("failed to mark order paid: %w", err)
		}
		status = "paid"
	}

	slog.Info("payment processed", "order_id", order.ID, "payment_id", *resp.Payment.ID, "status", status, "amount_cents", order.TotalCents)

	return &PaymentResult{
		OrderID:   order.ID,
		PaymentID: *resp.Payment.ID,
		Status:    status,
	}, nil
}

func (s *service) HandleWebhook(ctx context.Context, payload []byte, signatureHeader string) error {
	err := s.square.Webhooks.VerifySignature(ctx, &square.VerifySignatureRequest{
		RequestBody:     string(payload),
		SignatureHeader: signatureHeader,
		SignatureKey:    s.webhookSigKey,
		NotificationURL: s.webhookNotifURL,
	})
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	var event struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse webhook event: %w", err)
	}

	slog.Info("webhook received", "type", event.Type)

	if event.Type != "payment.updated" {
		return nil
	}

	var data struct {
		Object struct {
			Payment struct {
				Status  string `json:"status"`
				OrderID string `json:"order_id"`
			} `json:"payment"`
		} `json:"object"`
	}
	if err := json.Unmarshal(event.Data, &data); err != nil {
		return fmt.Errorf("failed to parse payment data: %w", err)
	}

	if data.Object.Payment.Status != "COMPLETED" {
		return nil
	}

	squareOrderID := data.Object.Payment.OrderID
	if squareOrderID == "" {
		return nil
	}

	// look up the square order to get our reference_id (our order UUID)
	whCtx, whCancel := context.WithTimeout(ctx, 10*time.Second)
	defer whCancel()

	squareOrder, err := s.square.Orders.Get(whCtx, &square.GetOrdersRequest{OrderID: squareOrderID})
	if err != nil {
		return fmt.Errorf("failed to get square order: %w", err)
	}

	if squareOrder.Order == nil || squareOrder.Order.ReferenceID == nil {
		return fmt.Errorf("square order missing reference_id")
	}

	orderID, err := uuid.Parse(*squareOrder.Order.ReferenceID)
	if err != nil {
		return fmt.Errorf("invalid reference_id in square order: %w", err)
	}

	slog.Info("webhook marking order paid", "order_id", orderID)
	return s.orders.MarkPaid(ctx, orderID)
}
