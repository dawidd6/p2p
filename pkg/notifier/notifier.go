package notifier

type Notifier struct {
	channel chan struct{}
}

func New(blocking bool) *Notifier {
	if blocking {
		return &Notifier{
			channel: make(chan struct{}),
		}
	}

	return &Notifier{
		channel: make(chan struct{}, 1),
	}
}

func (notifier *Notifier) Notify() {
	if len(notifier.channel) > 0 {
		<-notifier.channel
	}
	notifier.channel <- struct{}{}
}

func (notifier *Notifier) Wait() {
	<-notifier.channel
}
