package main

import (
	"aletheiaware.com/perspectivefungo"
	"errors"
	"fmt"
	"io"
	"syscall/js"
)

// TinyGo doesn't support embed yet, so they are fetched from the server separately

func loadAssets() error {
	if data, err := fetch("assets/Shader.vp"); err != nil {
		return err
	} else {
		perspectivefungo.VertexShader = string(data)
	}
	if data, err := fetch("assets/Shader.fp"); err != nil {
		return err
	} else {
		perspectivefungo.FragmentShader = string(data)
	}
	if data, err := fetch("assets/Block.off"); err != nil {
		return err
	} else {
		perspectivefungo.Block = data
	}
	if data, err := fetch("assets/Goal.off"); err != nil {
		return err
	} else {
		perspectivefungo.Goal = data
	}
	if data, err := fetch("assets/Player.off"); err != nil {
		return err
	} else {
		perspectivefungo.Player = data
	}
	if data, err := fetch("assets/Portal.off"); err != nil {
		return err
	} else {
		perspectivefungo.Portal = data
	}
	if data, err := fetch("assets/Start.off"); err != nil {
		return err
	} else {
		perspectivefungo.Start = data
	}
	if data, err := fetch("assets/GameOver.off"); err != nil {
		return err
	} else {
		perspectivefungo.GameOver = data
	}
	if data, err := fetch("assets/Retry.off"); err != nil {
		return err
	} else {
		perspectivefungo.Retry = data
	}
	if data, err := fetch("assets/Share.off"); err != nil {
		return err
	} else {
		perspectivefungo.Share = data
	}
	if data, err := fetch("assets/Zero.off"); err != nil {
		return err
	} else {
		perspectivefungo.Zero = data
	}
	if data, err := fetch("assets/One.off"); err != nil {
		return err
	} else {
		perspectivefungo.One = data
	}
	if data, err := fetch("assets/Two.off"); err != nil {
		return err
	} else {
		perspectivefungo.Two = data
	}
	if data, err := fetch("assets/Three.off"); err != nil {
		return err
	} else {
		perspectivefungo.Three = data
	}
	if data, err := fetch("assets/Four.off"); err != nil {
		return err
	} else {
		perspectivefungo.Four = data
	}
	if data, err := fetch("assets/Five.off"); err != nil {
		return err
	} else {
		perspectivefungo.Five = data
	}
	if data, err := fetch("assets/Six.off"); err != nil {
		return err
	} else {
		perspectivefungo.Six = data
	}
	if data, err := fetch("assets/Seven.off"); err != nil {
		return err
	} else {
		perspectivefungo.Seven = data
	}
	if data, err := fetch("assets/Eight.off"); err != nil {
		return err
	} else {
		perspectivefungo.Eight = data
	}
	if data, err := fetch("assets/Nine.off"); err != nil {
		return err
	} else {
		perspectivefungo.Nine = data
	}
	if data, err := fetch("assets/Point.off"); err != nil {
		return err
	} else {
		perspectivefungo.Point = data
	}
	if data, err := fetch("assets/Seconds.off"); err != nil {
		return err
	} else {
		perspectivefungo.Seconds = data
	}
	return nil
}

func fetch(uri string) ([]byte, error) {
	response := await(js.Global().Get("fetch").Invoke(uri))[0]

	b := response.Get("body")
	if b.IsUndefined() || b.IsNull() {
		return nil, fmt.Errorf("body is undefined or null")
	}
	var (
		done   bool
		data   []byte
		stream = b.Call("getReader")
		dc     = make(chan []byte, 1)
		ec     = make(chan error, 1)
	)
	success := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := args[0]
		if result.Get("done").Bool() {
			ec <- io.EOF
			return nil
		}
		value := make([]byte, result.Get("value").Get("byteLength").Int())
		js.CopyBytesToGo(value, result.Get("value"))
		dc <- value
		return nil
	})
	defer success.Release()
	failure := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ec <- errors.New(args[0].Get("message").String())
		return nil
	})
	defer failure.Release()
	for !done {
		stream.Call("read").Call("then", success, failure)
		select {
		case d := <-dc:
			data = append(data, d...)
		case e := <-ec:
			done = true
			if e != io.EOF {
				return nil, e
			}
		}
	}
	return data, nil
}

func await(awaitable js.Value) []js.Value {
	ch := make(chan []js.Value)
	awaitable.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args
		return nil
	}))
	return <-ch
}
