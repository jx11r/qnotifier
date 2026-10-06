package discord

type Payload struct {
	Username string  `json:"username"`
	Embeds   []Embed `json:"embeds"`
}

type Embed struct {
	Title  string       `json:"title"`
	URL    string       `json:"url"`
	Color  int          `json:"color"`
	Image  *EmbedImage  `json:"image,omitempty"`
	Author *EmbedAuthor `json:"author,omitempty"`
}

type EmbedImage struct {
	URL string `json:"url"`
}

type EmbedAuthor struct {
	Name string `json:"name"`
	Icon string `json:"icon_url"`
	URL  string `json:"url"`
}
