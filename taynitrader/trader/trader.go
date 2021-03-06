package trader

import (
	"fmt"

	"github.com/lagarciag/tayni/kredis"
	"github.com/lagarciag/tayni/twitter"
	"github.com/looplab/fsm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	StartState       = "StartState"
	IdleState        = "IdleState"
	TradingState     = "TradingState"
	TestTradingState = "TestTradingState"
	HoldState        = "HoldState"
	TestHoldState    = "TestHoldState"
	ShutdownState    = "ShutdownState"

	// --------------
	// Buy States
	// --------------

	Minute1BuyState   = "Minute1BuyState"
	Minute120BuyState = "Minute120BuyState"
	Minute60BuyState  = "Minute60BuyState"
	Minute30BuyState  = "Minute30BuySate"

	DoBuyState     = "DoBuyState"
	TestDoBuyState = "TestDoBuyState"

	// -------------
	//   Sell Sates
	// -------------

	Minute1SellState   = "Minute1SellState"
	Minute120SellState = "Minute120SellState"
	Minute60SellState  = "Minute60SellState"
	Minute30SellState  = "Minute30SellState"

	DoSellState     = "DoSellState"
	TestDoSellState = "TestDoSellState"
)

const (
	StartEvent    = "startEvent"
	StopEvent     = "stopEvent"
	ShutdownEvent = "shutdownEvent"
	TradeEvent    = "TradeEvent"
	HoldEvent     = "HoldEvent"
	Test1Event    = "Test2Event"
	// --------------
	// Buy Events
	// --------------

	Minute1BuyEvent   = "Minute1BuyEvent"
	Minute120BuyEvent = "Minute120BuyEvent"
	Minute60BuyEvent  = "Minute60BuyEvent"
	Minute30BuyEvent  = "Minute30BuySate"

	NotMinute1BuyEvent   = "NotMinute1BuyEvent"
	NotMinute120BuyEvent = "NotMinute120BuyEvent"
	NotMinute60BuyEvent  = "NotMinute60BuyEvent"
	NotMinute30BuyEvent  = "NotMinute30BuySate"

	DoBuyEvent            = "DoBuyEvent"
	TestDoBuyEvent        = "TestDoBuyEvent"
	DoSellEvent           = "DoSellEvent"
	TestDoSellEvent       = "TestDoSellEvent"
	BuyCompleteEvent      = "BuyCompleteEvent"
	TestBuyCompleteEvent  = "TestBuyCompleteEvent"
	SellCompleteEvent     = "SellCompleteEvent"
	TestSellCompleteEvent = "TestSellCompleteEvent"

	// -------------
	//   Sell Sates
	// -------------

	Minute1SellEvent   = "Minute1SellEvent"
	Minute120SellEvent = "Minute120SellEvent"
	Minute60SellEvent  = "Minute60SellEvent"
	Minute30SellEvent  = "Minute30SellEvent"

	NotMinute1SellEvent   = "NotMinute1SellEvent"
	NotMinute120SellEvent = "NotMinute120SellEvent"
	NotMinute60SellEvent  = "NotMinute60SellEvent"
	NotMinute30SellEvent  = "NotMinute30SellEvent"
)

type TradeFsm struct {
	tc           *twitter.TwitterClient
	kr           *kredis.Kredis
	To           string
	FSM          *fsm.FSM
	pairID       string
	holdingFunds bool

	BuyStates     []string
	SellStates    []string
	ControlStates []string
	TradingStates []string
	AllStates     []string

	BuyEvents     []string
	SellEvents    []string
	NotBuyEvents  []string
	NotSellEvents []string
	ControlEvents []string
	TradingEvents []string
	AllEvents     []string

	fsmBuyEventsDescriptors     []fsm.EventDesc
	fsmSellEventsDescriptors    []fsm.EventDesc
	fsmNotBuyEventsDescriptors  []fsm.EventDesc
	fsmNotSellEventsDescriptors []fsm.EventDesc

	// ------------
	// Events
	// ------------
	startEvent    fsm.EventDesc
	stopEvent     fsm.EventDesc
	shutdownEvent fsm.EventDesc
	tradeEvent    fsm.EventDesc
	holdEvent     fsm.EventDesc

	minute1BuyEvent   fsm.EventDesc
	minute120BuyEvent fsm.EventDesc
	minute60BuyEvent  fsm.EventDesc
	minute30BuyEvent  fsm.EventDesc

	notMinute1BuyEvent   fsm.EventDesc
	notMinute120BuyEvent fsm.EventDesc
	notMinute60BuyEvent  fsm.EventDesc
	notMinute30BuyEvent  fsm.EventDesc

	minute1SellEvent   fsm.EventDesc
	minute120SellEvent fsm.EventDesc
	minute60SellEvent  fsm.EventDesc
	minute30SellEvent  fsm.EventDesc

	notMinute1SellEvent   fsm.EventDesc
	notMinute120SellEvent fsm.EventDesc
	notMinute60SellEvent  fsm.EventDesc
	notMinute30SellEvent  fsm.EventDesc

	doBuyEvent        fsm.EventDesc
	doSellEvent       fsm.EventDesc
	buyCompleteEvent  fsm.EventDesc
	sellCompleteEvent fsm.EventDesc

	testDoBuyEvent        fsm.EventDesc
	testDoSellEvent       fsm.EventDesc
	testBuyCompleteEvent  fsm.EventDesc
	testSellCompleteEvent fsm.EventDesc

	// ------------------
	// Events List
	// ------------------
	eventsList fsm.Events

	// ------------------
	// Fsm Callbacks
	callbacks fsm.Callbacks

	ChanStartEvent    chan bool
	ChanStopEvent     chan bool
	ChanShutdownEvent chan bool
	ChanTradeEvent    chan bool
	ChanHoldEvent     chan bool

	// --------------
	// Buy Events
	// --------------
	ChanMinute1BuyEvent   chan bool
	ChanMinute120BuyEvent chan bool
	ChanMinute60BuyEvent  chan bool
	ChanMinute30BuyEvent  chan bool

	ChanNotMinute1BuyEvent chan bool
	//ChanNotMinute120BuyEvent chan bool
	//ChanNotMinute60BuyEvent  chan bool
	//ChanNotMinute30BuyEvent  chan bool

	ChanMinute1SellEvent   chan bool
	ChanMinute120SellEvent chan bool
	ChanMinute60SellEvent  chan bool
	ChanMinute30SellEvent  chan bool

	//ChanNotMinute1SellEvent   chan bool
	//ChanNotMinute120SellEvent chan bool
	//ChanNotMinute60SellEvent  chan bool
	//ChanNotMinute30SellEvent  chan bool

	ChanDoBuyEvent        chan bool
	ChanDoSellEvent       chan bool
	ChanBuyCompleteEvent  chan bool
	ChanSellCompleteEvent chan bool

	testChanDoBuyEvent        chan bool
	testChanDoSellEvent       chan bool
	testChanBuyCompleteEvent  chan bool
	testChanSellCompleteEvent chan bool

	ChanMap map[string]chan bool
}

func NewTradeFsm(pairID string) *TradeFsm {
	log.Info("Creating new trading fsm for pair: ", pairID)

	tFsm := &TradeFsm{}

	//---------------------
	//Event structures
	//---------------------

	tFsm.BuyStates = []string{Minute120BuyState, Minute60BuyState, Minute30BuyState}
	tFsm.SellStates = []string{Minute120SellState, Minute60SellState, Minute30SellState}
	tFsm.ControlStates = []string{StartState, IdleState, TradingState, HoldState, ShutdownState}
	tFsm.TradingStates = []string{DoBuyState, DoSellState}

	tFsm.AllStates = make([]string, len(tFsm.BuyStates))
	copy(tFsm.AllStates, tFsm.BuyStates)
	tFsm.AllStates = append(tFsm.AllStates, tFsm.SellStates...)
	tFsm.AllStates = append(tFsm.AllStates, tFsm.ControlStates...)
	tFsm.AllStates = append(tFsm.AllStates, tFsm.TradingStates...)

	tFsm.BuyEvents = []string{Minute120BuyEvent, Minute60BuyEvent, Minute30BuyEvent}
	tFsm.SellEvents = []string{Minute120SellEvent, Minute60SellEvent, Minute30SellEvent}

	tFsm.NotBuyEvents = []string{NotMinute120BuyEvent, NotMinute60BuyEvent, NotMinute30BuyEvent}
	tFsm.NotSellEvents = []string{NotMinute120SellEvent, NotMinute60SellEvent, NotMinute30SellEvent}

	tFsm.ControlEvents = []string{StartEvent, StopEvent, TradeEvent, HoldEvent, ShutdownEvent}
	tFsm.TradingEvents = []string{DoBuyEvent, DoSellEvent}

	tFsm.AllEvents = make([]string, len(tFsm.BuyEvents))
	copy(tFsm.AllEvents, tFsm.BuyEvents)
	tFsm.AllEvents = append(tFsm.AllEvents, tFsm.SellEvents...)
	tFsm.AllEvents = append(tFsm.AllEvents, tFsm.ControlEvents...)
	tFsm.AllEvents = append(tFsm.AllEvents, tFsm.TradingEvents...)

	// ----------------------
	// Kredis configuration
	// ----------------------
	tFsm.kr = kredis.NewKredis(1000000)
	tFsm.kr.Start()

	// ----------------------
	// Twitter configuration
	config := twitter.Config{}
	vTwitterConfig := viper.Get("twitter").(map[string]interface{})
	config.Twit = vTwitterConfig["twit"].(bool)
	config.ConsumerKey = vTwitterConfig["consumer_key"].(string)
	config.ConsumerSecret = vTwitterConfig["consumer_secret"].(string)
	config.AccessToken = vTwitterConfig["access_token"].(string)
	config.AccessTokenSecret = vTwitterConfig["access_token_secret"].(string)

	if config.ConsumerKey == "" {
		log.Fatal("bad consumerkey")
	}

	tFsm.tc = twitter.NewTwitterClient(config)
	tFsm.pairID = pairID

	// ------------
	// Events
	// ------------
	startEvent := fsm.EventDesc{Name: StartEvent, Src: []string{StartState}, Dst: IdleState}

	stopEvent := fsm.EventDesc{Name: StopEvent,
		Src: []string{
			TradingState,
			HoldState,
		}, Dst: IdleState}

	stopEvent.Src = append(stopEvent.Src, tFsm.BuyStates...)
	stopEvent.Src = append(stopEvent.Src, tFsm.SellStates...)

	tradeEvent := fsm.EventDesc{Name: TradeEvent, Src: []string{IdleState}, Dst: TradingState}
	holdEvent := fsm.EventDesc{Name: HoldEvent, Src: []string{IdleState}, Dst: HoldState}

	// ----------------------
	// Buying related events
	// ----------------------

	tFsm.fsmBuyEventsDescriptors = make([]fsm.EventDesc, len(tFsm.BuyEvents))

	for ID, eventName := range tFsm.BuyEvents {

		if ID == 0 {
			tFsm.fsmBuyEventsDescriptors[ID] = fsm.EventDesc{Name: eventName,
				Src: []string{TradingState},
				Dst: tFsm.BuyStates[ID]}
		} else {
			tFsm.fsmBuyEventsDescriptors[ID] = fsm.EventDesc{Name: eventName,
				Src: []string{tFsm.BuyStates[ID-1]},
				Dst: tFsm.BuyStates[ID]}
		}

	}

	// -------------
	// Sell Events
	// -------------
	tFsm.fsmSellEventsDescriptors = make([]fsm.EventDesc, len(tFsm.SellEvents))
	for ID, eventName := range tFsm.SellEvents {

		if ID == 0 {
			tFsm.fsmSellEventsDescriptors[ID] = fsm.EventDesc{Name: eventName,
				Src: []string{HoldState},
				Dst: tFsm.SellStates[ID]}
		} else {
			tFsm.fsmSellEventsDescriptors[ID] = fsm.EventDesc{Name: eventName,
				Src: []string{tFsm.SellStates[ID-1]},
				Dst: tFsm.SellStates[ID]}
		}

	}

	tFsm.fsmNotBuyEventsDescriptors = make([]fsm.EventDesc, len(tFsm.NotBuyEvents))

	for ID, eventName := range tFsm.NotBuyEvents {

		var dstState string

		if ID == 0 {
			dstState = TradingState
		} else {
			dstState = tFsm.BuyStates[ID-1]
		}

		tFsm.fsmNotBuyEventsDescriptors[ID] = fsm.EventDesc{Name: eventName,
			Src: []string{},
			Dst: dstState}

		for i := len(tFsm.BuyEvents) - 1; i >= ID; i-- {
			tFsm.fsmNotBuyEventsDescriptors[ID].Src = append(tFsm.fsmNotBuyEventsDescriptors[ID].Src, tFsm.BuyStates[i])
		}
	}

	tFsm.fsmNotSellEventsDescriptors = make([]fsm.EventDesc, len(tFsm.NotSellEvents))

	for ID, eventName := range tFsm.NotSellEvents {

		var dstState string

		if ID == 0 {
			dstState = HoldState
		} else {
			dstState = tFsm.SellStates[ID-1]
		}

		tFsm.fsmNotSellEventsDescriptors[ID] = fsm.EventDesc{Name: eventName,
			Src: []string{},
			Dst: dstState}

		for i := len(tFsm.NotSellEvents) - 1; i >= ID; i-- {
			tFsm.fsmNotSellEventsDescriptors[ID].Src = append(tFsm.fsmNotSellEventsDescriptors[ID].Src, tFsm.SellStates[i])
		}
	}

	log.Debug("NotBuyEvents Autogen List: ", tFsm.fsmNotBuyEventsDescriptors)
	log.Debug("NotSellEvents Autogen List: ", tFsm.fsmNotSellEventsDescriptors)

	// ----------------------
	// Selling related events
	// ----------------------

	doBuyEvent := fsm.EventDesc{Name: DoBuyEvent,
		Src: []string{Minute30BuyState},
		Dst: DoBuyState}

	doSellEvent := fsm.EventDesc{Name: DoSellEvent,
		Src: []string{Minute30SellState},
		Dst: DoSellState}

	buyCompleteEvent := fsm.EventDesc{Name: BuyCompleteEvent,
		Src: []string{DoBuyState},
		Dst: HoldState}

	sellCompleteEvent := fsm.EventDesc{Name: SellCompleteEvent,
		Src: []string{DoSellState},
		Dst: TradingState}

	tFsm.startEvent = startEvent
	tFsm.stopEvent = stopEvent
	tFsm.tradeEvent = tradeEvent
	tFsm.holdEvent = holdEvent

	tFsm.doBuyEvent = doBuyEvent
	tFsm.doSellEvent = doSellEvent

	tFsm.buyCompleteEvent = buyCompleteEvent
	tFsm.sellCompleteEvent = sellCompleteEvent

	tFsm.shutdownEvent = fsm.EventDesc{Name: ShutdownEvent,
		Src: []string{StartState,
			IdleState},
		Dst: ShutdownState}

	// ----------------------
	// Events List Registry
	// ----------------------
	tFsm.eventsList = fsm.Events{
		tFsm.startEvent,
		tFsm.stopEvent,
		tFsm.shutdownEvent,
		tFsm.tradeEvent,
		tFsm.holdEvent,
		tFsm.doBuyEvent,
		tFsm.doSellEvent,

		tFsm.buyCompleteEvent,
		tFsm.sellCompleteEvent,
	}

	tFsm.eventsList = append(tFsm.eventsList, tFsm.fsmBuyEventsDescriptors...)
	tFsm.eventsList = append(tFsm.eventsList, tFsm.fsmSellEventsDescriptors...)
	tFsm.eventsList = append(tFsm.eventsList, tFsm.fsmNotBuyEventsDescriptors...)
	tFsm.eventsList = append(tFsm.eventsList, tFsm.fsmNotSellEventsDescriptors...)
	// -------------------
	// Callbacks registry
	// -------------------
	tFsm.callbacks = fsm.Callbacks{
		StartState:    tFsm.CallBackInGenericState,
		IdleState:     tFsm.CallBackInGenericState,
		ShutdownState: tFsm.CallBackInGenericState,
		TradingState:  tFsm.CallBackInGenericState,
		HoldState:     tFsm.CallBackInGenericState,

		Minute120BuyState: tFsm.CallBackInGenericState,
		Minute60BuyState:  tFsm.CallBackInGenericState,
		Minute30BuyState:  tFsm.CallBackInGenericState, //tFsm.CallBackInMinute30BuyState,

		Minute120SellState: tFsm.CallBackInGenericState,
		Minute60SellState:  tFsm.CallBackInGenericState,
		Minute30SellState:  tFsm.CallBackInGenericState, //tFsm.CallBackInMinute30SellState,

		DoBuyEvent:            tFsm.CallBackInDoBuyState,
		DoSellEvent:           tFsm.CallBackInDoSellState,
		BuyCompleteEvent:      tFsm.CallBackInBuyCompleteState,
		SellCompleteEvent:     tFsm.CallBackInSellCompleteState,
		TestSellCompleteEvent: tFsm.CallBackInTestSellCompleteState,
	}

	// ------------------
	// Event Channels
	// ------------------

	tFsm.ChanStartEvent = make(chan bool)
	tFsm.ChanStopEvent = make(chan bool)
	tFsm.ChanShutdownEvent = make(chan bool)
	tFsm.ChanTradeEvent = make(chan bool)
	tFsm.ChanHoldEvent = make(chan bool)

	// --------------
	// Buy Events
	// --------------
	tFsm.ChanMinute1BuyEvent = make(chan bool)
	tFsm.ChanMinute120BuyEvent = make(chan bool)
	tFsm.ChanMinute60BuyEvent = make(chan bool)
	tFsm.ChanMinute30BuyEvent = make(chan bool)

	tFsm.ChanNotMinute1BuyEvent = make(chan bool)

	tFsm.ChanMinute1SellEvent = make(chan bool)
	tFsm.ChanMinute120SellEvent = make(chan bool)
	tFsm.ChanMinute60SellEvent = make(chan bool)
	tFsm.ChanMinute30SellEvent = make(chan bool)

	tFsm.ChanDoBuyEvent = make(chan bool, 1)
	tFsm.ChanDoSellEvent = make(chan bool, 1)
	tFsm.ChanBuyCompleteEvent = make(chan bool, 1)
	tFsm.ChanSellCompleteEvent = make(chan bool, 1)

	tFsm.testChanDoBuyEvent = make(chan bool, 1)
	tFsm.testChanDoSellEvent = make(chan bool, 1)
	tFsm.testChanBuyCompleteEvent = make(chan bool, 1)
	tFsm.testChanSellCompleteEvent = make(chan bool, 1)

	// -------------
	// FSM creation
	// -------------
	tFsm.FSM = fsm.NewFSM(StartState,
		tFsm.eventsList,
		tFsm.callbacks)

	tFsm.ChanMap = make(map[string]chan bool)
	tFsm.ChanMap[fmt.Sprintf("CEXIO_%s_MS_30_BUY", tFsm.pairID)] = tFsm.ChanMinute30BuyEvent
	tFsm.ChanMap[fmt.Sprintf("CEXIO_%s_MS_60_BUY", tFsm.pairID)] = tFsm.ChanMinute60BuyEvent
	tFsm.ChanMap[fmt.Sprintf("CEXIO_%s_MS_120_BUY", tFsm.pairID)] = tFsm.ChanMinute120BuyEvent

	tFsm.ChanMap[fmt.Sprintf("CEXIO_%s_MS_30_SELL", tFsm.pairID)] = tFsm.ChanMinute30SellEvent
	tFsm.ChanMap[fmt.Sprintf("CEXIO_%s_MS_60_SELL", tFsm.pairID)] = tFsm.ChanMinute60SellEvent
	tFsm.ChanMap[fmt.Sprintf("CEXIO_%s_MS_120_SELL", tFsm.pairID)] = tFsm.ChanMinute120SellEvent
	//ChanDoBuyEvent
	tFsm.ChanMap[fmt.Sprintf("%s_BUY", tFsm.pairID)] = tFsm.ChanDoBuyEvent
	tFsm.ChanMap[fmt.Sprintf("%s_SELL", tFsm.pairID)] = tFsm.ChanDoSellEvent

	tFsm.ChanMap["TRADE"] = tFsm.ChanTradeEvent
	tFsm.ChanMap["HOLD"] = tFsm.ChanHoldEvent
	tFsm.ChanMap["START"] = tFsm.ChanStartEvent
	tFsm.ChanMap["STOP"] = tFsm.ChanStartEvent

	return tFsm

}

func (tFsm *TradeFsm) Kredis() *kredis.Kredis {
	return tFsm.kr
}

func (tFsm *TradeFsm) SignalChannelsMap() map[string]chan bool {
	return tFsm.ChanMap
}

func (tFsm *TradeFsm) FsmController() {

	logMap := make(map[string]bool)
	logMap[Minute120BuyEvent] = false
	logMap[Minute60BuyEvent] = false
	logMap[Minute30BuyEvent] = false

	logMap[NotMinute60BuyEvent] = false
	logMap[NotMinute120BuyEvent] = false
	logMap[NotMinute30BuyEvent] = false

	logMap[Minute120SellEvent] = false
	logMap[Minute60SellEvent] = false
	logMap[Minute30SellEvent] = false

	logMap[NotMinute120SellEvent] = false
	logMap[NotMinute60SellEvent] = false
	logMap[NotMinute30SellEvent] = false

	log.Info("Starting tFsm controlloer for : ", tFsm.pairID)

	for {
		select {

		case _ = <-tFsm.ChanStartEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, StartEvent)
				if err := tFsm.FSM.Event(StartEvent); err != nil {
					//log.Warn(err.Error())
				}

			}
		case _ = <-tFsm.ChanShutdownEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, ShutdownEvent)
				if err := tFsm.FSM.Event(ShutdownEvent); err != nil {
					//log.Warn(err.Error())
				}
			}

		case _ = <-tFsm.ChanStopEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, ShutdownEvent)
				if err := tFsm.FSM.Event(StopEvent); err != nil {
					//log.Warn(err.Error())
				}
			}

		case _ = <-tFsm.ChanTradeEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, TradeEvent)
				if err := tFsm.FSM.Event(TradeEvent); err != nil {
					//log.Warn(err.Error())
				}
			}
		case _ = <-tFsm.ChanHoldEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, HoldEvent)
				if err := tFsm.FSM.Event(HoldEvent); err != nil {
					//log.Warn(err.Error())
				}
			}

		case ev := <-tFsm.ChanMinute120BuyEvent:
			{

				doLog := false
				if logMap[Minute120BuyEvent] != ev {
					doLog = true
				}
				logMap[Minute120BuyEvent] = ev

				if ev {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s , %v", tFsm.pairID, Minute120BuyEvent, ev)
					}
					if err := tFsm.FSM.Event(Minute120BuyEvent); err != nil {
						//log.Warn(err.Error())
					}
				} else {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s , %v", tFsm.pairID, NotMinute120BuyEvent, ev)
					}
					if err := tFsm.FSM.Event(NotMinute120BuyEvent); err != nil {
						//log.Warn(err.Error())
					}
				}

			}
		case ev := <-tFsm.ChanMinute60BuyEvent:
			{

				doLog := false
				if logMap[Minute60BuyEvent] != ev {
					doLog = true
				}
				logMap[Minute60BuyEvent] = ev

				if ev {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, Minute60BuyEvent, ev)
					}

					if err := tFsm.FSM.Event(Minute60BuyEvent); err != nil {
						//log.Warn(err.Error())
					}
				} else {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, NotMinute60BuyEvent, ev)
					}

					if err := tFsm.FSM.Event(NotMinute60BuyEvent); err != nil {
						//log.Warn(err.Error())
					}
				}

			}
		case ev := <-tFsm.ChanMinute30BuyEvent:
			{
				doLog := false
				if logMap[Minute30BuyEvent] != ev {
					doLog = true
				}
				logMap[Minute30BuyEvent] = ev

				if ev {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, Minute30BuyEvent, ev)
					}

					if err := tFsm.FSM.Event(Minute30BuyEvent); err != nil {
						//log.Warn(err.Error())
					}
				} else {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, NotMinute30BuyEvent, ev)
					}

					if err := tFsm.FSM.Event(NotMinute30BuyEvent); err != nil {
						//log.Warn(err.Error())
					}
				}
			}
			/*
				case _ = <-tFsm.ChanNotMinute120BuyEvent:
					{
						log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, NotMinute120BuyEvent)
						if err := tFsm.FSM.Event(NotMinute120BuyEvent); err != nil {
							//log.Warn(err.Error())
						}
					}
				case _ = <-tFsm.ChanNotMinute60BuyEvent:
					{
						log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, NotMinute60BuyEvent)
						if err := tFsm.FSM.Event(NotMinute60BuyEvent); err != nil {
							//log.Warn(err.Error())
						}
					}
				case _ = <-tFsm.ChanNotMinute30BuyEvent:
					{
						log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, NotMinute30BuyEvent)
						if err := tFsm.FSM.Event(NotMinute30BuyEvent); err != nil {
							//log.Warn(err.Error())
						}
					}
			*/
		case ev := <-tFsm.ChanMinute120SellEvent:
			{

				doLog := false
				if logMap[Minute120SellEvent] != ev {
					doLog = true
				}
				logMap[Minute120SellEvent] = ev

				if ev {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, Minute120SellEvent, ev)
					}

					if err := tFsm.FSM.Event(Minute120SellEvent); err != nil {
						//log.Warn(err.Error())
					}
				} else {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, NotMinute120SellEvent, ev)
					}
					if err := tFsm.FSM.Event(NotMinute120SellEvent); err != nil {
						//log.Warn(err.Error())
					}
				}

			}
		case ev := <-tFsm.ChanMinute60SellEvent:
			{

				doLog := false
				if logMap[Minute60SellEvent] != ev {
					doLog = true
				}
				logMap[Minute60SellEvent] = ev

				if ev {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, Minute60SellEvent, ev)
					}

					if err := tFsm.FSM.Event(Minute60SellEvent); err != nil {
						//log.Warn(err.Error())
					}

				} else {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, NotMinute60SellEvent, ev)
					}

					if err := tFsm.FSM.Event(NotMinute60SellEvent); err != nil {
						//log.Warn(err.Error())
					}

				}
			}
		case ev := <-tFsm.ChanMinute30SellEvent:
			{

				doLog := false
				if logMap[Minute30SellEvent] != ev {
					doLog = true
				}
				logMap[Minute30SellEvent] = ev

				if ev {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, Minute30SellEvent, ev)
					}

					if err := tFsm.FSM.Event(Minute30SellEvent); err != nil {
						//log.Warn(err.Error())
					}
				} else {
					if doLog {
						log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, NotMinute30SellEvent, ev)
					}
					if err := tFsm.FSM.Event(NotMinute30SellEvent); err != nil {
						//log.Warn(err.Error())
					}
				}
			}
			/*
				case _ = <-tFsm.ChanNotMinute120SellEvent:

					{
						log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, NotMinute120SellEvent)
						if err := tFsm.FSM.Event(NotMinute120SellEvent); err != nil {
							//log.Warn(err.Error())
						}
					}
				case _ = <-tFsm.ChanNotMinute60SellEvent:
					{
						log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, NotMinute30SellEvent)
						if err := tFsm.FSM.Event(NotMinute30SellEvent); err != nil {
							//log.Warn(err.Error())
						}
					}
				case _ = <-tFsm.ChanNotMinute30SellEvent:
					{
						log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, NotMinute30SellEvent)
						if err := tFsm.FSM.Event(NotMinute30SellEvent); err != nil {
							//log.Warn(err.Error())
						}
					}
			*/
		case _ = <-tFsm.ChanDoBuyEvent:
			{
				log.Infof("tFsm %s Controller Event: %s, %v", tFsm.pairID, DoBuyEvent)
				if err := tFsm.FSM.Event(DoBuyEvent); err != nil {
					//log.Warn(err.Error())
				}
			}
		case _ = <-tFsm.ChanDoSellEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, DoSellEvent)
				if err := tFsm.FSM.Event(DoSellEvent); err != nil {
					//log.Warn(err.Error())
				}
			}
		case _ = <-tFsm.ChanBuyCompleteEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, BuyCompleteEvent)
				if err := tFsm.FSM.Event(BuyCompleteEvent); err != nil {
					//log.Warn(err.Error())
				}
			}
		case _ = <-tFsm.ChanSellCompleteEvent:
			{
				log.Infof("tFsm %s Controller Event: %s", tFsm.pairID, SellCompleteEvent)
				if err := tFsm.FSM.Event(SellCompleteEvent); err != nil {
					//log.Warn(err.Error())
				}
			}

		}
		//time.Sleep(time.Second)
	}
}
