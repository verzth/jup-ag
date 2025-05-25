package tests

import (
	jupiter "git.verzth.work/go/jup-ag"
	"testing"
)

func TestWclient(t *testing.T) {
	cl := jupiter.NewJupagWorker("https://worker.jup.ag", true)

	lstApys, err := cl.LstApys(jupiter.LstApysParams{})

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("lst-apys: %+v", lstApys)
}
