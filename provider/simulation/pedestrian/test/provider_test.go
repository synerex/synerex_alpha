
import (
	"testing"
)

// Spec
// 起動時にGetAreaが呼ばれる

func TestDuck_Say(t *testing.T) {
	t.Run("it says quack", func(t *testing.T) {
        actual := duck.Say()
        expected := "tarou says quack"
        if actual != expected {
            t.Errorf("got: %v\nwant: %v", actual, expected)
        }
    })
}