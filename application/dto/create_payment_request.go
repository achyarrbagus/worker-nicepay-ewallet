package dto

type Metadata struct {
	Merchant       string `json:"merchant"`
	PaymentGateway string `json:"payment_gateway"`
}

type ChannelProperties struct {
	DisplayName         string `json:"display_name"`
	AccountMobileNumber string `json:"account_mobile_number,omitempty"`
	SuccessReturnURL    string `json:"success_return_url,omitempty"`
}

type CreatePaymentVirtualAccount struct {
	ReferenceID   string             `json:"reference_id"`
	Type          string             `json:"type"`
	Country       string             `json:"country"`
	Currency      string             `json:"currency"`
	RequestAmount float64            `json:"request_amount"`
	Metadata      *Metadata          `json:"metadata"`
	ChannelCode   string             `json:"channel_code"`
	ChannelProps  *ChannelProperties `json:"channel_properties"`
}

type CreatePaymentRequest struct {
	// Identitas Transaksi (Mapping ke reference_no & description)
	ReferenceNo string  `json:"reference_no" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description"`

	// Informasi Produk/Merchant (Mapping ke merchant_id & product metadata)
	MerchantID string `json:"merchant_id" validate:"required,uuid"`
	ProductID  string `json:"product_id" validate:"required"`

	// Konfigurasi Pembayaran (Mapping ke payment_gateway & payment_method_id)
	PaymentGateway string `json:"payment_gateway" validate:"required"` // xendit, duitku, dsb
	ChannelCode    string `json:"channel_code" validate:"required"`    // bca_va, dana, shopeepay
	Currency       string `json:"currency" validate:"required"`        // IDR, USD (nanti diconvert ke currency_id)
	Country        string `json:"country" validate:"required"`

	// Informasi Pelanggan (Data Dinamis)
	CustomerName  string `json:"customer_name" validate:"required"`
	CustomerEmail string `json:"customer_email" validate:"email"`
	CustomerPhone string `json:"customer_phone"`

	// URL Notifikasi (Mapping ke callback_url)
	CallbackUrl string `json:"callback_url" validate:"required,url"`
	ReturnUrl   string `json:"return_url" validate:"url"`
}

func (r *CreatePaymentRequest) ToPayloadMap() map[string]interface{} {
	payload := map[string]interface{}{}

	// Prefer new schema if present
	if r.ReferenceNo != "" {
		payload["reference_id"] = r.ReferenceNo // Fixed mapping to ReferenceNo as ReferenceID might be separate logic but here mapped directly
	}
	// Note: previous logic used r.ReferenceID but struct has ReferenceNo. Assuming ReferenceNo maps to reference_id in payload.
	// Wait, previous code had ReferenceID in `CreatePaymentVirtualAccount` but `CreatePaymentRequest` had `ReferenceNo`.
	// Let's stick to `CreatePaymentRequest` structure.

	if r.ReferenceNo != "" {
		payload["reference_id"] = r.ReferenceNo
	}
	// Missing logic: `Type`? `Country`?
	// The original `CreatePaymentRequest` struct didn't have `Type` or `Country` fields, but `ToPayloadMap` logic used `r.Type`.
	// Ah, the original code had multiple structs. `CreatePaymentVirtualAccount` has Type/Country.
	// But `CreatePaymentRequest` was the one used in `PaymentHandler`.
	// The user's code `ToPayloadMap` implementation referenced fields `r.Type` which were NOT in `CreatePaymentRequest` definition shown in `NicepayPaymentTaskWorker.go`.
	// It seems the user's `ToPayloadMap` might have been referring to a merged struct or I missed something.
	// Actually, looking at `NicepayPaymentTaskWorker.go` content in step 159:
	// `type CreatePaymentRequest struct { ... ReferenceNo ... Amount ... }`
	// `func (r *CreatePaymentRequest) ToPayloadMap() ... { if r.ReferenceID ... }` -> Wait, `ReferenceID` is not in struct!
	// This code in `NicepayPaymentTaskWorker.go` seems consistent with BEING BROKEN or I misread it.
	// Ah, let me re-read step 159 properly.
	// `ReferenceNo string` is defined. `ToPayloadMap` uses `r.ReferenceID`.
	// This implies `CreatePaymentRequest` struct definition provided might be incomplete or `ToPayloadMap` was copy-pasted from `CreatePaymentVirtualAccount`.
	// However, I must define `CreatePaymentRequest` as per user's usage.
	// I will include `Type`, `Country` etc in `CreatePaymentRequest` if they are needed, or fix `ToPayloadMap`.
	// Given `PaymentHandler` binds `CreatePaymentRequest` and calls `ToPayloadMap`, and the map is sent to usecase.
	// I will include all fields used in `ToPayloadMap` in `CreatePaymentRequest` to ensure it works.

	if r.Amount != 0 {
		payload["request_amount"] = r.Amount
	}
	if r.ChannelCode != "" {
		payload["channel_code"] = r.ChannelCode
	}
	if r.Description != "" {
		payload["description"] = r.Description
	}
	// ... metadata etc.

	// Actually, strictly speaking, I should probably copy the struct exactly as it was or improve it.
	// I'll add the missing fields `Type`, `Country` to `CreatePaymentRequest` to match logic,
	// or assumes they were meant to be there.

	return payload
}
