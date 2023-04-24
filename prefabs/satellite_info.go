package prefabs

type Satellite struct{
    NAME  string  `json:"name"`
    NumSat string  `json:"num_sat"`
    Inter  string `json:"inter"`
    Year  string  `json:"year"`
    Day string  `json:"day"`
	FirstMotion  string `json:"first_motion"`
	SecondMotion  string  `json:"second_motion"`
	Drag  string `json:"drag"`
	Number string `json:"number"`
	Incl string `json:"incl"`
	R_A string `json:"r_a"`
	Eccentricity string `json:"eccentricity"`
	ArgPer string `json:"arg_per"`
	Anomaly string 	`json:"anomaly"`
	Motion string `json:"motion"`
	Epoch string `json:"epoch"`
}