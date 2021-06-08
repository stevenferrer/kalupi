# Kalupi REST API

**Table of Contents**
----
  - [**Create wallet account**](#create-wallet-account)
  - [**Get wallet account**](#get-wallet-account)
  - [**List wallet accounts**](#list-wallet-accounts)
  - [**Make cash deposit**](#make-cash-deposit)
  - [**Make cash withdrawal**](#make-cash-withdrawal)
  - [**Make cash payment**](#make-cash-payment)
  - [**List cash payments**](#list-cash-payments)

**Create wallet account**
----
  Creates a wallet account.

* **URL**

  `/accounts`

* **Method:**

  `POST`
  
*  **URL Params**

    None

* **Data Params**

    ```json
    {
        "account_id": [alphanumeric],
        "currency": [ISO 4217 e.g. USD]
    }
    ```

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** None
 
* **Error Response:**

  * **Code:** 422 UNPROCESSABLE ENTITY <br />
    **Content:**
    ```json
    {
      "error": "validation error; account_id: must have length between 6 and 64."
    }
    ```
    or
    ```json
    {
      "error": "account already exist"
    }
    ```


**Get wallet account**
----
  Retrieves the wallet account

* **URL**

  `/accounts/{account_id}`

* **Method:**

  `GET`
  
* **URL Params**

  None

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 
    ```json
    {
      "account": {
        "id": "johndoe",
        "currency": "USD",
        "balance": "56.068"
      }
    }
    ```
 
* **Error Response:**

  * **Code** 404 NOT FOUND <br />
    **Content:**
    ```json
    {
      "error": "account not found"
    }
    ```

**List wallet accounts**
----
  Retrieves all wallet accounts

* **URL**

  `/accounts`

* **Method:**

  `GET`
  
* **URL Params**

  None

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 
    ```json
    {
      "accounts": [
        {
          "id": "johndoe",
          "currency": "USD",
          "balance": "56.068"
        },
        {
          "id": "maryjane",
          "currency": "USD",
          "balance": "10.398"
        }
      ]
    }
    ```
 
* **Error Response:**

  * **Code** 500 INTERNAL SERVER ERROR <br />
    **Content:**
    ```json
    {
      "error": "internal server error"
    }
    ```

**Make cash deposit**
----
  Make cash deposit.

* **URL**

  `/t/deposit`

* **Method:**

  `POST`
  
* **URL Params**

  None

* **Data Params**

    ```json
    {
        "account_id": [alphanumeric],
        "amount": [Non-zero and non-negative decimal]
    }
    ```

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** None
 
* **Error Response:**

  * **Code** 500 INTERNAL SERVER ERROR <br />
    **Content:**
    ```json
    {
      "error": "internal server error"
    }
    ```

**Make cash withdrawal**
----
  Make cash withdrawal.

* **URL**

  `/t/withdraw`

* **Method:**

  `POST`
  
* **URL Params**

  None

* **Data Params**

    ```json
    {
        "account_id": [alphanumeric],
        "amount": [Non-zero and non-negative decimal]
    }
    ```

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** None
 
* **Error Response:**

  * **Code** 422 UNPROCESSABLE ENTITY <br />
    **Content:**
    ```json
    {
      "error": "insufficient balance"
    }
    ```

  * **Code** 500 INTERNAL SERVER ERROR <br />
    **Content:**
    ```json
    {
      "error": "internal server error"
    }
    ```

**Make cash payment**
----
  Make cash payment.

* **URL**

  `/t/payments`

* **Method:**

  `POST`
  
* **URL Params**

  None

* **Data Params**

    ```json
    {
        "from_account": [alphanumeric],
        "to_account": [alphanumeric],
        "amount": [Non-zero, non-negative decimal]
    }
    ```

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** None
 
* **Error Response:**

  * **Code** 422 UNPROCESSABLE ENTITY<br />
    **Content:**
    ```json
    {
      "error": "sending account not found"
    }
    ```
    or
    ```json
    {
      "error": "receiving account not found"
    }
    ```
    or
    ```json
    {
      "error": "insufficient balance"
    }
    ```

**List cash payments**
----
  List cash payments.

* **URL**

  `/t/payments`

* **Method:**

  `GET`
  
* **URL Params**

  None

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:**
    ```json
    {
      "payments": [
        {
          "xact_no": "LM4I8FHC05X0",
          "account": "johndoe",
          "amount": "23.938",
          "direction": "outgoing",
          "to_account": "maryjane"
        },
        {
          "xact_no": "LM4I8FHC05X0",
          "account": "maryjane",
          "amount": "23.938",
          "direction": "incoming",
          "from_account": "johndoe"
        },
        ...
      ]
    }
    ```
 
* **Error Response:**

  * **Code** 500 INTERNAL SERVER ERROR <br />
    **Content:**
    ```json
    {
      "error": "internal server error"
    }
    ```