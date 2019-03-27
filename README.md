# Payment APP

## 1. Description
Payment App provides a RESTful API for managing Form3 payments.

## 2. Dependencies
* MongoDb 

## 3. Endpoints
The API provides the following endpoints for managing payments:
* `GET /v1/payments` - retrieve all payments
* `GET /v1/payments/{paymentID}` - retrieve the payment having the provided ID
* `POST /v1/payments` - create new payment
* `PUT /v1/payments/{paymentID}` - update payment 
* `DELETE /v1/payments/{paymentID}` - delete the payment with the provided ID

## 4. How to run the project locally
### Running the project locally without Docker
* Prerequisite:
   * GO must be installed on your machine and the GOPATH set accordingly
   * MongoDB should either be installed and running on your machine or a remote MongoDB should be accessible from your local machine. The address of the running MongoDB instance is to be passed to the application using the environment variable `MONGO_URI`

Run the following command from terminal:
```
MONGO_URI=localhost:27017 go run main.go
```
(The value for MONGO_URI should be the address of the MongoDB instance you are connecting to)
### Running the project with Docker
* Prerequisite:
   * Docker should be installed on your local machine

1. Build the image for the payment app using:
```
docker build -t payment-app .
```
2. Run docker compose in order to spin a container for the payment app and a container for mongodb:
```
docker-compose up
```

### Check the application is working correctly
Perform a curl to any of the above endpoints:
i.e. `curl -i "localhost:8080/v1/payments"`

## 5. How to run the tests
### Running unit tests
```
go test -v ./...
```
### Running integration tests
1. Run docker-compose in order to spin up a MongoDB container (as described above)
2. Instal Godog `go get github.com/DATA-DOG/godog/cmd/godog`
3. Run the following command in terminal `godog features/payment.feature`

## 6. API Examples

Create new payment

```
POST /v1/payments
{  
   "type":"Payment",
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
}
```


Update payment

```
PUT /v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43
{  
   "type":"Payment",
   "id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
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
}
```

## 7. To Do
- add health check endpoint
- persist payment id as Mongo DB UUID field
