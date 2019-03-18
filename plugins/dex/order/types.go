package order

import (
	"fmt"

	"github.com/binance-chain/node/common/types"
	me "github.com/binance-chain/node/plugins/dex/matcheng"
)

// The types here are shared between order and pub package

type ChangeType uint8

const (
	Ack ChangeType = iota
	Canceled
	Expired
	IocNoFill
	PartialFill
	FullyFill
	FailedBlocking
	FailedMatching
)

// True for should not remove order in these status from OrderInfoForPub
// False for remove
func (tpe ChangeType) IsOpen() bool {
	// FailedBlocking tx doesn't effect OrderInfoForPub, should not be put into closedToPublish
	return tpe == Ack ||
		tpe == PartialFill ||
		tpe == FailedBlocking
}

func (tpe ChangeType) String() string {
	switch tpe {
	case Ack:
		return "Ack"
	case Canceled:
		return "Canceled"
	case Expired:
		return "Expired"
	case IocNoFill:
		return "IocNoFill"
	case PartialFill:
		return "PartialFill"
	case FullyFill:
		return "FullyFill"
	case FailedBlocking:
		return "FailedBlocking"
	case FailedMatching:
		return "FailedMatching"
	default:
		return "Unknown"
	}
}

type ExecutionType uint8

const (
	NEW ExecutionType = iota
)

func (this ExecutionType) String() string {
	switch this {
	case NEW:
		return "NEW"
	default:
		return "Unknown"
	}
}

type OrderChange struct {
	Id             string
	Tpe            ChangeType
	MsgForFailedTx interface{} // pointer to NewOrderMsg or CancelOrderMsg
}

func (oc OrderChange) String() string {
	return fmt.Sprintf("id: %s, tpe: %s", oc.Id, oc.Tpe.String())
}

func (oc OrderChange) failedBlockingMsg() *OrderInfo {
	switch msg := oc.MsgForFailedTx.(type) {
	case NewOrderMsg:
		return &OrderInfo{
			NewOrderMsg: msg,
		}
	case CancelOrderMsg:
		return &OrderInfo{
			NewOrderMsg: NewOrderMsg{Sender: msg.Sender, Id: msg.RefId, Symbol: msg.Symbol},
		}
	default:
		return nil
	}
}

func (oc OrderChange) ResolveOrderInfo(orderInfos OrderInfoForPublish) *OrderInfo {
	switch oc.Tpe {
	case FailedBlocking:
		return oc.failedBlockingMsg()
	default:
		return orderInfos[oc.Id]
	}
}

// provide an easy way to retrieve order related static fields during generate executed order status
type OrderInfoForPublish map[string]*OrderInfo
type OrderChanges []OrderChange // clean after publish each block's EndBlock and before next block's BeginBlock

type ChangedPriceLevelsMap map[string]ChangedPriceLevelsPerSymbol

type ChangedPriceLevelsPerSymbol struct {
	Buys  map[int64]int64
	Sells map[int64]int64
}

type TradeHolder struct {
	OId    string
	Trade  *me.Trade
	Symbol string
}

func (fh TradeHolder) String() string {
	return fmt.Sprintf("oid: %s, bid: %s, sid: %s", fh.OId, fh.Trade.Bid, fh.Trade.Sid)
}

type ExpireHolder struct {
	OrderId string
}

type FeeHolder map[string]*types.Fee
