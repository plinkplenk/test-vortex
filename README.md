# test-vortex

## To run this project
 
- rename example.env
    ```bash
    mv example.env .env
    ```
 
- setup .env

- install dependencies
    ```bash
  go mod download
    ```

- run migrations
    ```bash
    make migrations-up
    ```

- run app
     ```bash
    make run
     ```
- or run via docker compose
    ```bash
    docker compose up -d
    ```
  
## Endpoints

-  **[GET] /orders/{exchange}/{pair}**
    ```bash
    curl --location 'http://localhost:8080/orders/{exchange}/{pair}'
    ```

- **[POST] /orders**
    ```bash
    curl --location 'http://localhost:8080/orders' \
    --header 'Content-Type: application/json' \
    --data '{   
      "exchangeName": "some-exchange",
      "pair": "A_B",
      "depth": [       
        {
            "price": 0.01,
            "baseQty": 1
        }
      ]
    }'
    ```

- **[GET] /orders/history/{clientName}/{exchangeName}?label={label}&pair={pair}**
    ```bash
    curl --location 'http://localhost:8080/orders/history/{clientName}/{exchangeName}?label={label}&pair={pair}'
    ```
 
- **[POST] /orders/history**
    ```bash
    curl --location 'http://localhost:8080/orders/history/' \
    --header 'Content-Type: application/json' \
    --data '{
      "client": {
        "clientName": "client",
        "exchangeName": "some-exchange",
        "label": "label",
        "pair": "A_B"
      },
      "orderHistory": {
        "side": "some-side",
        "type": "some-type",
        "baseQty": 1,
        "price": 0.01,
        "algorithmNamePlaced": "algo",
        "lowestSellPrc": 0.0099,
        "highestBuyPrc": 1.0,
        "commissionQuoteQty": 0.00001
      }
    }'
    ```
