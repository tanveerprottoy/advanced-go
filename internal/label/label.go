package label

func LabelExample() {
	a := [5][5]int{{1, 2}, {2, 5}}
OuterLoop:
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			switch a[i][j] {
			case 0:
				break OuterLoop
			case 1:
				break OuterLoop
			}
		}
	}
}
