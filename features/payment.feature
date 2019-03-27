Feature: manage payments
  As a payment API user
  I would like to perform CRUP operations on payments
  So that I can fulfill customers' payment requests

  Scenario: should successfully retrieve a newly created payment
    Given a payment:
    """
      {
        "type":"Payment",
        "version":0,
        "organisation_id":"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
        "attributes":{
            "amount":"100.21",
            "beneficiary_party":{
              "account_name":"Wilfred Owens",
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
    """
    When the user requests to add the payment
    Then the payment is successfully saved
    And the payment has an id assigned
    When the user requests to retrieve the payment
    Then the payment is successfully retrieved
    And the retrieved payment details are identical to the provided one

  Scenario: should successfully delete a payment
    Given a payment:
    """
      {
        "type":"Payment",
        "version":0,
        "organisation_id":"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
        "attributes":{
            "amount":"306.5",
            "beneficiary_party":{
              "account_name":"J Doe",
              "account_number":"81226819",
              "account_number_code":"BBAN",
              "account_type":0,
              "address":"Alberta street SE17",
              "bank_id":"320000",
              "bank_id_code":"GBDSC",
              "name":"Jane M Doe"
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
            "end_to_end_reference":"Season tickets",
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
            "reference":"Payment for Jane Brown's season ticket",
            "scheme_payment_sub_type":"InternetBanking",
            "scheme_payment_type":"ImmediatePayment",
            "sponsor_party":{
              "account_number":"56781234",
              "bank_id":"123123",
              "bank_id_code":"GBDSC"
            }
        }
      }
    """
    When the user requests to add the payment
    Then the payment is successfully saved
    And the payment has an id assigned
    When the user requests to retrieve the payment
    Then the payment is successfully retrieved
    When the user requests to delete the payment
    Then the payment is successfully deleted
    When the user requests to retrieve the payment
    Then the payment is not found