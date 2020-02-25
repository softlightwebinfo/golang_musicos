package libs

import "fmt"

func Base64Image(image string) string {
	if image != "" {
		return fmt.Sprintf("data:image/png;base64,%s", image)
	}
	return image
}
