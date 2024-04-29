package egriden

import (
	"testing"
)

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

	if len(goodSeq.frames) != 3 {
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
}
