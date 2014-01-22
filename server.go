package main

import (
    "encoding/json"
    "fmt"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "math"
    "net/http"
    "os"
    "strconv"
    "strings"
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

type Similarity struct {
    Left     int64   `db:"fid1"`
    Right    int64   `db:"fid2"`
    Title    float64 `db:"title"`
    Subtitle float64 `db:"sub"`
    Combined float64 `db:"combined"`
    Authors  float64 `db:"authors"`
}

type Entry struct {
    FincId      int64  `db:"finc_id"`
    RecordId    string `db:"record_id"`
    SourceId    int64  `db:"source_id"`
    YearControl string `db:"yearctrl"`
    YearData    string `db:"yeardata"`
    Title       string `db:"combined"`
    Authors     string `db:"authors"`
    Edition     string `db:"edition"`
    Isbns       string `db:"isbns"`
}

type Pair struct {
    Left         Entry
    Rigth        Entry
    Similarities Similarity
}

func main() {

    // load configuration on startup
    file, _ := os.Open("server.json")
    decoder := json.NewDecoder(file)
    configuration := &Configuration{}
    decoder.Decode(&configuration)

    m := martini.Classic()

    m.Map(configuration)
    m.Use(martini.Static("assets"))
    m.Use(render.Renderer(render.Options{
        Directory:  "templates",
        Layout:     "layout",
        Extensions: []string{".tmpl", ".html"},
        Charset:    "UTF-8",
    }))

    m.Get("/s/:sim", func(r render.Render, params martini.Params) {

        db, err := sqlx.Open("sqlite3", "./test.db")
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()
        var entries []Entry
        ids := strings.Split(params["sim"], "-")
        for _, id := range ids {
            entry := Entry{}
            err = db.Get(&entry, `SELECT finc_id, record_id, source_id,
                                yearctrl, yeardata, combined,
                                authors, edition, isbns
                                FROM item WHERE finc_id = ?`, id)
            if err != nil {
                log.Fatal(err)
            }
            log.Println(entry)
            entries = append(entries, entry)
        }
        r.HTML(200, "details", entries)
    })

    m.Get("/", func(r render.Render) {
        r.Redirect("/o/0")
    })

    m.Get("/o", func(r render.Render) {
        r.Redirect("/o/0")
    })

    // show the list
    m.Get("/o/:o", func(r render.Render, params martini.Params) {
        // access the sim db here to build up a list of `Pairs`
        db, err := sqlx.Open("sqlite3", "./test.db")
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()

        offset := 0
        if val, ok := params["o"]; ok {
            offset, _ = strconv.Atoi(val)
        }

        sims := []Similarity{}
        err = db.Select(&sims, `SELECT fid1, fid2, title, sub, combined, authors 
                                FROM similarity ORDER BY fid1 LIMIT 1000 OFFSET ?`, offset)

        total := 0
        row := db.QueryRow(`SELECT count(*) FROM similarity`)
        row.Scan(&total)

        vars := make(map[string]interface{})

        vars["total"] = total
        vars["sims"] = sims
        vars["offset"] = offset
        vars["previous"] = math.Max(0.0, float64(offset-1000))
        vars["next"] = math.Min(float64(offset+1000), float64(total-1000))
        vars["last"] = float64(total - 1000)
        r.HTML(200, "list", vars)
    })

    // testing mysql access
    m.Get("/fid/:fid", func(params martini.Params, r render.Render, c *Configuration) {
        db, err := sqlx.Open(c.Vendor, c.DSN)
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

    http.ListenAndServe("139.18.19.111:3000", m)
    // m.Run()
}
