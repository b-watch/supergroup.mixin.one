package models

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/lib/pq"
	"github.com/objcoding/wxpay"
)

const order_DDL = `
CREATE TABLE IF NOT EXISTS orders (
	order_id         VARCHAR(36) PRIMARY KEY CHECK (order_id ~* '^[0-9a-f-]{36,36}$'),
	trace_id         BIGSERIAL,
	user_id          VARCHAR(36) NOT NULL CHECK (user_id ~* '^[0-9a-f-]{36,36}$'),
	prepay_id        VARCHAR(36) DEFAULT '',
	state            VARCHAR(32) NOT NULL,
	asset_id         VARCHAR(36) NOT NULL,
	amount           VARCHAR(128) NOT NULL,
	pay_method          VARCHAR(32) NOT NULL,
	transaction_id   VARCHAR(32) DEFAULT '',
	created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	paid_at          TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS order_created_paidx ON orders(created_at, paid_at);
`

type Order struct {
	OrderId       string      `json:"order_id"`
	UserId        string      `json:"user_id"`
	TraceId       int64       `json:"trace_id"`
	PrepayId      string      `json:"prepay_id"`
	State         string      `json:"state"`
	AssetId       string      `json:"asset_id"`
	Amount        string      `json:"amount"`
	PayMethod     string      `json:"pay_method"`
	TransactionId string      `json:"transaction_id"`
	CreatedAt     time.Time   `json:"created_at"`
	PaidAt        pq.NullTime `json:"paid_at"`
}

const WX_TN_PREFIX = "tn-"

var orderColumns = []string{"order_id", "user_id", "trace_id", "prepay_id", "state", "asset_id", "amount", "pay_method", "transaction_id", "created_at", "paid_at"}

func (o *Order) values() []interface{} {
	return []interface{}{o.OrderId, o.UserId, o.TraceId, o.PrepayId, o.State, o.AssetId, o.Amount, o.PayMethod, o.TransactionId, o.CreatedAt, o.PaidAt}
}

func orderFromRow(row durable.Row) (*Order, error) {
	var o Order
	err := row.Scan(&o.OrderId, &o.UserId, &o.TraceId, &o.PrepayId, &o.State, &o.AssetId, &o.Amount, &o.PayMethod, &o.TransactionId, &o.CreatedAt, &o.PaidAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &o, err
}

func GetNotPaidOrders(ctx context.Context, limit int64) ([]*Order, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM orders WHERE
		state='NOTPAID'
			AND created_at > NOW() - INTERVAL '%d minute'
			AND user_id NOT IN
				(SELECT user_id FROM users WHERE state='paid')
		ORDER BY created_at`, strings.Join(orderColumns, ","), limit)
	// query := fmt.Sprintf("SELECT %s FROM orders WHERE state='NOTPAID' AND created_at > NOW() - INTERVAL '30 minute' ORDER BY created_at", strings.Join(orderColumns, ","))
	rows, err := session.Database(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		order, err := orderFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func CreateWechatOrder(ctx context.Context, userId, amount, wxOpenId string) (*Order, wxpay.Params, wxpay.Params, error) {
	var order *Order
	var err error
	// create an order
	order, err = createOrder(ctx, userId, "", amount, "PENDING", PayMethodWechat)
	if err != nil {
		return nil, nil, nil, session.TransactionError(ctx, err)
	}

	order, err = GetOrder(ctx, order.OrderId)
	if err != nil {
		return nil, nil, nil, err
	}

	// create wx payment request
	var wxp wxpay.Params
	var jswxp wxpay.Params
	client := CreateWxClient()
	wxp, err = CreateWxPayment(client, order.TraceId, order.Amount, wxOpenId)
	if err != nil {
		return nil, nil, nil, err
	}

	// sign params for jsapi
	jswxp = GetPayJsParams(client, wxp)

	// update record
	order.State = "NOTPAID"
	query := "UPDATE orders SET state=$1, prepay_id=$2 WHERE order_id=$3"
	_, err = session.Database(ctx).ExecContext(ctx, query, order.State, wxp["prepay_id"], order.OrderId)
	if err != nil {
		return nil, nil, nil, err
	}

	return order, wxp, jswxp, nil
}

func CreateMixinOrder(ctx context.Context, userId, assetId, amount string) (*Order, error) {
	var order *Order
	var err error
	// create an order
	order, err = createOrder(ctx, userId, assetId, amount, "NOTPAID", PayMethodMixin)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return order, nil
}

func createOrder(ctx context.Context, userId, assetId, amount, state, method string) (*Order, error) {
	order := &Order{
		OrderId:       bot.UuidNewV4().String(),
		UserId:        userId,
		TraceId:       0,
		PrepayId:      "",
		State:         state,
		AssetId:       assetId,
		Amount:        amount,
		PayMethod:     method,
		TransactionId: "",
	}

	// create an order
	var err error
	query := "INSERT INTO orders (order_id, user_id, prepay_id, state, asset_id, amount, pay_method) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = session.Database(ctx).ExecContext(ctx, query,
		order.OrderId, order.UserId, order.PrepayId, order.State, order.AssetId, order.Amount, order.PayMethod)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return order, nil
}

func MarkOrderAsPaidByTraceId(ctx context.Context, traceId int64, transactionId string) (*Order, error) {
	var order *Order
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		order, err = getOrderByTraceId(ctx, tx, traceId)
		if err != nil || order == nil {
			return err
		}
		query := "UPDATE orders SET state='PAID', transaction_id=$1, paid_at=$2 WHERE order_id=$3"
		_, err = tx.ExecContext(ctx, query, transactionId, time.Now(), order.OrderId)
		if err != nil {
			return err
		}
		user, err := findUserById(ctx, tx, order.UserId)
		if err != nil {
			return err
		}
		return user.paymentInTx(ctx, tx, PayMethodWechat)
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}

	plugin.Trigger(plugin.EventTypeOrderPaid, *order)

	return order, nil
}

func MarkOrderAsPaidByOrderId(ctx context.Context, orderId string) (*Order, error) {
	var order *Order
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		order, err = getOrderByOrderId(ctx, tx, orderId)
		if err != nil || order == nil {
			return err
		}
		query := "UPDATE orders SET state='PAID', paid_at=$1, pay_method=$2 WHERE order_id=$3"
		_, err = tx.ExecContext(ctx, query, time.Now(), PayMethodMixin, order.OrderId)
		if err != nil {
			return err
		}
		user, err := findUserById(ctx, tx, order.UserId)
		if err != nil {
			return err
		}
		return user.paymentInTx(ctx, tx, PayMethodMixin)
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}

	plugin.Trigger(plugin.EventTypeOrderPaid, *order)

	return order, nil
}

func getOrderByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (*Order, error) {
	query := fmt.Sprintf("SELECT %s FROM orders WHERE order_id=$1 ORDER BY created_at LIMIT 1", strings.Join(orderColumns, ","))
	row := tx.QueryRowContext(ctx, query, orderId)
	order, err := orderFromRow(row)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return order, nil
}

func getOrderByTraceId(ctx context.Context, tx *sql.Tx, traceId int64) (*Order, error) {
	query := fmt.Sprintf("SELECT %s FROM orders WHERE trace_id=$1 ORDER BY created_at LIMIT 1", strings.Join(orderColumns, ","))
	row := tx.QueryRowContext(ctx, query, traceId)
	order, err := orderFromRow(row)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return order, nil
}

func GetOrder(ctx context.Context, orderId string) (*Order, error) {
	var order *Order
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_order, err := getOrderByOrderId(ctx, tx, orderId)
		order = _order
		return err
	})
	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return nil, sessionErr
		}
		return nil, session.TransactionError(ctx, err)
	}
	return order, nil
}

func CreateWxClient() *wxpay.Client {
	cfg := config.AppConfig
	account := wxpay.NewAccount(cfg.Wechat.AppId, cfg.Wechat.MchId, cfg.Wechat.MchKey, false)
	client := wxpay.NewClient(account)
	account.SetCertData("./cert_test.p12")
	client.SetAccount(account)
	client.SetHttpConnectTimeoutMs(2000)
	client.SetHttpReadTimeoutMs(1000)
	client.SetSignType(wxpay.MD5)
	return client
}

func CreateWxPayment(client *wxpay.Client, traceId int64, amount, wxOpenId string) (wxpay.Params, error) {
	fs, _ := strconv.ParseFloat(amount, 32)
	tradeNo := WX_TN_PREFIX + strconv.FormatInt(traceId, 10)
	params := make(wxpay.Params)
	params.
		SetString("out_trade_no", tradeNo).
		SetInt64("total_fee", int64(math.Ceil(fs*100))).
		// I don't know what's the meaning of the IP
		SetString("spbill_create_ip", "123.12.12.123").
		// I don't have the permission of notify url.
		// so I pull to get order state
		// @TODO need to implement an method to handle it.
		SetString("notify_url", config.AppConfig.Wechat.NotifyUrl).
		// drop some shits here.
		SetString("body", "Mixin-PayToJoin").
		// only support jsapi trade type for now. No permission for "H5" trade type.
		SetString("trade_type", "JSAPI").
		SetString("openid", wxOpenId)

	p, err := client.UnifiedOrder(params)

	return p, err
}

func FetchWxPayment(client *wxpay.Client, traceId int64) (wxpay.Params, error) {
	tradeNo := WX_TN_PREFIX + strconv.FormatInt(traceId, 10)
	params := make(wxpay.Params)
	params.SetString("out_trade_no", tradeNo)
	return client.OrderQuery(params)
}

func GetPayJsParams(client *wxpay.Client, params wxpay.Params) wxpay.Params {
	// for JSAPI payment, we have to sign again for slight different params.
	// be careful about the stupid fields spelling, WeChat's API design is horrible.
	payParams := make(wxpay.Params)
	payParams.SetString("appId", params["appid"])
	payParams.SetString("timeStamp", strconv.FormatInt(time.Now().Unix(), 10))
	payParams.SetString("nonceStr", strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	// for JSAPI payment, use the prepay_id which get from UnifiedOrder() in this stupid form
	payParams.SetString("package", "prepay_id="+params["prepay_id"])
	payParams.SetString("signType", wxpay.MD5)
	// No 'sign', please use 'paySign'. The stupid WeChat online sign tool say NO? Leave that!.
	payParams.SetString("paySign", client.Sign(payParams))
	return payParams
}
