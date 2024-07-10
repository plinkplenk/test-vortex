package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/plinkplenk/test-vortex/internal/api/schemas"
	"github.com/plinkplenk/test-vortex/internal/api/validators"
	"github.com/plinkplenk/test-vortex/internal/orders"
	mock_service "github.com/plinkplenk/test-vortex/internal/orders/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var loggerStub = slog.New(slog.NewTextHandler(io.Discard, nil))

var (
	exchangeNameNotProvidedResponse, _ = json.Marshal(j{"error": validators.ErrExchangeNameNotProvided.Error()})
	pairNotProvidedResponse, _         = json.Marshal(j{"error": validators.ErrPairNotProvided.Error()})
	depthNotProvidedResponse, _        = json.Marshal(j{"error": validators.ErrDepthNotProvided.Error()})
	clientNotProvidedResponse, _       = json.Marshal(j{"error": validators.ErrClientNotProvided.Error()})
	orderHistoryNotProvidedResponse, _ = json.Marshal(j{"error": validators.ErrOrderHistoryNotProvided.Error()})
	orderBookNotFoundResponse, _       = json.Marshal(j{"error": ErrOrderBookNotFound.Error()})
	orderHistoryNotFoundResponse       = func(label, pair string) string {
		r, _ := json.Marshal(
			j{
				"error": fmt.Sprintf(
					"Order history with lable[%s] and pair[%s] not found",
					label,
					pair,
				),
			},
		)
		return string(r)
	}
	labelOrPairNotProvidedResponse, _ = json.Marshal(j{"error": ErrLabelOrPairNotProvided.Error()})
)

func TestOrdersHandler_GetOrderHistory(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrdersService, client orders.Client)
	testTable := []struct {
		name               string
		mockBehavior       mockBehavior
		inputClient        orders.Client
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "SUCCESS",
			mockBehavior: func(s *mock_service.MockOrdersService, client orders.Client) {
				s.EXPECT().GetOrderHistory(context.Background(), client).Return(
					[]*orders.History{
						{
							Client:              client,
							Side:                "side",
							Type:                "type",
							BaseQty:             1,
							Price:               0.1,
							AlgorithmNamePlaced: "algo",
							LowestSellPrc:       0.09,
							HighestBuyPrc:       1,
							CommissionQuoteQty:  0.005,
							TimePlaced:          time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC),
						},
					},
					nil,
				)
			},
			inputClient: orders.Client{
				ClientName:   "client",
				ExchangeName: "exchange",
				Label:        "label",
				Pair:         "A_B",
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"orderHistory":[{"clientName":"client","exchangeName":"exchange","label":"label","pair":"A_B","side":"side","type":"type","baseQty":1,"price":0.1,"algorithmNamePlaced":"algo","lowestSellPrc":0.09,"highestBuyPrc":1,"commissionQuoteQty":0.005,"timePlaced":"2024-01-01T01:01:01Z"}]}`,
		},
		{
			name: "NOT FOUND",
			mockBehavior: func(s *mock_service.MockOrdersService, client orders.Client) {
				s.EXPECT().GetOrderHistory(context.Background(), client).Return(
					nil,
					nil,
				)
			},
			inputClient: orders.Client{
				ClientName:   "client",
				ExchangeName: "exchange",
				Label:        "label",
				Pair:         "A_B",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       orderHistoryNotFoundResponse("label", "A_B"),
		},
		{
			name: "INVALID INPUT (NO PAIR)",
			mockBehavior: func(s *mock_service.MockOrdersService, client orders.Client) {
				s.EXPECT().GetOrderHistory(context.Background(), client).Times(0)
			},
			inputClient: orders.Client{
				ClientName:   "client",
				ExchangeName: "exchange",
				Label:        "label",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       string(labelOrPairNotProvidedResponse),
		},
		{
			name: "INVALID INPUT (NO LABEL)",
			mockBehavior: func(s *mock_service.MockOrdersService, client orders.Client) {
				s.EXPECT().GetOrderHistory(context.Background(), client).Times(0)
			},
			inputClient: orders.Client{
				ClientName:   "client",
				ExchangeName: "exchange",
				Pair:         "A_B",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       string(labelOrPairNotProvidedResponse),
		},
	}

	for _, test := range testTable {
		t.Run(
			test.name, func(t *testing.T) {
				c := gomock.NewController(t)
				defer c.Finish()

				orderService := mock_service.NewMockOrdersService(c)
				test.mockBehavior(orderService, test.inputClient)

				handler := NewOrdersHandler(orderService, loggerStub)
				router := chi.NewRouter()
				router.Get("/history/{client_name}/{exchange_name}", handler.GetOrderHistory(context.Background()))

				w := httptest.NewRecorder()
				path, _ := url.JoinPath("/history", test.inputClient.ClientName, test.inputClient.ExchangeName)
				path += fmt.Sprintf("?label=%s&pair=%s", test.inputClient.Label, test.inputClient.Pair)
				r := httptest.NewRequest(http.MethodGet, path, bytes.NewBufferString(""))
				router.ServeHTTP(w, r)

				assert.Equal(t, test.expectedStatusCode, w.Code)
				assert.Equal(t, test.expectedBody, w.Body.String())
			},
		)
	}
}

func TestOrdersHandler_GetOrderBook(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrdersService, exchangeName string, pair string)
	testTable := []struct {
		name               string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "SUCCESS",
			mockBehavior: func(s *mock_service.MockOrdersService, exchangeName string, pair string) {
				s.EXPECT().GetOrderBook(context.Background(), exchangeName, pair).Return(
					[]orders.Depth{
						{
							Price:   0.5,
							BaseQty: 1,
						},
					},
					nil,
				)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"depthOrders":[{"price":0.5,"baseQty":1}]}`,
		},
		{
			name: "NOT FOUND",
			mockBehavior: func(s *mock_service.MockOrdersService, exchangeName string, pair string) {
				s.EXPECT().GetOrderBook(context.Background(), exchangeName, pair).Return(
					nil,
					nil,
				)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       string(orderBookNotFoundResponse),
		},
	}

	for _, test := range testTable {
		t.Run(
			test.name, func(t *testing.T) {
				c := gomock.NewController(t)
				defer c.Finish()

				exchangeName := "some-exchange"
				pair := "A_B"
				orderService := mock_service.NewMockOrdersService(c)
				test.mockBehavior(orderService, exchangeName, pair)

				handler := NewOrdersHandler(orderService, loggerStub)
				router := chi.NewRouter()
				router.Get("/{exchange_name}/{pair}", handler.GetOrderBook(context.Background()))

				w := httptest.NewRecorder()
				path, _ := url.JoinPath("/", exchangeName, pair)
				r := httptest.NewRequest(http.MethodGet, path, bytes.NewBufferString(""))
				router.ServeHTTP(w, r)

				assert.Equal(t, test.expectedStatusCode, w.Code)
				assert.Equal(t, test.expectedBody, w.Body.String())
			},
		)
	}
}

func TestOrdersHandler_SaveOrderBook(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrdersService, orderBook schemas.OrderBookCreate)
	testTable := []struct {
		name               string
		inputBody          string
		inputOrderBook     schemas.OrderBookCreate
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "SUCCESS",
			inputBody: `{"exchangeName": "some-exchange", "pair": "A_B", "depth": [{ "price": 0.01, "baseQty": 0.01 }]}`,
			inputOrderBook: schemas.OrderBookCreate{
				ExchangeName: "some-exchange",
				Pair:         "A_B",
				Depth:        []orders.Depth{{Price: 0.01, BaseQty: 0.01}},
			},
			mockBehavior: func(s *mock_service.MockOrdersService, orderBook schemas.OrderBookCreate) {
				s.EXPECT().SaveOrderBook(
					context.Background(),
					orderBook.ExchangeName,
					orderBook.Pair,
					orderBook.Depth,
				).Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "",
		},
		{
			name:      "INVALID INPUT (NO EXCHANGE NAME)",
			inputBody: `{ "pair": "A_B", "depth": [ { "price": 0.01, "baseQty": 0.01 } ]}`,
			inputOrderBook: schemas.OrderBookCreate{
				Pair:  "A_B",
				Depth: []orders.Depth{{Price: 0.01, BaseQty: 0.01}},
			},
			mockBehavior: func(s *mock_service.MockOrdersService, orderBook schemas.OrderBookCreate) {
				s.EXPECT().SaveOrderBook(
					context.Background(),
					orderBook.ExchangeName, orderBook.Pair, orderBook.Depth,
				).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       string(exchangeNameNotProvidedResponse),
		},
		{
			name:      "INVALID INPUT (NO PAIR)",
			inputBody: `{"exchangeName": "some-exchange", "depth": [{ "price": 0.01, "baseQty": 0.01 }]}`,
			inputOrderBook: schemas.OrderBookCreate{
				ExchangeName: "some-exchange",
				Depth:        []orders.Depth{{Price: 0.01, BaseQty: 0.01}},
			},
			mockBehavior: func(s *mock_service.MockOrdersService, orderBook schemas.OrderBookCreate) {
				s.EXPECT().SaveOrderBook(
					context.Background(),
					orderBook.ExchangeName,
					orderBook.Pair,
					orderBook.Depth,
				).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       string(pairNotProvidedResponse),
		},
		{
			name:      "INVALID INPUT (NO DEPTH)",
			inputBody: `{"exchangeName": "some-exchange", "pair": "A_B"}`,
			inputOrderBook: schemas.OrderBookCreate{
				ExchangeName: "some-exchange",
				Pair:         "A_B",
			},
			mockBehavior: func(s *mock_service.MockOrdersService, orderBook schemas.OrderBookCreate) {
				s.EXPECT().SaveOrderBook(
					context.Background(),
					orderBook.ExchangeName,
					orderBook.Pair,
					orderBook.Depth,
				).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       string(depthNotProvidedResponse),
		},
	}
	for _, test := range testTable {
		t.Run(
			test.name, func(t *testing.T) {
				c := gomock.NewController(t)
				defer c.Finish()

				orderService := mock_service.NewMockOrdersService(c)
				test.mockBehavior(orderService, test.inputOrderBook)

				handler := NewOrdersHandler(orderService, loggerStub)
				router := chi.NewRouter()
				router.Post("/orders", handler.SaveOrderBook(context.Background()))

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(test.inputBody))
				router.ServeHTTP(w, r)
				assert.Equal(t, test.expectedStatusCode, w.Code)
				assert.Equal(t, test.expectedBody, w.Body.String())

			},
		)
	}
}

func TestOrdersHandler_SaveOrder(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrdersService, client orders.Client, history *orders.History)
	testTable := []struct {
		name               string
		inputBody          string
		inputOrderHistory  orders.History
		inputOrderClient   orders.Client
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "SUCCESS",
			inputBody: `{
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
						"price": 0.1,
						"algorithmNamePlaced": "algo",
						"lowestSellPrc": 0.1,
						"highestBuyPrc": 1,
						"commissionQuoteQty": 0.01
					}
				}`,
			inputOrderClient: orders.Client{
				ClientName:   "client",
				ExchangeName: "some-exchange",
				Label:        "label",
				Pair:         "A_B",
			},
			inputOrderHistory: orders.History{
				Side:                "some-side",
				Type:                "some-type",
				BaseQty:             1,
				Price:               0.1,
				AlgorithmNamePlaced: "algo",
				LowestSellPrc:       0.1,
				HighestBuyPrc:       1,
				CommissionQuoteQty:  0.01,
			},
			mockBehavior: func(s *mock_service.MockOrdersService, client orders.Client, history *orders.History) {
				s.EXPECT().SaveOrder(
					context.Background(),
					client,
					history,
				).Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "",
		}, {
			name: "INVALID INPUT (NO CLIENT)",
			inputBody: `{
			    "orderHistory": {
			        "side": "some-side",
			        "type": "some-type",
			        "baseQty": 1,
			        "price": 0.1,
			        "algorithmNamePlaced": "algo",
			        "lowestSellPrc": 0.1,
			        "highestBuyPrc": 1,
			        "commissionQuoteQty": 0.01
			    }
			}`,
			inputOrderClient: orders.Client{},
			inputOrderHistory: orders.History{
				Side:                "some-side",
				Type:                "some-type",
				BaseQty:             1,
				Price:               0.1,
				AlgorithmNamePlaced: "algo",
				LowestSellPrc:       0.1,
				HighestBuyPrc:       1,
				CommissionQuoteQty:  0.01,
			},
			mockBehavior: func(s *mock_service.MockOrdersService, client orders.Client, history *orders.History) {
				s.EXPECT().SaveOrder(
					context.Background(),
					client,
					history,
				).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       string(clientNotProvidedResponse),
		},
		{
			name: "INVALID INPUT (NO ORDER HISTORY)",
			inputBody: `{
				"client": {
        			"clientName": "client",
        			"exchangeName": "some-exchange",
        			"label": "label",
        			"pair": "A_B"
    			}
			}`,
			inputOrderClient: orders.Client{
				ClientName: "client", ExchangeName: "some-exchange", Label: "label", Pair: "A_B",
			},
			inputOrderHistory: orders.History{},
			mockBehavior: func(s *mock_service.MockOrdersService, client orders.Client, history *orders.History) {
				s.EXPECT().SaveOrder(
					context.Background(),
					client,
					history,
				).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       string(orderHistoryNotProvidedResponse),
		},
	}
	for _, test := range testTable {
		t.Run(
			test.name, func(t *testing.T) {
				c := gomock.NewController(t)
				defer c.Finish()

				orderService := mock_service.NewMockOrdersService(c)
				test.mockBehavior(orderService, test.inputOrderClient, &test.inputOrderHistory)
				handler := NewOrdersHandler(orderService, loggerStub)
				router := chi.NewRouter()
				router.Post("/history", handler.SaveOrder(context.Background()))

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/history", bytes.NewBufferString(test.inputBody))
				router.ServeHTTP(w, r)

				assert.Equal(t, test.expectedStatusCode, w.Code)
				assert.Equal(t, test.expectedBody, w.Body.String())
			},
		)
	}
}
