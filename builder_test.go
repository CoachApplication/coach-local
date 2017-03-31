package local_test

func isInSlice(vl string, sl []string) bool {
	for _, each := range sl {
		if each == vl {
			return true
		}
	}
	return false
}
