package backend

import (
	"errors"
	"testing"
	"time"

	"github.com/okex/okchain-go-sdk/module/backend/types"
	"github.com/okex/okchain-go-sdk/types/params"

	"github.com/golang/mock/gomock"
	"github.com/okex/okchain-go-sdk/mocks"
	sdk "github.com/okex/okchain-go-sdk/types"
	"github.com/stretchr/testify/require"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	addr    = "okchain1dcsxvxgj374dv3wt9szflf9nz6342juzzkjnlz"
	product = "btc-000_okt"
)

func TestBackendClient_QueryCandles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := sdk.NewClientConfig("testURL", "testChain", sdk.BroadcastBlock, "0.01okt", 200000)
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewBackendClient(mockCli.MockBaseClient))

	mockCandles := [][]string{
		{"1.024", "2.048", "4.096", "8.192"},
		{"10.24", "20.48", "40.96", "81.92"},
	}
	granularity, size := 60, 1

	expectedRet := mockCli.BuildBackendCandlesBytes(mockCandles)
	expectedCdc := mockCli.GetCodec()

	queryParams := params.NewQueryKlinesParams(product, granularity, size)
	queryBytes := expectedCdc.MustMarshalJSON(queryParams)

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(3)
	mockCli.EXPECT().Query(types.CandlesPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil)

	candles, err := mockCli.Backend().QueryCandles(product, granularity, size)
	require.NoError(t, err)
	require.Equal(t, "1.024", candles[0][0])
	require.Equal(t, "2.048", candles[0][1])
	require.Equal(t, "4.096", candles[0][2])
	require.Equal(t, "8.192", candles[0][3])
	require.Equal(t, "10.24", candles[1][0])
	require.Equal(t, "20.48", candles[1][1])
	require.Equal(t, "40.96", candles[1][2])
	require.Equal(t, "81.92", candles[1][3])

	mockCli.EXPECT().Query(types.CandlesPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryCandles(product, granularity, size)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.CandlesPath, cmn.HexBytes(queryBytes)).Return(append(expectedRet, '}'), nil)
	_, err = mockCli.Backend().QueryCandles(product, granularity, size)
	require.Error(t, err)
}

func TestBackendClient_QueryTickers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := sdk.NewClientConfig("testURL", "testChain", sdk.BroadcastBlock, "0.01okt", 200000)
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewBackendClient(mockCli.MockBaseClient))

	open, close, high, low := "2.048", "2.048", "4.096", "1.024"
	timestamp := time.Now().String()
	price, volume, change := "2.048", "1024", "100"
	expectedRet := mockCli.BuildBackendTickersBytes(product, product, timestamp, open, close, high, low, price, volume, change)
	expectedCdc := mockCli.GetCodec()

	queryParams := params.NewQueryTickerParams(product, 10, true)
	queryBytes := expectedCdc.MustMarshalJSON(queryParams)

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(7)
	mockCli.EXPECT().Query(types.TickersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil).Times(2)

	tickers, err := mockCli.Backend().QueryTickers(product)
	require.NoError(t, err)

	_, err = mockCli.Backend().QueryTickers(product, 10)
	require.NoError(t, err)
	require.Equal(t, product, tickers[0].Symbol)
	require.Equal(t, product, tickers[0].Product)
	require.Equal(t, timestamp, tickers[0].Timestamp)
	require.Equal(t, open, tickers[0].Open)
	require.Equal(t, close, tickers[0].Close)
	require.Equal(t, high, tickers[0].High)
	require.Equal(t, low, tickers[0].Low)
	require.Equal(t, price, tickers[0].Price)
	require.Equal(t, volume, tickers[0].Volume)
	require.Equal(t, change, tickers[0].Change)

	mockCli.EXPECT().Query(types.TickersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryTickers(product)
	require.Error(t, err)

	queryParams = params.NewQueryTickerParams("", 10, true)
	queryBytes = expectedCdc.MustMarshalJSON(queryParams)
	mockCli.EXPECT().Query(types.TickersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil).Times(2)

	_, err = mockCli.Backend().QueryTickers("")
	require.NoError(t, err)

	tickers, err = mockCli.Backend().QueryTickers("", 10)
	require.NoError(t, err)
	require.Equal(t, product, tickers[0].Symbol)
	require.Equal(t, product, tickers[0].Product)
	require.Equal(t, timestamp, tickers[0].Timestamp)
	require.Equal(t, open, tickers[0].Open)
	require.Equal(t, close, tickers[0].Close)
	require.Equal(t, high, tickers[0].High)
	require.Equal(t, low, tickers[0].Low)
	require.Equal(t, price, tickers[0].Price)
	require.Equal(t, volume, tickers[0].Volume)
	require.Equal(t, change, tickers[0].Change)

	_, err = mockCli.Backend().QueryTickers(product, 1, 1)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTickers("", 1, 1)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTickers(product, -1)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTickers("", -1)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.TickersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryTickers("")
	require.Error(t, err)

	mockCli.EXPECT().Query(types.TickersPath, cmn.HexBytes(queryBytes)).Return(append(expectedRet, '}'), nil)
	_, err = mockCli.Backend().QueryTickers("")
	require.Error(t, err)

}

func TestBackendClient_QueryDeals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := sdk.NewClientConfig("testURL", "testChain", sdk.BroadcastBlock, "0.01okt", 200000)
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewBackendClient(mockCli.MockBaseClient))

	timestamp, blockHeight := time.Now().Unix()*1000, int64(1024)
	orderID, side, fee := "ID0000000000-1", "BUY", "0.00100000btc-000"
	price, quantity := 1024.1024, 2048.2048
	start, end, page, perPage := 0, 0, 1, 30

	expectedRet := mockCli.BuildBackendDealsResultBytes(timestamp, blockHeight, orderID, addr, product, side, fee, price, quantity)
	expectedCdc := mockCli.GetCodec()

	queryParams := params.NewQueryDealsParams(addr, product, int64(start), int64(end), page, perPage, side)
	queryBytes := expectedCdc.MustMarshalJSON(queryParams)

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(3)
	mockCli.EXPECT().Query(types.DealsPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil)

	deals, err := mockCli.Backend().QueryDeals(addr, product, side, start, end, page, perPage)
	require.NoError(t, err)
	require.Equal(t, timestamp, deals[0].Timestamp)
	require.Equal(t, blockHeight, deals[0].BlockHeight)
	require.Equal(t, orderID, deals[0].OrderID)
	require.Equal(t, addr, deals[0].Sender)
	require.Equal(t, product, deals[0].Product)
	require.Equal(t, side, deals[0].Side)
	require.Equal(t, price, deals[0].Price)
	require.Equal(t, quantity, deals[0].Quantity)
	require.Equal(t, fee, deals[0].Fee)

	_, err = mockCli.Backend().QueryDeals(addr[1:], product, side, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryDeals(addr, "", side, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryDeals(addr, product, "BUY&&SELL", start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryDeals(addr, product, side, end+1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryDeals(addr, product, side, -1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryDeals(addr, product, side, start, -1, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryDeals(addr, product, side, start, end, -1, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryDeals(addr, product, side, start, end, page, -1)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.DealsPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryDeals(addr, product, side, start, end, page, perPage)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.DealsPath, cmn.HexBytes(queryBytes)).Return(expectedRet[1:], nil)
	_, err = mockCli.Backend().QueryDeals(addr, product, side, start, end, page, perPage)
	require.Error(t, err)

}

func TestBackendClient_QueryOpenOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := sdk.NewClientConfig("testURL", "testChain", sdk.BroadcastBlock, "0.01okt", 200000)
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewBackendClient(mockCli.MockBaseClient))

	txHash, orderID, side := "default tx hash", "ID0000000000-1", "BUY"
	price, quantity, filledAvgQuantity, remainQuantity := "1024.1024", "2048.2048", "1024.1024", "4096.4096"
	status, timestamp := int64(0), time.Now().Unix()*1000
	start, end, page, perPage := 0, 0, 1, 30

	expectedRet := mockCli.BuildBackendOrdersResultBytes(txHash, orderID, addr, product, side, price, quantity, filledAvgQuantity,
		remainQuantity, status, timestamp)
	expectedCdc := mockCli.GetCodec()

	queryParams := params.NewQueryOrderListParams(addr, product, side, page, perPage, int64(start), int64(end), false)
	queryBytes := expectedCdc.MustMarshalJSON(queryParams)

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(3)
	mockCli.EXPECT().Query(types.OpenOrdersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil)

	openOrders, err := mockCli.Backend().QueryOpenOrders(addr, product, side, start, end, page, perPage)
	require.NoError(t, err)
	require.Equal(t, txHash, openOrders[0].TxHash)
	require.Equal(t, orderID, openOrders[0].OrderID)
	require.Equal(t, addr, openOrders[0].Sender)
	require.Equal(t, product, openOrders[0].Product)
	require.Equal(t, side, openOrders[0].Side)
	require.Equal(t, price, openOrders[0].Price)
	require.Equal(t, quantity, openOrders[0].Quantity)
	require.Equal(t, status, openOrders[0].Status)
	require.Equal(t, filledAvgQuantity, openOrders[0].FilledAvgPrice)
	require.Equal(t, remainQuantity, openOrders[0].RemainQuantity)
	require.Equal(t, timestamp, openOrders[0].Timestamp)

	_, err = mockCli.Backend().QueryOpenOrders(addr[1:], product, side, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryOpenOrders(addr, "", side, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryOpenOrders(addr, product, "BUY&&SELL", start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryOpenOrders(addr, product, side, end+1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryOpenOrders(addr, product, side, -1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryOpenOrders(addr, product, side, start, -1, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryOpenOrders(addr, product, side, start, end, -1, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryOpenOrders(addr, product, side, start, end, page, -1)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.OpenOrdersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryOpenOrders(addr, product, side, start, end, page, perPage)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.OpenOrdersPath, cmn.HexBytes(queryBytes)).Return(expectedRet[1:], nil)
	_, err = mockCli.Backend().QueryOpenOrders(addr, product, side, start, end, page, perPage)
	require.Error(t, err)
}

func TestBackendClient_QueryClosedOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := sdk.NewClientConfig("testURL", "testChain", sdk.BroadcastBlock, "0.01okt", 200000)
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewBackendClient(mockCli.MockBaseClient))

	txHash, orderID, side := "default tx hash", "ID0000000000-1", "BUY"
	price, quantity, filledAvgQuantity, remainQuantity := "1024.1024", "2048.2048", "1024.1024", "4096.4096"
	status, timestamp := int64(0), time.Now().Unix()*1000
	start, end, page, perPage := 0, 0, 1, 30

	expectedRet := mockCli.BuildBackendOrdersResultBytes(txHash, orderID, addr, product, side, price, quantity, filledAvgQuantity,
		remainQuantity, status, timestamp)
	expectedCdc := mockCli.GetCodec()

	queryParams := params.NewQueryOrderListParams(addr, product, side, page, perPage, int64(start), int64(end), false)
	queryBytes := expectedCdc.MustMarshalJSON(queryParams)

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(3)
	mockCli.EXPECT().Query(types.ClosedOrdersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil)

	openOrders, err := mockCli.Backend().QueryClosedOrders(addr, product, side, start, end, page, perPage)
	require.NoError(t, err)
	require.Equal(t, txHash, openOrders[0].TxHash)
	require.Equal(t, orderID, openOrders[0].OrderID)
	require.Equal(t, addr, openOrders[0].Sender)
	require.Equal(t, product, openOrders[0].Product)
	require.Equal(t, side, openOrders[0].Side)
	require.Equal(t, price, openOrders[0].Price)
	require.Equal(t, quantity, openOrders[0].Quantity)
	require.Equal(t, status, openOrders[0].Status)
	require.Equal(t, filledAvgQuantity, openOrders[0].FilledAvgPrice)
	require.Equal(t, remainQuantity, openOrders[0].RemainQuantity)
	require.Equal(t, timestamp, openOrders[0].Timestamp)

	_, err = mockCli.Backend().QueryClosedOrders(addr[1:], product, side, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryClosedOrders(addr, "", side, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryClosedOrders(addr, product, "BUY&&SELL", start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryClosedOrders(addr, product, side, end+1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryClosedOrders(addr, product, side, -1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryClosedOrders(addr, product, side, start, -1, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryClosedOrders(addr, product, side, start, end, -1, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryClosedOrders(addr, product, side, start, end, page, -1)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.ClosedOrdersPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryClosedOrders(addr, product, side, start, end, page, perPage)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.ClosedOrdersPath, cmn.HexBytes(queryBytes)).Return(expectedRet[1:], nil)
	_, err = mockCli.Backend().QueryClosedOrders(addr, product, side, start, end, page, perPage)
	require.Error(t, err)
}

func TestBackendClient_QueryRecentTxRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := sdk.NewClientConfig("testURL", "testChain", sdk.BroadcastBlock, "0.01okt", 200000)
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewBackendClient(mockCli.MockBaseClient))

	timestamp, blockHeight := time.Now().Unix()*1000, int64(1024)
	price, quantity := 1024.1024, 2048.2048
	start, end, page, perPage := 0, 0, 1, 30

	expectedRet := mockCli.BuildBackendMatchResultBytes(timestamp, blockHeight, product, price, quantity)
	expectedCdc := mockCli.GetCodec()

	queryParams := params.NewQueryMatchParams(product, int64(start), int64(end), page, perPage)
	queryBytes := expectedCdc.MustMarshalJSON(queryParams)

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(3)
	mockCli.EXPECT().Query(types.RecentTxRecordPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil)

	txRecord, err := mockCli.Backend().QueryRecentTxRecord(product, start, end, page, perPage)
	require.NoError(t, err)
	require.Equal(t, timestamp, txRecord[0].Timestamp)
	require.Equal(t, product, txRecord[0].Product)
	require.Equal(t, price, txRecord[0].Price)
	require.Equal(t, quantity, txRecord[0].Quantity)
	require.Equal(t, blockHeight, txRecord[0].BlockHeight)

	_, err = mockCli.Backend().QueryRecentTxRecord("", start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryRecentTxRecord(product, end+1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryRecentTxRecord(product, -1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryRecentTxRecord(product, start, -1, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryRecentTxRecord(product, start, end, -1, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryRecentTxRecord(product, start, end, page, -1)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.RecentTxRecordPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryRecentTxRecord(product, start, end, page, perPage)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.RecentTxRecordPath, cmn.HexBytes(queryBytes)).Return(expectedRet[1:], nil)
	_, err = mockCli.Backend().QueryRecentTxRecord(product, start, end, page, perPage)
	require.Error(t, err)
}

func TestBackendClient_QueryTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := sdk.NewClientConfig("testURL", "testChain", sdk.BroadcastBlock, "0.01okt", 200000)
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewBackendClient(mockCli.MockBaseClient))

	txHash, quantity, fee := "default tx hash", "1024.1024", "0.00000000okt"
	txType, side, timestamp := 2, int64(1), time.Now().Unix()*1000
	start, end, page, perPage := 0, 0, 1, 30

	expectedRet := mockCli.BuildBackendTransactionsResultBytes(txHash, addr, product, quantity, fee, int64(txType), side, timestamp)
	expectedCdc := mockCli.GetCodec()

	queryParams := params.NewQueryTxListParams(addr, int64(txType), int64(start), int64(end), page, perPage)
	queryBytes := expectedCdc.MustMarshalJSON(queryParams)

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(3)
	mockCli.EXPECT().Query(types.TransactionsPath, cmn.HexBytes(queryBytes)).Return(expectedRet, nil)

	txs, err := mockCli.Backend().QueryTransactions(addr, txType, start, end, page, perPage)
	require.NoError(t, err)
	require.Equal(t, txHash, txs[0].TxHash)
	require.Equal(t, int64(txType), txs[0].Type)
	require.Equal(t, addr, txs[0].Address)
	require.Equal(t, product, txs[0].Symbol)
	require.Equal(t, side, txs[0].Side)
	require.Equal(t, quantity, txs[0].Quantity)
	require.Equal(t, fee, txs[0].Fee)
	require.Equal(t, timestamp, txs[0].Timestamp)

	_, err = mockCli.Backend().QueryTransactions(addr[1:], txType, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTransactions(addr, -1, start, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTransactions(addr, txType, end+1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTransactions(addr, txType, -1, end, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTransactions(addr, txType, start, -1, page, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTransactions(addr, txType, start, end, -1, perPage)
	require.Error(t, err)

	_, err = mockCli.Backend().QueryTransactions(addr, txType, start, end, page, -1)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.TransactionsPath, cmn.HexBytes(queryBytes)).Return(expectedRet, errors.New("default error"))
	_, err = mockCli.Backend().QueryTransactions(addr, txType, start, end, page, perPage)
	require.Error(t, err)

	mockCli.EXPECT().Query(types.TransactionsPath, cmn.HexBytes(queryBytes)).Return(expectedRet[1:], nil)
	_, err = mockCli.Backend().QueryTransactions(addr, txType, start, end, page, perPage)
	require.Error(t, err)
}
