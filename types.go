package main

type Item struct {
	ProductId string
	VariantId string
	Quantity  int
	Price     int64
}

func (i Item) Id() string {
	if i.VariantId != "" {
		return i.VariantId
	}

	return i.ProductId
}

type Customer struct {
	ID string `csv:"Id_"`

	EmailAddress string `csv:"Email"`
	FirstName    string `csv:"FirstName"`
	LastName     string `csv:"LastName"`

	// Misc
	OptInStatus bool
	Company     string
}

type Order struct {
	ID                    string `csv:"Id_"`
	UserId                string `csv:"UserId"`
	Company               string `csv:"Company"`
	CurrencyCode          string `csv:"Currency"`
	OrderTotal            int64  `csv:"Total"`
	TaxTotal              int64  `csv:"Tax"`
	ShippingTotal         int64  `csv:"Shipping"`
	FinancialStatus       string `csv:"PaymentStatus"`
	FulfillmentStatus     string `csv:"FulfillmentStatus"`
	MailchimpCampaignID   string `csv:"Mailchimp.CampaignId"`
	MailchimpTrackingCode string `csv:"Mailchimp.TrackingCode"`

	ProcessedAtForeign string `csv:"CreatedAt"`
	CancelledAtForeign string `csv:"CancelledAt"`
	UpdatedAtForeign   string `csv:"UpdatedAt"`

	ItemsJSON string `csv:"Items_"`
	Items     []Item

	ShippingAddressLine1       string `csv:"ShippingAddress.Line1"`
	ShippingAddressLine2       string `csv:"ShippingAddress.Line2"`
	ShippingAddressCity        string `csv:"ShippingAddress.City"`
	ShippingAddressState       string `csv:"ShippingAddress.State"`
	ShippingAddressPostalCode  string `csv:"ShippingAddress.PostalCode"`
	ShippingAddressCountryCode string `csv:"ShippingAddress.CountryCode"`

	BillingAddressLine1       string `csv:"BillingAddress.Line1"`
	BillingAddressLine2       string `csv:"BillingAddress.Line2"`
	BillingAddressCity        string `csv:"BillingAddress.City"`
	BillingAddressState       string `csv:"BillingAddress.State"`
	BillingAddressPostalCode  string `csv:"BillingAddress.PostalCode"`
	BillingAddressCountryCode string `csv:"BillingAddress.CountryCode"`

	// Misc
	Customer *Customer
}
