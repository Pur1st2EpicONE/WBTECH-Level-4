package or

import "sync"

func Or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan any)
	go func() {
		var once sync.Once
		for _, channel := range channels {
			go func(channel <-chan any) {
				select {
				case <-channel:
					once.Do(func() { close(orDone) })
				case <-orDone:
				}
			}(channel)
		}
	}()
	return orDone
}
