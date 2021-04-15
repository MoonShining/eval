package eval

// condition ast
type Condition struct {
	AND []*Condition
	OR  []*Condition

	Left  interface{}
	Op    string
	Right interface{}
}

// ast walk
func (c *Condition) Match(env map[string]interface{}) bool {
	var match bool

	if len(c.AND) > 0 {
		match = true
		for _, andConf := range c.AND {
			if !andConf.Match(env) {
				match = false
				break
			}
		}
	} else if len(c.OR) > 0 {
		match = false
		for _, orConf := range c.OR {
			if orConf.Match(env) {
				match = true
				break
			}
		}
	} else {
		l := env[c.Left.(string)]
		r := c.Right

		switch c.Op {
		case "==":
			match = l == r
		case "!=":
			match = l != r
		}
	}

	return match
}

