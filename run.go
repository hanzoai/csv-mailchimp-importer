package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/zeekay/gochimp3"
)

func FileExists(dir string) bool {
	if _, err := os.Stat(dir); err == nil {
		return true
	}
	return false
}

func MkDir(dir string) {
	if !FileExists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func Prepare(o *Order, c *Customer) {
	o.Customer = c
	c.Company = o.Company
	c.OptInStatus = true
	DecodeBytes([]byte(o.ItemsJSON), &o.Items)
}

func main() {
	dataPath := Config.DataPath

	filename := dataPath + "/order.csv"
	ofh, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer ofh.Close()

	filename = dataPath + "/user.csv"
	ufh, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer ufh.Close()

	var ords []*Order
	var usrs []*Customer

	err = gocsv.UnmarshalFile(ofh, &ords)
	if err != nil {
		panic(err)
	}

	err = gocsv.UnmarshalFile(ufh, &usrs)
	if err != nil {
		panic(err)
	}

	users := make(map[string]*Customer)

	for _, usr := range usrs {
		users[usr.ID] = usr
	}

	client := gochimp3.New(Config.APIKey)
	client.Transport = &http.Transport{}
	client.Debug = true

	// resumed := false

	stor, err := client.GetStore(Config.DefaultStore, nil)
	if err != nil {
		panic(err)
	}

	list, err := client.GetList(Config.ListId, nil)
	if err != nil {
		panic(err)
	}

	for _, ord := range ords {
		// if ord.ID == "6DimAjn6fDHEGOqdIJRoJA" {
		// 	resumed = true
		// }
		// if !resumed {
		// 	continue
		// }

		usr, ok := users[ord.UserId]

		if ok {
			Prepare(ord, usr)
			CreateOrder(client, ord, stor)
			SubscribeCustomer(client, usr, list)
		} else {
			fmt.Printf("ORPHANNED %v\n", Encode(ord))
		}

		// return
	}
}

func centsToFloat(v int64, currency string) float64 {
	return float64(v) / 100.00
}

var timeFormat = "2006-01-02T15:04:05"

func SubscribeCustomer(client *gochimp3.API, usr *Customer, list *gochimp3.ListResponse) {
	status := "subscribed"

	req := &gochimp3.MemberRequest{
		EmailAddress: usr.EmailAddress,
		Status:       status,
		// MergeFields:  s.MergeFields(),
		Interests: make(map[string]interface{}),
		// Language:  s.Client.Language,
		VIP: false,
		Location: &gochimp3.MemberLocation{
			Latitude:  0.0,
			Longitude: 0.0,
			GMTOffset: 0,
			DSTOffset: 0,
			// CountryCode: s.Client.Country,
			Timezone: "",
		},
	}

	h := md5.New()
	io.WriteString(h, usr.EmailAddress)
	md5 := fmt.Sprintf("%x", h.Sum(nil))

	// Try to update subscriber, create new member if that fails.
	if _, err := list.UpdateMember(md5, req); err != nil {
		_, err = list.CreateMember(req)
	}

	// fmt.Printf("ADDING %v\n", Encode(req))
}

func CreateOrder(client *gochimp3.API, ord *Order, stor *gochimp3.Store) {
	usr := ord.Customer

	// Fetch user
	lines := make([]gochimp3.LineItem, 0)
	for _, line := range ord.Items {
		lines = append(lines, gochimp3.LineItem{
			ID:               ord.ID + line.Id(),
			ProductID:        line.ProductId,
			ProductVariantID: line.Id(),
			Quantity:         line.Quantity,
			Price:            centsToFloat(line.Price, ord.CurrencyCode),
		})
	}

	processedAtForeign, _ := time.Parse(timeFormat, ord.ProcessedAtForeign)
	cancelledAtForeign, _ := time.Parse(timeFormat, ord.CancelledAtForeign)
	updatedAtForeign, _ := time.Parse(timeFormat, ord.UpdatedAtForeign)

	// Create Order
	req := &gochimp3.Order{
		// Required
		ID:           ord.ID,
		CurrencyCode: strings.ToUpper(ord.CurrencyCode),
		OrderTotal:   centsToFloat(ord.OrderTotal, ord.CurrencyCode),
		Customer: gochimp3.Customer{
			// Required
			ID: usr.ID,

			// Optional
			EmailAddress: usr.EmailAddress,
			OptInStatus:  true,
			Company:      ord.Company,
			FirstName:    usr.FirstName,
			LastName:     usr.LastName,
			// OrdersCount:  1,
			// TotalSpent:   centsToFloat(usr.Total, usr.Currency),
			Address: &gochimp3.Address{
				Address1:     ord.ShippingAddressLine1,
				Address2:     ord.ShippingAddressLine2,
				City:         ord.ShippingAddressCity,
				ProvinceCode: ord.ShippingAddressState,
				PostalCode:   ord.ShippingAddressPostalCode,
				CountryCode:  ord.ShippingAddressCountryCode,
			},
		},
		Lines: lines,

		// Optional
		TaxTotal:          centsToFloat(ord.TaxTotal, ord.CurrencyCode),
		ShippingTotal:     centsToFloat(ord.ShippingTotal, ord.CurrencyCode),
		FinancialStatus:   string(ord.FinancialStatus),
		FulfillmentStatus: string(ord.FulfillmentStatus),
		CampaignID:        ord.MailchimpCampaignID,
		TrackingCode:      ord.MailchimpTrackingCode,
		BillingAddress: &gochimp3.Address{
			Address1:     ord.BillingAddressLine1,
			Address2:     ord.BillingAddressLine2,
			City:         ord.BillingAddressCity,
			ProvinceCode: ord.BillingAddressState,
			PostalCode:   ord.BillingAddressPostalCode,
			CountryCode:  ord.BillingAddressCountryCode,
		},
		ShippingAddress: &gochimp3.Address{
			Address1:     ord.ShippingAddressLine1,
			Address2:     ord.ShippingAddressLine2,
			City:         ord.ShippingAddressCity,
			ProvinceCode: ord.ShippingAddressState,
			PostalCode:   ord.ShippingAddressPostalCode,
			CountryCode:  ord.ShippingAddressCountryCode,
		},
		ProcessedAtForeign: processedAtForeign,
		CancelledAtForeign: cancelledAtForeign,
		UpdatedAtForeign:   updatedAtForeign,
	}

	_, err := stor.UpdateOrder(req)
	if err != nil {
		_, err = stor.CreateOrder(req)
	}

	// fmt.Printf("ADDING %v, %v\n", err, Encode(req))
}
