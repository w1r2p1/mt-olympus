package engine

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "matcher/proto"
)

// Process an order and return the trades generated before adding the remaining amount to the market
func (ob *OrderBook) Process(order Order) (pb.OrderConf, []Trade) {
	orderConf := pb.OrderConf{
		UserId:               order.UserId,
		OrderId:              order.OrderId,
		Amount:               order.Amount,
		Price:                order.Price,
		Side:                 order.Side,
		Type:                 order.Type,
		Message:              "Confirmed",
		CreatedAt:            ptypes.TimestampNow(),
	}

	if order.Side == pb.Side_BUY {
		return orderConf, ob.processLimitBuy(order)
	}

	return orderConf, ob.processLimitSell(order)
}

// Process a limit buy order
func (ob *OrderBook) processLimitBuy(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	numSells := len(ob.SellOrders)

	// Check if we have at least one matching order
	if numSells > 0 && ob.SellOrders[numSells-1].Price <= order.Price {
		// Traverse all orders that match
		for i := numSells - 1; i >= 0; i-- {
			sellOrder := ob.SellOrders[i]

			if sellOrder.Price > order.Price {
				break
			}

			// Fill the entire order
			if sellOrder.Amount >= order.Amount {
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    sellOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   sellOrder.OrderId,
					Amount:     order.Amount,
					Price:      sellOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}

				//trade := Trade{order.ID, sellOrder.ID, order.Amount, sellOrder.Price}
				trades = append(trades, trade)

				sellOrder.Amount -= order.Amount
				if sellOrder.Amount == 0 {
					ob.removeSellOrder(i)
				}

				return trades
			}
			// Fill a partial order and continue
			if sellOrder.Amount < order.Amount {
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    sellOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   sellOrder.OrderId,
					Amount:     sellOrder.Amount,
					Price:      sellOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}

				//trade := Trade{order.ID, sellOrder.ID, sellOrder.Amount, sellOrder.Price}
				trades = append(trades, trade)

				order.Amount -= sellOrder.Amount
				ob.removeSellOrder(i)
			}
		}
	}
	// finally add the remaining order to the list
	ob.addBuyOrder(order)

	return trades
}

// Process a limit sell order
func (ob *OrderBook) processLimitSell(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	numBuys := len(ob.BuyOrders)

	// Check if we have at least one matching order
	if numBuys > 0 && ob.BuyOrders[numBuys-1].Price >= order.Price {
		// Traverse all orders that match
		for i := numBuys - 1; i >= 0; i-- {
			buyOrder := ob.BuyOrders[i]

			if buyOrder.Price < order.Price {
				break
			}

			// Fill the entire order
			if buyOrder.Amount >= order.Amount {
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    buyOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   buyOrder.OrderId,
					Amount:     order.Amount,
					Price:      buyOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}

				//trades = append(trades, Trade{order.ID, buyOrder.ID, order.Amount, buyOrder.Price})
				trades = append(trades, trade)
				buyOrder.Amount -= order.Amount
				if buyOrder.Amount == 0 {
					ob.removeBuyOrder(i)
				}

				return trades
			}

			// Fill a partial order and continue
			if buyOrder.Amount < order.Amount {
				//trades = append(trades, Trade{order.ID, buyOrder.ID, buyOrder.Amount, buyOrder.Price})
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    buyOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   buyOrder.OrderId,
					Amount:     buyOrder.Amount,
					Price:      buyOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}

				trades = append(trades, trade)

				order.Amount -= buyOrder.Amount
				ob.removeBuyOrder(i)
			}
		}
	}

	// Finally add the remaining order to the list
	ob.addSellOrder(order)

	return trades
}
