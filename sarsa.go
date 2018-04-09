package sarsa

import "math/rand"

//import "fmt"

//import "time"

type State interface {
	GetRandomFirstPosition() State
	GetActions() []string
	GetActiveTiles(string) [][]int
	InGoalState() bool
	TakeAction(string) (State, float64)
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

func SemiGradientSarsa(state State, valueFunction *ValueFunction) int {

	//random position in range (-0.6,-0.4)
	currentState := state.GetRandomFirstPosition()

	currentAction := getAction(state, valueFunction)

	steps := 0
	idx := 0
	for idx < 10000 { //!currentState.InGoalState() {
		idx++
		//fmt.Println(steps)
		steps += 1
		//applies the current action to the current state
		newState, reward := currentState.TakeAction(currentAction)
		//Get best action given a position and a velocity
		newAction := getAction(newState, valueFunction)

		target := valueOf(newState, newAction, valueFunction) + reward
		//if target > 20 {
		//time.Sleep(2000 * time.Millisecond)
		//fmt.Println("Pesos: ", valueFunction.Weights[:100])
		//fmt.Println("Valor del siguiente estado: ", target-reward)
		//fmt.Println("Y la recompensa del estado actual es:", reward)

		learn(currentState, currentAction, target, valueFunction)

		currentState = newState
		currentAction = newAction

	}
	return steps
}

func getAction(state State, vf *ValueFunction) string {
	values := make([]float64, 0)
	actions := state.GetActions()
	for _, action := range actions {
		values = append(values, valueOf(state, action, vf))
	}
	ac := actions[getIdxMax(values)]
	//fmt.Println("agarre", ac, " porque es el mas grande en ", values, getIdxMax(values), actions)
	return ac
}

func valueOf(state State, action string, vf *ValueFunction) float64 {
	if state.InGoalState() {
		return 0.0
	}

	activeTiles := state.GetActiveTiles(action)
	estimations := make([]float64, vf.Features)

	for feature := 0; feature < vf.Features; feature++ {
		for idx := 0; idx < vf.Tilings; idx++ {
			estimations[feature] += vf.Weights[activeTiles[feature][idx]]
		}
	}
	val := 0.0
	for estimation := range estimations {
		val += estimations[estimation]
	}

	return val

	/*idxActiveTiles := state.GetActiveTiles(action)
	val := 0.0
	for _, tile := range idxActiveTiles {
		for feature := range vf.Weights {
			val += vf.Weights[feature][tile]
		}
	}
	return val*/
}

func learn(state State, action string, target float64, vf *ValueFunction) {

	//fmt.Println("Obtuve un target de: ", target)
	activeTiles := state.GetActiveTiles(action)

	estimations := make([]float64, vf.Features)

	for feature := 0; feature < vf.Features; feature++ {
		for idx := 0; idx < vf.Tilings; idx++ {
			estimations[feature] += vf.Weights[activeTiles[feature][idx]]
		}
	}
	//fmt.Println("Estimations: ", estimations)
	//fmt.Println("Active Tiles: ", activeTiles)
	delta := make([]float64, len(estimations))

	for idx := 0; idx < len(delta); idx++ {
		delta[idx] = /*vf.Alpha **/ 0.0001 * (target - estimations[idx])
	}
	//fmt.Println("Deltas: ", delta)
	for feature := range delta {
		for tile := range activeTiles[feature] {
			vf.Weights[activeTiles[feature][tile]] += delta[feature]
		}
	}
}

func getIdxMax(slice []float64) int {
	idx := 0
	max := slice[idx]
	//get the idx of the biggest element
	for i := 1; i < len(slice); i++ {
		if max < slice[i] {
			idx = i
			max = slice[i]
		}
		//If max and slice are equal, we randomly change so have more exploration in the algorithm
		if max == slice[i] {
			if rand.Float64() <= 0.5 {
				idx = i
				max = slice[i]
			}
		}
	}
	return idx
}
