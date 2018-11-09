/*
Package embed contains builders to create Discord rich messages.

Here is an example of how to create a complex embed :

	e := embed.New().
		Author(
			embed.NewAuthor().
				Name("author name").
				IconURL("http://example.org/icon.png").
				Build(),
		).
		Title("title").
		Description("description").
		Color(0x3277ce).
		Fields(
			embed.NewField().Name("field 1").Value("value 1").Build(),
			embed.NewField().Name("field 2").Value("value 2").Build(),
		Footer(embed.NewFooter().Text("footer").Build()).
		Build()
*/
package embed
