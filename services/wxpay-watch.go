package services

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/objcoding/wxpay"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
)

type WxPaymentService struct{}

func StartWxPaymentWatch(name string, db *durable.Database) {
	context := session.WithDatabase(context.Background(), db)
	client := models.CreateWxClient()
	ctx := session.WithLogger(context, durable.BuildLogger())
	var orders []*models.Order
	var err error
	var params wxpay.Params
	for true {
		// check orders with state "NOTPAID" in every 7 seconds
		// the window is 15 min
		// @TODO
		// 1. do not check the orders which of owners who have paid.
		// 2. handle notify_url for better performance.
		orders, err = models.GetNotPaidOrders(ctx)
		if err != nil {
			time.Sleep(time.Duration(20) * time.Second)
			log.Printf("Error in StartWxPaymentWatch's Loop: %v\n", err)
			continue
		}
		if len(orders) != 0 {
			log.Printf("Handle %v orders in StartWxPaymentWatch's Loop\n", len(orders))
		}
		for _, order := range orders {
			params, err = models.FetchWxPayment(client, order.TraceId)
			if err != nil {
				time.Sleep(time.Duration(1) * time.Second)
				continue
			}
			if params["result_code"] == "SUCCESS" && params["trade_state"] == "SUCCESS" {
				tn := params["out_trade_no"]
				transactionId := params["transaction_id"]
				if strings.HasPrefix(tn, models.WX_TN_PREFIX) {
					if tnId, err := strconv.ParseInt(tn[3:], 10, 64); err == nil {
						models.UpdateOrderStateByTraceId(ctx, tnId, "PAID", transactionId)
						if user, err := models.FindUser(ctx, order.UserId); err == nil {
							user.Payment(ctx)
						}
					}
				}
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
		time.Sleep(time.Duration(7) * time.Second)
	}
}
