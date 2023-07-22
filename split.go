package dist

func Split[T any](src <-chan T, cnt int) []<-chan T {
	dests := make([]<-chan T, 0, cnt)

	for i := 0; i < cnt; i++ {
		ch := make(chan T)
		dests = append(dests, ch)

		go func() {
			defer close(ch)
			for v := range src {
				ch <- v
			}
		}()
	}

	return dests
}
