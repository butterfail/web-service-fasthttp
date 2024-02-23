package main

import (
	"encoding/json"
	"fmt"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func Albums(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	json, err := json.Marshal(albums)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(json)
}

func postAlbums(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	var album album
	err := json.Unmarshal(ctx.PostBody(), &album)
	if err != nil {
		ctx.Write([]byte(`{"message": "error unmarshalling the request"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	albums = append(albums, album)

	json, err := json.Marshal(albums)
	if err != nil {
		ctx.Write([]byte(`{"message": "error marshalling the albums"}`))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(json)
}

func getAlbumById(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id := ctx.UserValue("id").(string)

	for _, album := range albums {
		if album.ID == id {
			json, err := json.Marshal(album)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			ctx.Write(json)
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.Write([]byte(`{"message": "album not found"}`))
}

func main() {
	r := router.New()
	r.GET("/albums", Albums)
	r.GET("/albums/{id}", getAlbumById)
	r.POST("/albums", postAlbums)

	s := &fasthttp.Server{
		Handler: r.Handler,
		Name:    "Albums API",
	}

	if err := s.ListenAndServe(":80"); err != nil {
		fmt.Println("Error in ListenAndServe: ", err)
	}
}
