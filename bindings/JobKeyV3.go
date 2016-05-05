package bindings

type JobKeyV3 struct {
	*KeyV3
	/* INHERITED: Name (string representation) for this Key.
	Name string: "" `json:"name"`
	*/

	/* INHERITED: Name (string representation) for the type of Keyed this Key points to.
	Type string: "" `json:"type"`
	*/

	/* INHERITED: URL for the resource that this Key points to, if one exists.
	Url string: "" `json:"URL"`
	*/
}

func newJobKeyV3() *JobKeyV3 {
	return &JobKeyV3{
		KeyV3: &KeyV3{
			Name: "",
			Type: "",
			Url:  "",
		},
	}
}
