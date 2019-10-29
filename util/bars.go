package util

import "fmt"

var HBarChar = []string{
	` `,
	`▏`,
	`▎`,
	`▍`,
	`▌`,
	`▋`,
	`▊`,
	`▉`,
	`█`,
}
var VBarChar = []string{
	` `,
	`▁`,
	`▂`,
	`▃`,
	`▄`,
	`▅`,
	`▆`,
}
var BlockChar = `█`

// Draws multipart bar. arguments
//
// proportions is just list of ratio between vars
// colors will be colors used. empty string means color will not be set.
func MultipartBar(numChars int, proportion []float32, colors []string) (string,error) {
	if len(proportion) < len(colors)  {
		return "", fmt.Errorf("need more colors than number of variables provided in proportions")
	}
	if len(proportion) < 1 {
		return "", fmt.Errorf("need more colors than number of variables provided in proportions")
	}
	bar := make([]string,numChars)
	var pSum float32
	var pPos := make([]int,proportion)
	for _, v := range proportion {
		pSum += v
	}


	return bar, nil


}