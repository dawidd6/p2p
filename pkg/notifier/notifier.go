package notifier

type Notifier struct {
	channel chan struct{}
}

func NewBlocking() *Notifier {
	return &Notifier{
		channel: make(chan struct{}),
	}
}

func NewNotBlocking() *Notifier {
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

func (notifier *Notifier) Wait() <-chan struct{} {
	return notifier.channel
}
