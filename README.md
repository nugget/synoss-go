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

	var (
		json string
		err  error
	)

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
	defer nas.Logout()

	p := make(map[string]string)
	p["version"] = "8"
	p["basic"] = "true"

	json, err = nas.Raw("SYNO.SurveillanceStation.Camera", "List", p)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	fmt.Printf("\n-- \n%v\n-- \n\n", json)
}
```
