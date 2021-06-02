## statemachine

--



Statemachine is a simple state machine implementation in go.
We can set up a few events and its handlers in order to manage
correctly the satate and dispatche the new state according to the flow
definition.


#### Simple usage 


```go
func main() {
	sm := statemachine.New(
		&statemachine.Event{Name: eventFetch},
		map[string]statemachine.State{
			eventFetch: {
				EventHandler: handleFetchInfoFn,
				Destination:  []string{eventFetch, eventCalculateVolume},
			},
			eventCalculateVolume: {
				EventHandler: handleCalculateVolumeFn,
				Destination:  []string{eventFetch, eventBuySecondaryMarket},
			},
			eventBuySecondaryMarket: {
				EventHandler: ism.HandleBuySecondaryMarket,
				Destination:  []string{eventFetch, eventBuySecondaryMarket, eventCheckOrderExecuted},
			},
			eventCheckOrderExecuted: {
				EventHandler: handleCheckOrderExecutedFn,
				Destination:  []string{eventFinish, eventCheckOrderExecuted, eventCancelOrder},
			},
			eventCancelOrder: {
				EventHandler: handleCancelOrderFn,
				Destination:  []string{eventFinish, eventCancelOrder},
			},
			eventFinish: {
				EventHandler: handleFinishFn,
				Destination:  []string{statemachine.EventAbort},
			},
		},
		logger,
	)
}

func handleFinishFn(e *statemachine.Event) (*statemachine.Event, error) {
	return &statemachine.Event{Name: statemachine.EventAbort}, nil
}
```


You can also print a mermaid diagram based on your statemachine implementation and this is an example of the result

<img width="508" alt="Screen Shot 2021-04-24 at 6 45 01 PM" src="https://user-images.githubusercontent.com/8041435/115974796-6154f080-a52d-11eb-8b7c-19339d2fccbc.png">

For more use cases see the folder examples.

