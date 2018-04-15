package sarsa

//import "math/rand"

//import "fmt"
import "github.com/faiface/pixel/pixelgl"

//import "os"
//import "os/exec"

//import "time"

type State interface {
	GetRandomFirstPosition() State
	GetActions() []string
	GetActiveTiles(string) [][]int
	InGoalState() bool
	TakeAction(string) (State, float64)
	GetWin() *pixelgl.Window
}

type ValueFunction struct {
	Weights  []float64
	Tilings  int
	Alpha    float64
	Features int
}

//constructor
func (v *ValueFunction) New(feature, max_size, tiling int, alpha float64) {

	v.Weights = make([]float64, max_size)

	v.Tilings = tiling
	v.Alpha = alpha
	v.Features = feature

}

type ActionFunction func(State, *ValueFunction) string
type Valuefunction func(State, string, *ValueFunction) float64

func SemiGradientSarsa(state State, valueOf Valuefunction, GetAction ActionFunction, valueFunction *ValueFunction) int {
	//fmt.Println("hola")
	currentState := state.GetRandomFirstPosition()

	currentAction := GetAction(state, valueFunction)

	steps := 0
	for !currentState.InGoalState() {

		steps += 1
		newState, reward := currentState.TakeAction(currentAction)
		newAction := GetAction(newState, valueFunction)
		target := valueOf(newState, newAction, valueFunction) + reward
		learn(currentState, currentAction, target, valueFunction)
		currentState = newState
		currentAction = newAction

	}
	return steps
}

func learn(state State, action string, target float64, vf *ValueFunction) {

	activeTiles := state.GetActiveTiles(action)

	estimations := make([]float64, vf.Features)

	for feature := 0; feature < vf.Features; feature++ {
		for idx := 0; idx < vf.Tilings; idx++ {
			estimations[feature] += vf.Weights[activeTiles[feature][idx]]
		}
	}
	delta := make([]float64, len(estimations))

	for idx := 0; idx < len(delta); idx++ {
		delta[idx] = /*vf.Alpha **/ 0.0001 * (target - estimations[idx])

	}

	for feature := range delta {
		for tile := range activeTiles[feature] {
			vf.Weights[activeTiles[feature][tile]] += delta[feature]
		}
	}
}
