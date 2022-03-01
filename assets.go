package perspectivefungo

import _ "embed"

//go:embed assets/Shader.vp
var VertexShader string

//go:embed assets/Shader.fp
var FragmentShader string

//go:embed assets/Block.off
var Block []byte

//go:embed assets/Goal.off
var Goal []byte

//go:embed assets/Player.off
var Player []byte

//go:embed assets/Portal.off
var Portal []byte

//go:embed assets/Start.off
var Start []byte

//go:embed assets/GameOver.off
var GameOver []byte

//go:embed assets/Retry.off
var Retry []byte

//go:embed assets/Share.off
var Share []byte

//go:embed assets/Zero.off
var Zero []byte

//go:embed assets/One.off
var One []byte

//go:embed assets/Two.off
var Two []byte

//go:embed assets/Three.off
var Three []byte

//go:embed assets/Four.off
var Four []byte

//go:embed assets/Five.off
var Five []byte

//go:embed assets/Six.off
var Six []byte

//go:embed assets/Seven.off
var Seven []byte

//go:embed assets/Eight.off
var Eight []byte

//go:embed assets/Nine.off
var Nine []byte

//go:embed assets/Point.off
var Point []byte

//go:embed assets/Seconds.off
var Seconds []byte

func LoadAssets(d Driver) error {
	if err := LoadOFFMesh(d, "block", Block, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "goal", Goal, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "player", Player, true); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "portal", Portal, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "start", Start, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "gameover", GameOver, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "retry", Retry, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "share", Share, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "0", Zero, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "1", One, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "2", Two, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "3", Three, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "4", Four, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "5", Five, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "6", Six, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "7", Seven, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "8", Eight, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "9", Nine, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, ".", Point, false); err != nil {
		return err
	}

	if err := LoadOFFMesh(d, "s", Seconds, false); err != nil {
		return err
	}
	return nil
}
