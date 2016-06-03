package mytags

import (
    "testing"
    "fmt"
)

func TestGetMyTags(t *testing.T) {
	
	employertags :=[]string{"java","php"}
	
	out := GetMyTags("/home/juno/git/jobprotractor/gojobextractor/mytags.csv",employertags)
	
	fmt.Println(out)

}

