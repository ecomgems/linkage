package tunnel

func (t *Tunnel) errorWhileNotStarted(err error) error {
	t.stats.UpdateState(Stopped)
	t.Unlock()

	return err
}

func (t *Tunnel) errorWhenStarted(err error) {
	t.Lock()

	if t.isStarted() {
		t.cancel()
		t.stats.UpdateState(Stopped)
		t.errChannel <- err
	}

	t.Unlock()
}
