package payment

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"

	"github.com/tenzi-muraru/paymentApp/payment/model"
)

const paymentJSON1 = `{  
	"id":"7eb8277a-6c91-45e9-8a03-a27f82aca350",
	"organisation_id":"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
	"type":"Payment",
	"attributes":{  
	   "amount":"100.21",
	   "beneficiary_party":{  
		  "account_name":"W Owens",
		  "account_number":"31926819",
		  "account_number_code":"BBAN",
		  "address":"1 The Beneficiary Localtown SE2",
		  "bank_id":"403000",
		  "bank_id_code":"GBDSC",
		  "name":"Wilfred Jeremiah Owens"
	   },
	   "charges_information":{  
		  "bearer_code":"SHAR",
		  "sender_charges":[  
			 {  
				"amount":"5.00",
				"currency":"GBP"
			 },
			 {  
				"amount":"10.00",
				"currency":"USD"
			 }
		  ],
		  "receiver_charges_amount":"1.00",
		  "receiver_charges_currency":"USD"
	   },
	   "currency":"GBP",
	   "debtor_party":{  
		  "account_name":"EJ Brown Black",
		  "account_number":"GB29XABC10161234567801",
		  "account_number_code":"IBAN",
		  "address":"10 Debtor Crescent Sourcetown NE1",
		  "bank_id":"203301",
		  "bank_id_code":"GBDSC",
		  "name":"Emelia Jane Brown"
	   },
	   "end_to_end_reference":"Wil piano Jan",
	   "fx":{  
		  "contract_reference":"FX123",
		  "exchange_rate":"2.00000",
		  "original_amount":"200.42",
		  "original_currency":"USD"
	   },
	   "numeric_reference":"1002001",
	   "payment_id":"123456789012345678",
	   "payment_purpose":"Paying for goods/services",
	   "payment_scheme":"FPS",
	   "payment_type":"Credit",
	   "processing_date":"2017-01-18",
	   "reference":"Payment for Em's piano lessons",
	   "scheme_payment_sub_type":"InternetBanking",
	   "scheme_payment_type":"ImmediatePayment",
	   "sponsor_party":{  
		  "account_number":"56781234",
		  "bank_id":"123123",
		  "bank_id_code":"GBDSC"
	   }
	}
 }`

const paymentJSON2 = `{  
	"type":"Payment",
	"id":"216d4da9-e59a-4cc6-8df3-3da6e7580b77",
	"organisation_id":"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
	"attributes":{  
	   "amount":"100.21",
	   "beneficiary_party":{  
		  "account_name":"W Owens",
		  "account_number":"31926819",
		  "account_number_code":"BBAN",
		  "account_type":0,
		  "address":"1 The Beneficiary Localtown SE2",
		  "bank_id":"403000",
		  "bank_id_code":"GBDSC",
		  "name":"Wilfred Jeremiah Owens"
	   },
	   "charges_information":{  
		  "bearer_code":"SHAR",
		  "sender_charges":[  
			 {  
				"amount":"5.00",
				"currency":"GBP"
			 },
			 {  
				"amount":"10.00",
				"currency":"USD"
			 }
		  ],
		  "receiver_charges_amount":"1.00",
		  "receiver_charges_currency":"USD"
	   },
	   "currency":"GBP",
	   "debtor_party":{  
		  "account_name":"EJ Brown Black",
		  "account_number":"GB29XABC10161234567801",
		  "account_number_code":"IBAN",
		  "address":"10 Debtor Crescent Sourcetown NE1",
		  "bank_id":"203301",
		  "bank_id_code":"GBDSC",
		  "name":"Emelia Jane Brown"
	   },
	   "end_to_end_reference":"Wil piano Jan",
	   "fx":{  
		  "contract_reference":"FX123",
		  "exchange_rate":"2.00000",
		  "original_amount":"200.42",
		  "original_currency":"USD"
	   },
	   "numeric_reference":"1002001",
	   "payment_id":"123456789012345678",
	   "payment_purpose":"Paying for goods/services",
	   "payment_scheme":"FPS",
	   "payment_type":"Credit",
	   "processing_date":"2017-01-18",
	   "reference":"Payment for Em's piano lessons",
	   "scheme_payment_sub_type":"InternetBanking",
	   "scheme_payment_type":"ImmediatePayment",
	   "sponsor_party":{  
		  "account_number":"56781234",
		  "bank_id":"123123",
		  "bank_id_code":"GBDSC"
	   }
	}
 }`

const paymentJSON3 = `{  
	"type":"Payment",
	"organisation_id":"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
	"attributes":{  
	   "amount":"100.21",
	   "beneficiary_party":{  
		  "account_name":"J Doe",
		  "account_number":"11223344",
		  "account_number_code":"BBAN",
		  "account_type":0,
		  "address":"1 The Beneficiary Localtown SE2",
		  "bank_id":"403000",
		  "bank_id_code":"GBDSC",
		  "name":"John Doe"
	   },
	   "charges_information":{  
		  "bearer_code":"SHAR",
		  "sender_charges":[  
			 {  
				"amount":"5.00",
				"currency":"GBP"
			 },
			 {  
				"amount":"10.00",
				"currency":"USD"
			 }
		  ],
		  "receiver_charges_amount":"1.00",
		  "receiver_charges_currency":"USD"
	   },
	   "currency":"GBP",
	   "debtor_party":{  
		  "account_name":"Ja Doe",
		  "account_number":"GB29XABC10161234567801",
		  "account_number_code":"IBAN",
		  "address":"10 Debtor Crescent Sourcetown NE1",
		  "bank_id":"203301",
		  "bank_id_code":"GBDSC",
		  "name":"Jane Doe"
	   },
	   "end_to_end_reference":"Wil piano Jan",
	   "fx":{  
		  "contract_reference":"FX123",
		  "exchange_rate":"2.00000",
		  "original_amount":"200.42",
		  "original_currency":"USD"
	   },
	   "numeric_reference":"1002001",
	   "payment_id":"123456789012345678",
	   "payment_purpose":"Paying for goods/services",
	   "payment_scheme":"FPS",
	   "payment_type":"Credit",
	   "processing_date":"2017-01-18",
	   "reference":"Payment for Em's piano lessons",
	   "scheme_payment_sub_type":"InternetBanking",
	   "scheme_payment_type":"ImmediatePayment",
	   "sponsor_party":{  
		  "account_number":"56781234",
		  "bank_id":"123123",
		  "bank_id_code":"GBDSC"
	   }
	}
 }`

type MockPaymentRepository struct {
	err error
}

func (r MockPaymentRepository) GetAll() ([]model.Payment, error) {
	if r.err != nil {
		return nil, r.err
	}

	return []model.Payment{jsonToPayment(paymentJSON1), jsonToPayment(paymentJSON2)}, nil
}

func (r MockPaymentRepository) GetByID(paymentID string) (*model.Payment, error) {
	if r.err != nil {
		return nil, r.err
	}

	payment := jsonToPayment(paymentJSON1)
	return &payment, nil
}

func (r MockPaymentRepository) Add(payment model.Payment) (*model.Payment, error) {
	if r.err != nil {
		return nil, r.err
	}

	payment.ID = uuid.New().String()
	return &payment, nil
}

func (r MockPaymentRepository) Delete(paymentID string) error {
	if r.err != nil {
		return r.err
	}
	return nil
}

func (r MockPaymentRepository) Update(payment model.Payment) error {
	if r.err != nil {
		return r.err
	}

	return nil
}

func newMockPaymentRepository(err error) MockPaymentRepository {
	return MockPaymentRepository{
		err: err,
	}
}
func TestGetPaymentByID(t *testing.T) {
	handler := Handler{newMockPaymentRepository(nil)}

	req := buildGetByIDRequest(t)
	rr := httptest.NewRecorder() // Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.

	handler.GetPaymentByID(rr, req)

	checkResponseStatus(rr.Code, http.StatusOK, t)

	expectedPayment := jsonToPayment(paymentJSON1)
	actualPayment := jsonToPayment(rr.Body.String())

	assert.Equal(t, expectedPayment.ID, actualPayment.ID, "Payment id mismatch")
}

func TestGetPaymentByIDNotFound(t *testing.T) {
	handler := Handler{newMockPaymentRepository(ErrPaymentNotFound)}

	req := buildGetByIDRequest(t)
	rr := httptest.NewRecorder()

	handler.GetPaymentByID(rr, req)

	checkResponseStatus(rr.Code, http.StatusNotFound, t)
	checkPaymentErrorResponse(rr.Body.String(), model.PaymentNotFound, t)
}
func TestGetPaymentByIDWithError(t *testing.T) {
	handler := Handler{newMockPaymentRepository(errors.New("Generic error"))}

	req := buildGetByIDRequest(t)
	rr := httptest.NewRecorder()

	handler.GetPaymentByID(rr, req)

	checkResponseStatus(rr.Code, http.StatusInternalServerError, t)
	checkPaymentErrorResponse(rr.Body.String(), model.GenericError, t)
}

func TestAddPayment(t *testing.T) {
	handler := Handler{newMockPaymentRepository(nil)}

	req := buildAddPaymentRequest(t)
	rr := httptest.NewRecorder()

	handler.AddPayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusCreated, t)

	payment := jsonToPayment(paymentJSON3)
	createdPayment := jsonToPayment(rr.Body.String())
	assert.NotEmpty(t, createdPayment.ID, "Empty id for newly created payment")
	assert.Equal(t, payment.BeneficiaryParty.AccountName, createdPayment.BeneficiaryParty.AccountName, "Beneficiary name mismatch")
}

func TestAddPaymentWithInvalidInput(t *testing.T) {
	handler := Handler{newMockPaymentRepository(nil)}

	req := buildAddPaymentInvalidRequest(t)
	rr := httptest.NewRecorder()

	handler.AddPayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusBadRequest, t)

	response := jsonToPaymentError(rr.Body.String())
	assert.Equal(t, "INVALID_INPUT", response.Code, "Error code mismatch")
}

func TestAddPaymentWithError(t *testing.T) {
	handler := Handler{newMockPaymentRepository(errors.New("Generic error"))}

	req := buildAddPaymentRequest(t)
	rr := httptest.NewRecorder()

	handler.AddPayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusInternalServerError, t)
	checkPaymentErrorResponse(rr.Body.String(), model.GenericError, t)
}

func TestDeletePayment(t *testing.T) {
	handler := Handler{newMockPaymentRepository(nil)}

	req := buildDeletePaymentRequest(t)
	rr := httptest.NewRecorder()

	handler.DeletePayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusNoContent, t)
}

func TestDeletePaymentWithError(t *testing.T) {
	handler := Handler{newMockPaymentRepository(errors.New("Generic error"))}

	req := buildDeletePaymentRequest(t)
	rr := httptest.NewRecorder()

	handler.DeletePayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusInternalServerError, t)
	checkPaymentErrorResponse(rr.Body.String(), model.GenericError, t)
}

func TestGetAllPayments(t *testing.T) {
	handler := Handler{newMockPaymentRepository(nil)}

	req := buildGetAllPaymentsRequest(t)
	rr := httptest.NewRecorder()

	handler.GetAllPayments(rr, req)

	checkResponseStatus(rr.Code, http.StatusOK, t)
	payments := jsonToPayments(rr.Body.String())
	assert.Equal(t, 2, len(payments), "Incorrect payments count")
}

func TestGetAllPaymentsWithError(t *testing.T) {
	handler := Handler{newMockPaymentRepository(errors.New("Generic error"))}

	req := buildGetAllPaymentsRequest(t)
	rr := httptest.NewRecorder()

	handler.GetAllPayments(rr, req)

	checkResponseStatus(rr.Code, http.StatusInternalServerError, t)
	checkPaymentErrorResponse(rr.Body.String(), model.GenericError, t)
}

func TestUpdatePayment(t *testing.T) {
	handler := Handler{newMockPaymentRepository(nil)}

	req := buildUpdatePaymentRequest(t)
	rr := httptest.NewRecorder()

	handler.UpdatePayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusOK, t)
}

func TestUpdatePaymentWithInvalidInput(t *testing.T) {
	handler := Handler{newMockPaymentRepository(nil)}

	req := buildUpdatePaymentInvalidRequest(t)
	rr := httptest.NewRecorder()

	handler.UpdatePayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusBadRequest, t)

	response := jsonToPaymentError(rr.Body.String())
	assert.Equal(t, "INVALID_INPUT", response.Code, "Error code mismatch")
}

func TestUpdatePaymentWithError(t *testing.T) {
	handler := Handler{newMockPaymentRepository(errors.New("Generic error"))}

	req := buildUpdatePaymentRequest(t)
	rr := httptest.NewRecorder()

	handler.UpdatePayment(rr, req)

	checkResponseStatus(rr.Code, http.StatusInternalServerError, t)
	checkPaymentErrorResponse(rr.Body.String(), model.GenericError, t)
}

func jsonToPayment(paymentJSON string) model.Payment {
	var payment model.Payment
	if err := json.Unmarshal([]byte(paymentJSON), &payment); err != nil {
		log.Fatalf("Cannot unmarshall JSON to Payment %v", err)
	}
	return payment
}

func jsonToPayments(paymentJSON string) []model.Payment {
	var payments []model.Payment
	if err := json.Unmarshal([]byte(paymentJSON), &payments); err != nil {
		log.Fatalf("Cannot unmarshall JSON to Payment %v", err)
	}
	return payments
}

func jsonToPaymentError(paymentJSON string) model.PaymentError {
	var paymentError model.PaymentError
	if err := json.Unmarshal([]byte(paymentJSON), &paymentError); err != nil {
		log.Fatalf("Cannot unmarshall JSON to PaymentError %v", err)
	}
	return paymentError
}

func checkResponseStatus(actualStatus, expectedStatus int, t *testing.T) {
	assert.Equal(t, expectedStatus, actualStatus, "HTTP status code mismatch")
}

func checkPaymentErrorResponse(response string, expectedError model.PaymentError, t *testing.T) {
	actualError := jsonToPaymentError(response)
	assert.Equal(t, expectedError, actualError, "Error response mismatch")
}

func buildGetByIDRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "/v1/payments/7eb8277a-6c91-45e9-8a03-a27f82aca350", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func buildGetAllPaymentsRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "/v1/payments", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func buildUpdatePaymentRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("PUT", "/v1/payments/7eb8277a-6c91-45e9-8a03-a27f82aca350", strings.NewReader(paymentJSON1))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func buildUpdatePaymentInvalidRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("PUT", "/v1/payments/7eb8277a-6c91-45e9-8a03-a27f82aca350", strings.NewReader("Invalid request"))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func buildDeletePaymentRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("DELETE", "/v1/payments/7eb8277a-6c91-45e9-8a03-a27f82aca350", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func buildAddPaymentRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/v1/payments", strings.NewReader(paymentJSON3))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func buildAddPaymentInvalidRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/v1/payments", strings.NewReader("Invalid request"))
	if err != nil {
		t.Fatal(err)
	}
	return req
}
