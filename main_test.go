package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/DATA-DOG/godog/gherkin"

	"github.com/tenzi-muraru/paymentApp/payment/model"

	"github.com/DATA-DOG/godog"
)

const (
	paymentAppHost = "http://localhost:9999"

	paymentEndpoint = "/v1/payments"
)

type paymentFeature struct {
	resp         *http.Response
	inputPayment *model.Payment
	paymentID    string
}

func (f *paymentFeature) resetResponse(interface{}) {
	f.resp = nil
	f.inputPayment = nil
}

func (f *paymentFeature) aPayment(body *gherkin.DocString) error {
	var inputPayment model.Payment
	if err := json.Unmarshal([]byte(body.Content), &inputPayment); err != nil {
		return fmt.Errorf("Cannot unmarshal provided payment: %v", err)
	}

	f.inputPayment = &inputPayment
	return nil
}

func (f *paymentFeature) theUserRequestsToAddThePayment() error {
	paymentJSON, err := json.Marshal(f.inputPayment)
	if err != nil {
		return fmt.Errorf("Unable to marshal payment input: %v", err)
	}

	resp, err := http.Post(paymentAppHost+paymentEndpoint, "application/json", bytes.NewBuffer(paymentJSON))
	if err != nil {
		return fmt.Errorf("Unable to perform POST call to %s: %v", paymentEndpoint, err)
	}

	f.resp = resp
	return nil
}

func (f *paymentFeature) thePaymentIsSuccessfullySaved() error {
	return f.checkResponseCode(http.StatusCreated)
}

func (f *paymentFeature) thePaymentHasAnIDAssigned() error {
	payment, err := getPaymentResponse(f.resp.Body)
	if err != nil {
		return err
	}

	if len(payment.ID) == 0 {
		return fmt.Errorf("Missing payment ID on newly created payment")
	}

	f.paymentID = payment.ID
	return nil
}

func (f *paymentFeature) theUserRequestsToRetrieveThePayment() error {
	resp, err := performGetCall(fmt.Sprintf("%s/%s", paymentEndpoint, f.paymentID))
	if err != nil {
		return err
	}
	f.resp = resp
	return nil
}

func (f *paymentFeature) thePaymentIsSuccessfullyRetrieved() error {
	return f.checkResponseCode(http.StatusOK)
}

func (f *paymentFeature) theRetrievedPaymentDetailsAreIdenticalToTheProvidedOne() error {
	outputPayment, err := getPaymentResponse(f.resp.Body)
	if err != nil {
		return err
	}

	if outputPayment.ID != f.paymentID {
		return fmt.Errorf("Incorrect payment ID. Expected %s - actual %s", f.paymentID, outputPayment.ID)
	}

	if outputPayment.BeneficiaryParty != f.inputPayment.BeneficiaryParty {
		return fmt.Errorf("Incorrect payment ID. Expected %v - actual %v", f.inputPayment.BeneficiaryParty, outputPayment.BeneficiaryParty)
	}

	return nil
}

func (f *paymentFeature) theUserRequestsToDeleteThePayment() error {
	endpoint := fmt.Sprintf("%s/%s", paymentEndpoint, f.paymentID)
	deleteReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", paymentAppHost, endpoint), nil)
	if err != nil {
		return fmt.Errorf("Unable to create DELETE request to %s: %v", endpoint, err)
	}

	client := &http.Client{}
	f.resp, err = client.Do(deleteReq)
	if err != nil {
		return fmt.Errorf("Unable to create DELETE call to %s: %v", endpoint, err)
	}
	return nil
}

func (f *paymentFeature) thePaymentIsSuccessfullyDeleted() error {
	return f.checkResponseCode(http.StatusNoContent)
}

func (f *paymentFeature) thePaymentIsNotFound() error {
	return f.checkResponseCode(http.StatusNotFound)
}

func (f *paymentFeature) checkResponseCode(statusCode int) error {
	if f.resp.StatusCode != statusCode {
		return fmt.Errorf("Incorrect status code. Expected %d - actual %d", statusCode, f.resp.StatusCode)
	}

	return nil
}

func performGetCall(endpoint string) (*http.Response, error) {
	resp, err := http.Get(paymentAppHost + endpoint)
	if err != nil {
		return nil, fmt.Errorf("Unable to perform GET call to %s: %v", endpoint, err)
	}
	return resp, nil
}

func getPaymentResponse(responseBody io.ReadCloser) (*model.Payment, error) {
	var result model.Payment
	if err := json.NewDecoder(responseBody).Decode(&result); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal payment response: %v", err)
	}
	return &result, nil
}

func FeatureContext(s *godog.Suite) {
	// start payment API on localhost:9999
	os.Setenv("MONGO_URI", "localhost:27017")
	ctx, cancel := context.WithCancel(context.Background())
	go startApp(ctx, ":9999")
	defer cancel()
	// wait for the payment API to start
	time.Sleep(2 * time.Second)

	f := &paymentFeature{}

	s.BeforeScenario(f.resetResponse)

	s.Step(`^a payment:$`, f.aPayment)
	s.Step(`^the user requests to add the payment$`, f.theUserRequestsToAddThePayment)
	s.Step(`^the payment is successfully saved$`, f.thePaymentIsSuccessfullySaved)
	s.Step(`^the payment has an id assigned$`, f.thePaymentHasAnIDAssigned)
	s.Step(`^the user requests to retrieve the payment$`, f.theUserRequestsToRetrieveThePayment)
	s.Step(`^the payment is successfully retrieved$`, f.thePaymentIsSuccessfullyRetrieved)
	s.Step(`^the retrieved payment details are identical to the provided one$`, f.theRetrievedPaymentDetailsAreIdenticalToTheProvidedOne)

	s.Step(`^the user requests to delete the payment$`, f.theUserRequestsToDeleteThePayment)
	s.Step(`^the payment is successfully deleted$`, f.thePaymentIsSuccessfullyDeleted)
	s.Step(`^the payment is not found$`, f.thePaymentIsNotFound)
}
