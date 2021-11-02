package server

type Resource struct {
	Name        string
	Description string
	Allocated   float64
	Used        float64
}

func (res Resource) Usage() (Allocated, Used, Available float64) {
	return Allocated, Used, Allocated - Used
}
