# synoss-go

# Sample Usage

```
package main

import (
	"fmt"
	"os"

	"github.com/nugget/synoss-go"
)

func main() {
	nas := synoss.New()

	var err error

	err = nas.Connect("https://synology.example.org:5001")
	if err != nil {
		fmt.Println("CONNECT ERROR", err)
		os.Exit(1)
	}

	err = nas.Login("username", "password")
	if err != nil {
		fmt.Println("LOGIN ERROR", err)
		os.Exit(1)
	}

}
```
