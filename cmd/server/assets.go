package main

import (
	"aletheiaware.com/netgo/handler"
	"aletheiaware.com/perspectivefungo"
	"log"
	"net/http"
)

// TinyGo doesn't support embed yet, so they are fetched from the server separately

func AttachAssetHandlers(mux *http.ServeMux) {
	for a, d := range map[string][]byte{
		"/assets/Shader.vp":    []byte(perspectivefungo.VertexShader),
		"/assets/Shader.fp":    []byte(perspectivefungo.FragmentShader),
		"/assets/Block.off":    perspectivefungo.Block,
		"/assets/Goal.off":     perspectivefungo.Goal,
		"/assets/Player.off":   perspectivefungo.Player,
		"/assets/Portal.off":   perspectivefungo.Portal,
		"/assets/Start.off":    perspectivefungo.Start,
		"/assets/GameOver.off": perspectivefungo.GameOver,
		"/assets/Retry.off":    perspectivefungo.Retry,
		"/assets/Share.off":    perspectivefungo.Share,
		"/assets/Zero.off":     perspectivefungo.Zero,
		"/assets/One.off":      perspectivefungo.One,
		"/assets/Two.off":      perspectivefungo.Two,
		"/assets/Three.off":    perspectivefungo.Three,
		"/assets/Four.off":     perspectivefungo.Four,
		"/assets/Five.off":     perspectivefungo.Five,
		"/assets/Six.off":      perspectivefungo.Six,
		"/assets/Seven.off":    perspectivefungo.Seven,
		"/assets/Eight.off":    perspectivefungo.Eight,
		"/assets/Nine.off":     perspectivefungo.Nine,
		"/assets/Point.off":    perspectivefungo.Point,
		"/assets/Seconds.off":  perspectivefungo.Seconds,
	} {
		asset := a
		data := d
		mux.Handle(asset, handler.Log(handler.Compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(data)
			if err != nil {
				log.Println(err)
				return
			}
		}))))
	}
}
