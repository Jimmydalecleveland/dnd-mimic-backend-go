package datasources

type ItemResolver struct {
	result interface{}
}

func (r *ItemResolver) ToQuantifiedWeapon() (*QuantifiedWeapon, bool) {
	res, ok := r.result.(*QuantifiedWeapon)
	return res, ok
}

func (r *ItemResolver) ToQuantifiedArmor() (*QuantifiedArmor, bool) {
	res, ok := r.result.(*QuantifiedArmor)
	return res, ok
}

func (r *ItemResolver) ToQuantifiedAdventuringGear() (*QuantifiedAdventuringGear, bool) {
	res, ok := r.result.(*QuantifiedAdventuringGear)
	return res, ok
}
