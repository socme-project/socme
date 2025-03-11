package notifyme

type Notifier interface {
	Notify(title, message string) error
}

type Notifiers []Notifier

func (n Notifiers) Notify(title, message string) []error {
	var errs []error = nil
	for _, notifier := range n {
		if err := notifier.Notify(title, message); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
