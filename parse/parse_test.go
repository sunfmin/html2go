package parse_test

import (
	"strings"
	"testing"

	"github.com/sunfmin/html2go/parse"
	"github.com/theplant/testingutils"
)

func TestAll(t *testing.T) {
	var cases = []struct {
		name   string
		pkg    string
		html   string
		gocode string
	}{
		{
			name: "normal",
			html: `
<nav class="navbar navbar-expand-lg navbar-light bg-light">
  <a class="navbar-brand" href="#">Navbar</a>
  <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
    <span class="navbar-toggler-icon"></span>
  </button>
  <div class="collapse navbar-collapse" id="navbarNav">
    <ul class="navbar-nav">
      <li class="nav-item active">
        <a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="#">Features</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="#">Pricing</a>
      </li>
      <li class="nav-item">
        <a class="nav-link disabled" href="#" tabindex="-1" aria-disabled="true">Disabled</a>
      </li>
    </ul>
  </div>
</nav>
`,
			gocode: `package hello

var n = Body(
	Nav(
		A(
			Text("Navbar"),
		).Class("navbar-brand").
			Href("#"),
		Button("").Class("navbar-toggler").
			Type("button").
			Attr("data-toggle", "collapse").
			Attr("data-target", "#navbarNav").
			Attr("aria-controls", "navbarNav").
			Attr("aria-expanded", "false").
			Attr("aria-label", "Toggle navigation").
			Children(
				Span("").Class("navbar-toggler-icon"),
			),
		Div(
			Ul(
				Li(
					A(
						Text("Home"),
						Span("(current)").Class("sr-only"),
					).Class("nav-link").
						Href("#"),
				).Class("nav-item active"),
				Li(
					A(
						Text("Features"),
					).Class("nav-link").
						Href("#"),
				).Class("nav-item"),
				Li(
					A(
						Text("Pricing"),
					).Class("nav-link").
						Href("#"),
				).Class("nav-item"),
				Li(
					A(
						Text("Disabled"),
					).Class("nav-link disabled").
						Href("#").
						TabIndex(-1).
						Attr("aria-disabled", "true"),
				).Class("nav-item"),
			).Class("navbar-nav"),
		).Class("collapse navbar-collapse").
			Id("navbarNav"),
	).Class("navbar navbar-expand-lg navbar-light bg-light"),
)
`,
		},
		{
			name: "bool attribute always be true even if value is false",
			html: `
<nav class="navbar navbar-expand-lg navbar-light bg-light">
  <input readonly required disabled checked tabindex="-1">
  <input readonly="false">
</nav>
`,
			gocode: `package hello

var n = Body(
	Nav(
		Input("").Readonly(true).
			Required(true).
			Disabled(true).
			Checked(true).
			TabIndex(-1),
		Input("").Readonly(true),
	).Class("navbar navbar-expand-lg navbar-light bg-light"),
)
`,
		},
		{
			name: "text attr with text",
			html: `
<div>
  <span>Hello</span>
</div>
`,
			gocode: `package hello

var n = Body(
	Div(
		Span("Hello"),
	),
)
`,
		},

		{
			name: "text attr with more children",
			html: `
<div>
  <span>Hello<b>world</b></span>
</div>
`,
			gocode: `package hello

var n = Body(
	Div(
		Span("").
			Children(
				Text("Hello"),
				B("world"),
			),
	),
)
`,
		},
		{
			name: "text attr on tag children",
			html: `
<div>
  <span><b>world</b></span>
</div>
`,
			gocode: `package hello

var n = Body(
	Div(
		Span("").
			Children(
				B("world"),
			),
	),
)
`,
		},
		{
			name: "text attr with more children with pkg",
			pkg:  "h",
			html: `
<div>
  <span>Hello<b>world</b></span>
</div>
`,
			gocode: `package hello

var n = h.Body(
	h.Div(
		h.Span("").
			Children(
				h.Text("Hello"),
				h.B("world"),
			),
	),
)
`,
		},

		{
			name: "code with javascript",
			pkg:  "",
			html: `
<div class="antialiased min-h-screen bg-gray-100 flex items-center">
  
<div class="w-full max-w-sm mx-auto">
  <div class="bg-white shadow rounded-lg p-5 dark:bg-gray-800 w-full"
	x-data="{
		weatherData: null,
		fetchWeatherData() {
			fetch('https://api.weatherapi.com/v1/forecast.json?key=ff9b41622f994b1287a73535210809&q=Guwahati&days=3')
				.then(response => response.json())
				.then(json => this.weatherData = json)
		},
		formattedDateDisplay(date) {
			const options = {
				weekday: 'long',
				year: 'numeric',
				month: 'long',
				day: 'numeric'
		   };
		   
		   return (new Date(date)).toLocaleDateString('en-US', options);
		}
	}"
	x-init="fetchWeatherData()"
	x-cloak
>
	<h2 class="font-bold text-gray-800 text-lg dark:text-gray-400" x-text="formattedDateDisplay(new Date())"></h2>
	
	<template x-if="weatherData != null">
		<div>
			<div class="flex mt-4 mb-2">
				<div class="flex-1">
					<div class="text-gray-600 text-sm dark:text-gray-400" x-text="weatherData.location.name +', '+ weatherData.location.region"></div>
					<div class="text-3xl font-bold text-gray-800 dark:text-gray-300" x-html="|backquote|${weatherData.current.temp_c} &deg;C|backquote|"></div>
					<div x-text="weatherData.current.condition.text" class="text-xs text-gray-600 dark:text-gray-400"></div>
				</div>
				<div class="w-24">
					<img :src="|backquote|https:${weatherData.current.condition.icon}|backquote|" :alt="weatherData.current.condition.text" loading="lazy">
				</div>
			</div>

			<div class="flex space-x-2 justify-between border-t dark:border-gray-500">
				<template x-for="(forecast, key) in weatherData.forecast.forecastday.splice(1)">
					<div class="flex-1 text-center pt-3" :class="{'border-r dark:border-gray-500': key==0}">
						<div x-text="|backquote|${forecast.date.split('-')[2]}/${forecast.date.split('-')[1]}/${forecast.date.split('-')[0]}|backquote|" class="text-xs text-gray-500 dark:text-gray-400"></div>
						<img :src="|backquote|https:${forecast.day.condition.icon}|backquote|" :alt="forecast.day.condition.text" loading="lazy" class="mx-auto">
						<div x-html="|backquote|${forecast.day.maxtemp_c} &deg;C|backquote|" class="font-semibold text-gray-800 mt-1.5 dark:text-gray-300"></div>
						<div x-text="forecast.day.condition.text" class="text-xs text-gray-600 dark:text-gray-400"></div>
					</div>
				</template>
			</div>
		</div>
	</template>

	<template x-if="weatherData == null">
		<div class="animate-pulse">
			<div class="flex mt-4 mb-5">
				<div class="flex-1">
					<div class="rounded h-2 mb-1.5 bg-gray-200 w-1/2"></div>
					<div class="bg-gray-200 rounded h-4"></div>
					<div class="rounded h-2 mt-1.5 bg-gray-200 w-1/2"></div>
				</div>
				<div class="w-24">
					<div class="w-12 h-12 rounded-full bg-gray-100 mx-auto"></div>
				</div>
			</div>

			<div class="flex space-x-2 justify-between border-t h-32 dark:border-gray-500">
				<div class="flex-1 text-center pt-4 border-r px-5 dark:border-gray-500">
					<div class="rounded h-2 mb-2 bg-gray-200 w-1/2 mx-auto"></div>
					<div class="w-12 h-12 rounded-full bg-gray-100 mx-auto mb-2"></div>
					<div class="rounded h-3 mt-1 bg-gray-200 mt-1.5 mx-auto"></div>
					<div class="rounded h-2 mt-1 bg-gray-200 w-1/2 mx-auto"></div>

				</div>
				<div class="flex-1 text-center pt-4 px-5">
					<div class="rounded h-2 mb-2 bg-gray-200 w-1/2 mx-auto"></div>
					<div class="w-12 h-12 rounded-full bg-gray-100 mx-auto mb-2"></div>
					<div class="rounded h-3 mt-1 bg-gray-200 mt-1.5 mx-auto"></div>
					<div class="rounded h-2 mt-1 bg-gray-200 w-1/2 mx-auto"></div>
				</div>
			</div>
		</div>
	</template>
</div>

  <div class="mt-10">
                Build with <a class="underline text-purple-600" href="weatherapi.com">Weatherapi.com</a>, <a class="underline text-purple-600" href="https://tailwindcss.com/">TailwindCSS</a> & <a class="underline text-purple-600" href="https://alpinejs.dev/">Alpine.js</a>.
                <br>
                Github Gist: <a class="underline text-purple-600" href="https://gist.github.com/mithicher/70647625163c3d217c742dfc4a0b84c0">Weather Forecast Blade Component</a>
  </div>
</div>
</div
`,
			gocode: `package hello

var n = Body(
	Div(
		Div(
			Div(
				H2("").Class("font-bold text-gray-800 text-lg dark:text-gray-400").
					Attr("x-text", "formattedDateDisplay(new Date())"),
				Template(
					Div(
						Div(
							Div(
								Div().Class("text-gray-600 text-sm dark:text-gray-400").
									Attr("x-text", |backquote|weatherData.location.name +", "+ weatherData.location.region|backquote|),
								Div().Class("text-3xl font-bold text-gray-800 dark:text-gray-300").
									Attr("x-html", "|backquote|${weatherData.current.temp_c} °C|backquote|"),
								Div().Attr("x-text", "weatherData.current.condition.text").
									Class("text-xs text-gray-600 dark:text-gray-400"),
							).Class("flex-1"),
							Div(
								Img("").Attr("x-bind:src", "|backquote|https:${weatherData.current.condition.icon}|backquote|").
									Attr("x-bind:alt", "weatherData.current.condition.text").
									Attr("loading", "lazy"),
							).Class("w-24"),
						).Class("flex mt-4 mb-2"),
						Div(
							Template(
								Div(
									Div().Attr("x-text", "|backquote|${forecast.date.split(\"-\")[2]}/${forecast.date.split(\"-\")[1]}/${forecast.date.split(\"-\")[0]}|backquote|").
										Class("text-xs text-gray-500 dark:text-gray-400"),
									Img("").Attr("x-bind:src", "|backquote|https:${forecast.day.condition.icon}|backquote|").
										Attr("x-bind:alt", "forecast.day.condition.text").
										Attr("loading", "lazy").
										Class("mx-auto"),
									Div().Attr("x-html", "|backquote|${forecast.day.maxtemp_c} °C|backquote|").
										Class("font-semibold text-gray-800 mt-1.5 dark:text-gray-300"),
									Div().Attr("x-text", "forecast.day.condition.text").
										Class("text-xs text-gray-600 dark:text-gray-400"),
								).Class("flex-1 text-center pt-3").
									Attr("x-bind:class", |backquote|{"border-r dark:border-gray-500": key==0}|backquote|),
							).Attr("x-for", "(forecast, key) in weatherData.forecast.forecastday.splice(1)"),
						).Class("flex space-x-2 justify-between border-t dark:border-gray-500"),
					),
				).Attr("x-if", "weatherData != null"),
				Template(
					Div(
						Div(
							Div(
								Div().Class("rounded h-2 mb-1.5 bg-gray-200 w-1/2"),
								Div().Class("bg-gray-200 rounded h-4"),
								Div().Class("rounded h-2 mt-1.5 bg-gray-200 w-1/2"),
							).Class("flex-1"),
							Div(
								Div().Class("w-12 h-12 rounded-full bg-gray-100 mx-auto"),
							).Class("w-24"),
						).Class("flex mt-4 mb-5"),
						Div(
							Div(
								Div().Class("rounded h-2 mb-2 bg-gray-200 w-1/2 mx-auto"),
								Div().Class("w-12 h-12 rounded-full bg-gray-100 mx-auto mb-2"),
								Div().Class("rounded h-3 mt-1 bg-gray-200 mt-1.5 mx-auto"),
								Div().Class("rounded h-2 mt-1 bg-gray-200 w-1/2 mx-auto"),
							).Class("flex-1 text-center pt-4 border-r px-5 dark:border-gray-500"),
							Div(
								Div().Class("rounded h-2 mb-2 bg-gray-200 w-1/2 mx-auto"),
								Div().Class("w-12 h-12 rounded-full bg-gray-100 mx-auto mb-2"),
								Div().Class("rounded h-3 mt-1 bg-gray-200 mt-1.5 mx-auto"),
								Div().Class("rounded h-2 mt-1 bg-gray-200 w-1/2 mx-auto"),
							).Class("flex-1 text-center pt-4 px-5"),
						).Class("flex space-x-2 justify-between border-t h-32 dark:border-gray-500"),
					).Class("animate-pulse"),
				).Attr("x-if", "weatherData == null"),
			).Class("bg-white shadow rounded-lg p-5 dark:bg-gray-800 w-full").
				Attr("x-data", |backquote|{
		weatherData: null,
		fetchWeatherData() {
			fetch("https://api.weatherapi.com/v1/forecast.json?key=ff9b41622f994b1287a73535210809&q=Guwahati&days=3")
				.then(response => response.json())
				.then(json => this.weatherData = json)
		},
		formattedDateDisplay(date) {
			const options = {
				weekday: "long",
				year: "numeric",
				month: "long",
				day: "numeric"
		   };
		   
		   return (new Date(date)).toLocaleDateString("en-US", options);
		}
	}|backquote|).
				Attr("x-init", "fetchWeatherData()").
				Attr("x-cloak", ""),
			Div(
				Text("Build with"),
				A(
					Text("Weatherapi.com"),
				).Class("underline text-purple-600").
					Href("weatherapi.com"),
				Text(","),
				A(
					Text("TailwindCSS"),
				).Class("underline text-purple-600").
					Href("https://tailwindcss.com/"),
				Text("&"),
				A(
					Text("Alpine.js"),
				).Class("underline text-purple-600").
					Href("https://alpinejs.dev/"),
				Text("."),
				Br(),
				Text("Github Gist:"),
				A(
					Text("Weather Forecast Blade Component"),
				).Class("underline text-purple-600").
					Href("https://gist.github.com/mithicher/70647625163c3d217c742dfc4a0b84c0"),
			).Class("mt-10"),
		).Class("w-full max-w-sm mx-auto"),
	).Class("antialiased min-h-screen bg-gray-100 flex items-center"),
)
`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gocode := parse.GenerateHTMLGo(c.pkg, strings.NewReader(
				strings.ReplaceAll(c.html, "|backquote|", "`"),
			))
			diff := testingutils.PrettyJsonDiff(strings.ReplaceAll(c.gocode, "|backquote|", "`"), gocode)

			if len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}
