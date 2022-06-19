# supabase-go

Unofficial [Supabase](https://supabase.io) client for Go. It is an amalgamation of all the libraries similar to the [official Supabase client](https://supabase.io/docs/reference/javascript/supabase-client).

## Installation
```
go get github.com/grid-rbx/supabase-go
```

## Usage
```golang
import supabase "github.com/grid-rbx/supabase-go"

func main() {
  supabaseUrl := "<SUPABASE-URL>"
  supabaseKey := "<SUPABASE-KEY>"
  supabaseClient := supabase.CreateClient(supabaseUrl, supabaseKey)

  // Auth
  user, err := supabaseClient.Auth.SignIn(supabase.UserCredentials{
    email: "example@example.com",
    password: "password"
  })
  if err != nil {
    panic(err)
  }

  // DB
  var results map[string]interface{}
  err = supabaseClient.DB.From("something").Select("*").Single().Execute(&results)
  if err != nil {
    panic(err)
  }

  fmt.Println(results)

  // Storage
  supabaseClient.Storage.UploadFile("bucket-name", "file.txt", file)

  // Realtime
  supabaseClient.Realtime.Connect()

  channel, err := supabaseClient.Realtime.Channel(supabaseClient.Realtime.WithTable("realtime", "public", "my_table"))

  channel.OnInsert = func(m supabaseClient.Realtime.Message) {
    fmt.Println(m)
  }
}
```

## Roadmap
- [x] Auth support (1)
- [x] DB support (2)
- [x] Realtime
- [x] Storage
- [ ] Testing

(1) - Thin API wrapper. Does not rely on the GoTrue library for simplicity
(2) - Through `postgrest-go`

I just implemented features which I actually needed for my project for now. If you like to implement these features, feel free to submit a pull request as stated in the [Contributing](#contributing) section below.

## Design Goals
It tries to mimick as much as possible the official Javascript client library in terms of ease-of-use and in setup process.

# Contributing
## Submitting a pull request
- Fork it (https://github.com/grid-rbx/supabase-go/fork)
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create a new Pull Request

# Contributors
- [nedpals](https://github.com/nedpals) - creator and maintainer
- [dominictwlee](https://github.com/dominictwlee) - supabase-community/postgrest-go implementation
- [cursecodes](https://github.com/cursecodes) - storage implementation
- [overseedio](https://github.com/overseedio) - realtime implementation