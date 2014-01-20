package main

import (
    "encoding/json"
    "fmt"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "os"
    "time"
)

type FincMapping struct {
    FincId   int    `db:"finc_id" json:"finc_id"`
    SourceId int    `db:"source_id" json:"source_id"`
    RecordId string `db:"record_id" json:"record_id"`
}

type Configuration struct {
    Vendor string `json:"vendor"`
    DSN    string `json:"dsn"`
}

func main() {

    // load configuration on startup
    file, _ := os.Open("server.json")
    decoder := json.NewDecoder(file)
    configuration := &Configuration{}
    decoder.Decode(&configuration)

    m := martini.Classic()
    // m.Map(configuration)

    m.Use(render.Renderer(render.Options{
        Directory:  "templates",
        Layout:     "layout",
        Extensions: []string{".tmpl", ".html"},
        Charset:    "UTF-8",
    }))

    m.Get("/", func() string {
        return fmt.Sprintf("Hello world! (%s)", time.Now())
    })

    m.Get("/hello/:name", func(params martini.Params) string {
        return "Hello " + params["name"]
    })

    m.Get("/test", func(r render.Render) {
        r.HTML(200, "hello", "jeremy")
    })

    m.Get("/api", func(r render.Render) {
        r.JSON(200, map[string]interface{}{"hello": "world"})
    })

    m.Get("/fid/:fid", func(params martini.Params, r render.Render) {
        // TODO, ping first
        db, err := sqlx.Open(configuration.Vendor, configuration.DSN)
        if err != nil {
            r.JSON(500, map[string]interface{}{"error": fmt.Sprintf("%s", err)})
            return
        }
        defer db.Close()

        err = db.Ping()
        if err != nil {
            r.JSON(500, map[string]interface{}{"error": fmt.Sprintf("%s", err)})
            return
        }

        // log.Println(params)
        mapping := FincMapping{}
        err = db.Get(&mapping, `SELECT finc_id, record_id, source_id 
                                FROM finc_mapping WHERE finc_id = ?`, params["fid"])
        if err != nil {
            r.JSON(404, map[string]interface{}{"error": fmt.Sprintf("%s", err)})
        } else {
            r.JSON(200, mapping)
        }
    })

    m.Run()
}
