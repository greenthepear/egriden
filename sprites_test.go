package egriden

import (
	"embed"
	_ "image/png"
	"testing"
)

//go:embed examples/gridsweeper/Graphics/*
var GridsweeperGraphics embed.FS

func TestSequenceAndPackCreation(t *testing.T) {
	_, err := CreateImageSequenceFromFolder("bad", "./fakedir")
	if err == nil {
		t.Errorf("failed to report bad path when creating image sequence from folder")
	}

	goodSeq, err := CreateImageSequenceFromFolder("good",
		"./examples/gridsweeper/Graphics/unrevealed")

	if err != nil {
		t.Errorf("normal error while creating image sequence: %v", err)
	}

	if len(goodSeq.Frames) != 3 {
		t.Errorf("wrong number of frames loaded")
	}

	emptyPack := EmptySpritePack()
	if emptyPack.Sprite() != nil {
		t.Errorf("empty SpritePack didn't return a nil")
	}

	sp1 := NewSpritePack()
	sp1.AddImageSequence(goodSeq)

	if sp1.Sprite() == nil {
		t.Errorf("SpritePack with good image sequence returns nil sprite")
	}

	yaml := []byte(`
spritepacks:
  - name: sp2
    sequences:
    - name: good
      paths:
        - examples/gridsweeper/Graphics/unrevealed/unrevealed0.png
        - examples/gridsweeper/Graphics/unrevealed/unrevealed1.png
        - examples/gridsweeper/Graphics/unrevealed/unrevealed2.png
`)

	sp2map, err := CreateSpritePacksFromYaml(yaml, GridsweeperGraphics)
	if err != nil {
		t.Errorf("normal error from CreateSpritePacksFromYaml: %v", err)
	}
	sp2, ok := sp2map["sp2"]
	if !ok {
		t.Errorf("no value for key 'sp2' in SpritePack map from CreateSpritePacksFromYaml")
	}
	if len(sp2.sequences["good"].Frames) != 3 {
		t.Errorf("wrong number of frames loaded")
	}
}
