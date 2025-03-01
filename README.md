# Discord Video Embedder

A wrapper that helps with discord video embeds that doesn't work.

the wrapper uses:
- [discord.nfp.is](https://discord.nfp.is/) for embed link generation
- [catbox.moe](https://catbox.moe/) for file upload

### Get Module
```sh
go get -u https://github.com/reiyuchan/discord-video-embedder-go 
```

### Example:
```go
package main

import (
	"fmt"

	discordvideoembedder "github.com/reiyuchan/discord-video-embedder-go"
)

func main() {
	 c := discordembedder.New(nil)
	 fURL, err := c.UploadToCatBox("file.mp4")
	 if err != nil {
	 	fmt.Println(err)
	 }
	 fmt.Printf("%s\n", fURL) // https://files.catbox.moe/<id>.mp4
   eURL, err := c.GetURL(fURL) // fURL -> https://files.catbox.moe/<id>.mp4
	 if err != nil {
	 	fmt.Println(err)
	 }
	 fmt.Printf("%s\n", eURL) // eURL -> https://discord.nfp.is/<id>
}

```
