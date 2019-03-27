package model

import (
	"encoding/json"
	"time"
)

// Payment - representation of a Form3 payment
type Payment struct {
	// TODO ID should be stored as binary for better performance
	ID                string `json:"id,omitempty" bson:"id,omitempty"`
	OrganisationID    string `json:"organisation_id" bson:"organisation_id"`
	Type              string `json:"type" bson:"type"`
	PaymentAttributes `json:"attributes" bson:"attributes"`
}

// PaymentAttributes - representation for Form3 payment attributes
type PaymentAttributes struct {
	Amount               string            `json:"amount" bson:"amount"`
	BeneficiaryParty     PaymentParty      `json:"beneficiary_party" bson:"beneficiary_party"`
	ChargesInfo          PaymentChargeInfo `json:"charges_information" bson:"charges_information"`
	Currency             string            `json:"currency" bson:"currency"`
	DebtorParty          PaymentParty      `json:"debtor_party" bson:"debtor_party"`
	EndToEndReference    string            `json:"end_to_end_reference" bson:"end_to_end_reference"`
	FX                   ForeignExchange   `json:"fx" bson:"fx"`
	NumericReference     string            `json:"numeric_reference" bson:"numeric_reference"`
	ID                   string            `json:"payment_id" bson:"payment_id"`
	Purpose              string            `json:"payment_purpose" bson:"payment_purpose"`
	Scheme               string            `json:"payment_scheme" bson:"payment_scheme"`
	Type                 string            `json:"payment_type" bson:"payment_type"`
	ProcessingDate       Time              `json:"processing_date" bson:"processing_date"`
	Reference            string            `json:"reference" bson:"reference"`
	SchemePaymentSubType string            `json:"scheme_payment_sub_type" bson:"scheme_payment_sub_type"`
	SchemePaymentType    string            `json:"scheme_payment_type" bson:"scheme_payment_type"`
	SponsorParty         PaymentParty      `json:"sponsor_party" bson:"sponsor_party"`
}

// PaymentParty - representation for Form3 payment parties
type PaymentParty struct {
	AccountName       string `json:"account_name,omitempty" bson:"account_name,omitempty"`
	AccountNumber     string `json:"account_number" bson:"account_number"`
	AccountNumberCode string `json:"account_number_code,omitempty" bson:"account_number_code,omitempty"`
	AccountType       int    `json:"account_type,omitempty" bson:"account_type,omitempty"`
	Address           string `json:"address,omitempty" bson:"address,omitempty"`
	BankID            string `json:"bank_id" bson:"bank_id"`
	BankIDCode        string `json:"bank_id_code" bson:"bank_id_code"`
	Name              string `json:"name,omitempty" bson:"name,omitempty"`
}

// PaymentChargeInfo - representation for Form3 payment charge information
type PaymentChargeInfo struct {
	BearerCode              string          `json:"bearer_code" bson:"bearer_code"`
	SenderCharges           []PaymentCharge `json:"sender_charges" bson:"sender_charges"`
	ReceiverChargesAmount   string          `json:"receiver_charges_amount" bson:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `json:"receiver_charges_currency" bson:"receiver_charges_currency"`
}

// PaymentCharge - representation for Form3 payment charge
type PaymentCharge struct {
	Amount   string `json:"amount" bson:"amount"`
	Currency string `json:"currency" bson:"currency"`
}

// ForeignExchange - representation for Form3 FX
type ForeignExchange struct {
	ContractReference string `json:"contract_reference" bson:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate" bson:"exchange_rate"`
	OriginalAmount    string `json:"original_amount" bson:"original_amount"`
	OriginalCurrency  string `json:"original_currency" bson:"original_currency"`
}

// TimeLayout - time layout for Form3 payments
const TimeLayout = "2006-01-02"

// Time - time wrapper
type Time struct {
	time.Time
}

// MarshalJSON - performs the marshalling for the custom Time used within a Payment
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(TimeLayout))
}

// UnmarshalJSON - performs the unmarshalling for the custom Time used within a Payment
func (t *Time) UnmarshalJSON(data []byte) error {
	input := string(data)
	parsedTime, err := time.Parse(`"`+TimeLayout+`"`, input)
	if err != nil {
		return err
	}

	*t = Time{parsedTime}
	return nil
}

// type BsonUUID struct {
// 	primitive.Binary
// }

// func NewBsonUUID() BsonUUID {
// 	uuid := uuid.New()
// 	return BsonUUID{primitive.Binary{
// 		Subtype: 0x04,
// 		Data:    uuid[:],
// 	}}
// }

// func (b BsonUUID) MarshalJSON() ([]byte, error) {
// 	if len(b.Data) == 0 {
// 		return []byte(`""`), nil
// 	}

// 	u, err := uuid.FromBytes(b.Data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return json.Marshal(u.String())
// }

// func (b *BsonUUID) UnmarshalJSON(data []byte) error {
// 	if string(data) == "null" {
// 		return nil
// 	}
// 	u, err := uuid.Parse(string(data))
// 	if err != nil {
// 		return fmt.Errorf("invalid uuid %s", data)
// 	}

// 	*b = BsonUUID{primitive.Binary{
// 		Subtype: 0x04,
// 		Data:    u[:],
// 	}}
// 	return nil
// }
