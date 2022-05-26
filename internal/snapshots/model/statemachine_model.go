package model

import (
	"time"

	"github.com/99minutos/shipments-snapshot-service/pkg/pb"
)

type Shipment struct {
	ID           string  `json:"id" bson:"_id"`
	TrackingID   string  `json:"tracking_id" bson:"trackingId"`
	InternalKey  string  `json:"internal_key" bson:"internalKey"`
	OrderID      string  `json:"order_id" bson:"orderId"`
	Status       string  `json:"status" bson:"status"`
	DeliveryType string  `json:"delivery_type" bson:"deliveryType"`
	Origin       Address `json:"origin" bson:"origin"`
	Destination  Address `json:"destination" bson:"destination"`
	Recipient    Person  `json:"recipient" bson:"recipient"`
	Sender       Person  `json:"sender" bson:"sender"`
	Payment      Payment `json:"payment" bson:"payment"`
	Option       Option  `json:"option" bson:"option"`
	Items        []Item  `json:"items" bson:"items"`
}

type Address struct {
	StreetLine string  `json:"street_line" bson:"address"`
	Lat        float64 `json:"lat" bson:"lat"`
	Lng        float64 `json:"lng" bson:"lng"`
	Type       string  `json:"type" bson:"type"`
}

type Person struct {
	FirstName   string `json:"first_name" bson:"firstName"`
	LastName    string `json:"last_name" bson:"LastName"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phone_number" bson:"phone"`
}

type Payment struct {
	PaymentMethod string `json:"payment_method" bson:"paymentMethod"`
}

type Option struct {
	PickUpAfter            *time.Time      `json:"pick_up_after,omitempty" bson:"pickUpAfter"`
	RequiresIdentification *bool           `json:"requires_identification,omitempty" bson:"requiresIdentification"`
	RequiresSignature      *bool           `json:"requires_signature,omitempty" bson:"requiresSignature"`
	TwoFactorAuth          *bool           `json:"two_factor_auth,omitempty" bson:"twoFactorAuth"`
	LegacyMetadata         *LegacyMetadata `json:"legacy_metadata,omitempty" bson:"metadata"`
}

type LegacyMetadata struct {
	AmountCash   *string `json:"amount_cash,omitempty" bson:"AmountCash"`
	LegacyApikey *string `json:"legacy_apikey,omitempty" bson:"LegacyApikey"`
	Notes        *string `json:"notes,omitempty" bson:"Notes"`
}

type Item struct {
	Size   string  `json:"size" bson:"size"`
	Weight float32 `json:"weight" bson:"weight"`
}

type Event struct {
	EventName  string        `json:"event_name" bson:"eventName"`
	StatusCode string        `json:"status_code" bson:"statusCode"`
	StatusName string        `json:"status_name" bson:"statusName"`
	Data       EventData     `json:"data,omitempty" bson:"data"`
	Metadata   EventMetadata `json:"metadata,omitempty" bson:"metadata"`
}

type EventData struct {
	Station *string `json:"station,omitempty" bson:"station"`
}

type EventMetadata struct {
	Platform *string `json:"platform,omitempty" bson:"platform"`
}

func TransformShipmentToProto(shipment *Shipment) (*pb.Shipment, error) {
	origin := &pb.Address{
		StreetLine: shipment.Origin.StreetLine,
		Lat:        shipment.Origin.Lat,
		Lng:        shipment.Origin.Lng,
		Type:       shipment.Origin.Type,
	}

	destination := &pb.Address{
		StreetLine: shipment.Destination.StreetLine,
		Lat:        shipment.Destination.Lat,
		Lng:        shipment.Destination.Lng,
		Type:       shipment.Destination.Type,
	}

	recipient := &pb.Person{
		FirstName:   shipment.Recipient.FirstName,
		LastName:    shipment.Recipient.LastName,
		Email:       shipment.Recipient.Email,
		PhoneNumber: shipment.Recipient.PhoneNumber,
	}

	sender := &pb.Person{
		FirstName:   shipment.Sender.FirstName,
		LastName:    shipment.Sender.LastName,
		Email:       shipment.Sender.Email,
		PhoneNumber: shipment.Sender.PhoneNumber,
	}

	payment := &pb.Payment{
		PaymentMethod: shipment.Payment.PaymentMethod,
	}

	var legacyMetadata *pb.LegacyMetadata
	if shipment.Option.LegacyMetadata != nil {
		legacyMetadata = &pb.LegacyMetadata{
			AmountCash:   shipment.Option.LegacyMetadata.AmountCash,
			LegacyApikey: shipment.Option.LegacyMetadata.LegacyApikey,
			Notes:        shipment.Option.LegacyMetadata.Notes,
		}
	} else {
		legacyMetadata = nil
	}

	pickupAfter := ""
	if shipment.Option.PickUpAfter != nil {
		pickupAfter = shipment.Option.PickUpAfter.String()
	}

	option := &pb.Option{
		PickUpAfter:            &pickupAfter,
		RequiresIdentification: shipment.Option.RequiresIdentification,
		RequiresSignature:      shipment.Option.RequiresSignature,
		TwoFactoraAuth:         shipment.Option.TwoFactorAuth,
		LegacyMetadata:         legacyMetadata,
	}

	items := make([]*pb.Item, 0)
	for _, item := range shipment.Items {
		items = append(items, &pb.Item{
			Size:   item.Size,
			Weight: item.Weight,
		})
	}

	return &pb.Shipment{
		TrackingID:   shipment.TrackingID,
		InternalKey:  shipment.InternalKey,
		OrderID:      shipment.OrderID,
		Status:       shipment.Status,
		Origin:       origin,
		Destination:  destination,
		Recipient:    recipient,
		Sender:       sender,
		DeliveryType: shipment.DeliveryType,
		Payment:      payment,
		Option:       option,
		Items:        items,
	}, nil
}

func TransformEventToProto(event *Event) (*pb.Event, error) {
	legacy := &pb.Legacy{
		Station: event.Data.Station,
	}

	data := &pb.EventData{
		Legacy: legacy,
	}

	metadata := &pb.EventMetadata{
		Platform: event.Metadata.Platform,
	}

	return &pb.Event{
		EventName:  event.EventName,
		StatusCode: event.StatusCode,
		StatusName: event.StatusName,
		Data:       data,
		Metadata:   metadata,
	}, nil
}

func (s *Shipment) IsV2() bool {
	return false
}

func (s *Shipment) IsV3() bool {
	return true
}
