package statemachine

var (
	eventOne   = "one"
	eventTwo   = "two"
	eventThree = "three"
)

type implementation struct {
	maxCycles       int
	cycles          int
	eventOneCalls   []int
	eventTwoCalls   []int
	eventThreeCalls []int
}

func (s *implementation) eventOneHandler(e *Event) (*Event, error) {
	if s.cycles == s.maxCycles {
		return &Event{Name: EventAbort}, nil
	}

	s.eventOneCalls = append(s.eventOneCalls, len(s.eventOneCalls)+1)
	s.cycles = s.cycles + 1
	return &Event{Name: eventTwo}, nil
}

func (s *implementation) eventTwoHandler(e *Event) (*Event, error) {
	s.eventTwoCalls = append(s.eventTwoCalls, len(s.eventTwoCalls)+1)
	return &Event{Name: eventThree}, nil
}

func (s *implementation) eventThreeHandler(e *Event) (*Event, error) {
	s.eventThreeCalls = append(s.eventThreeCalls, len(s.eventThreeCalls)+1)
	return &Event{Name: eventOne}, nil
}
