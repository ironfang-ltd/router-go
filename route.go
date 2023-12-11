package router

type routeParams struct {
	Keys   []string
	Values []string
}

func (r *routeParams) get(key string) string {

	for i, k := range r.Keys {
		if k == key {
			return r.Values[i]
		}
	}

	return ""
}

func (r *routeParams) set(key, value string) {

	for i, k := range r.Keys {
		if k == key {
			r.Values[i] = value
			return
		}
	}

	r.Keys = append(r.Keys, key)
	r.Values = append(r.Values, value)
}
