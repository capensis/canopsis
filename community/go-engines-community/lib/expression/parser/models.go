package parser

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

type Expression struct {
	Or []*OrCondition `@@ { "OR" @@ }`
}

func (e *Expression) Query() bson.M {
	or := make([]bson.M, len(e.Or))
	for i, v := range e.Or {
		or[i] = v.Query()
	}

	if len(or) == 1 {
		return or[0]
	}

	return bson.M{"$or": or}
}

type OrCondition struct {
	And []*Condition `@@ { "AND" @@ }`
}

func (c *OrCondition) Query() bson.M {
	and := make([]bson.M, len(c.And))
	for i, v := range c.And {
		and[i] = v.Query()
	}

	if len(and) == 1 {
		return and[0]
	}

	return bson.M{"$and": and}
}

type Condition struct {
	Operand *ConditionOperand `  @@`
	Not     *Condition        `| "NOT" @@`
}

func (c *Condition) Query() bson.M {
	if c.Operand != nil {
		return c.Operand.Query()
	}
	if c.Not != nil {
		return bson.M{"$not": c.Not.Query()}
	}

	return nil
}

type ConditionOperand struct {
	Operand      *Operand      `@@`
	ConditionRHS *ConditionRHS `[ @@ ]`
}

func (o *ConditionOperand) Query() bson.M {
	left := ""
	var right interface{}
	if o.Operand != nil {
		left, _ = o.Operand.Val().(string)
	}
	if o.ConditionRHS != nil {
		right = o.ConditionRHS.Query()
	}

	return bson.M{left: right}
}

type ConditionRHS struct {
	Compare     *Compare     `  @@`
	Like        *Like        `| "LIKE" @@`
	NotLike     *NotLike     `| "NOT" "LIKE" @@`
	Contains    *Contains    `| "CONTAINS" @@`
	NotContains *NotContains `| "NOT" "CONTAINS" @@`
}

func (r *ConditionRHS) Query() bson.M {
	if r.Compare != nil {
		return r.Compare.Query()
	}
	if r.Like != nil {
		return r.Like.Query()
	}
	if r.NotLike != nil {
		return r.NotLike.Query()
	}
	if r.Contains != nil {
		return r.Contains.Query()
	}
	if r.NotContains != nil {
		return r.NotContains.Query()
	}

	return nil
}

type Compare struct {
	Operator string   `@( "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  *Operand `(  @@`
	Select   *Operand ` | @@ )`
}

func (c *Compare) Query() bson.M {
	var operand interface{}
	if c.Operand != nil {
		operand = c.Operand.Val()
	}
	if c.Select != nil {
		operand = c.Select.Val()
	}

	mapOperator := map[string]string{
		"<=": "$lte",
		"<":  "$lt",
		"=":  "$eq",
		"!=": "$ne",
		">=": "$gte",
		">":  "$gt",
	}

	return bson.M{mapOperator[c.Operator]: operand}
}

type Like struct {
	Operand *Operand `@@`
}

func (l *Like) Query() bson.M {
	var operand interface{}
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	return bson.M{"$regex": operand}
}

type NotLike struct {
	Operand *Operand `@@`
}

func (l *NotLike) Query() bson.M {
	var operand interface{}
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	return bson.M{"$not": bson.M{"$regex": operand}}
}

type Contains struct {
	Operand *Operand `@@`
}

func (l *Contains) Query() bson.M {
	var operand interface{}
	if l.Operand != nil {
		operand = l.Operand.Val()
		if reflect.TypeOf(operand).Kind() != reflect.Array {
			operand = []interface{}{operand}
		}
	}

	return bson.M{"$in": operand}
}

type NotContains struct {
	Operand *Operand `@@`
}

func (l *NotContains) Query() bson.M {
	var operand interface{}
	if l.Operand != nil {
		operand = l.Operand.Val()
		if reflect.TypeOf(operand).Kind() != reflect.Array {
			operand = []interface{}{operand}
		}
	}

	return bson.M{"$nin": operand}
}

type Operand struct {
	Terms []*Term `@@ { @@ }`
}

func (o *Operand) Val() interface{} {
	terms := make([]interface{}, len(o.Terms))
	for i, v := range o.Terms {
		terms[i] = v.Val()
	}

	if len(terms) == 1 {
		return terms[0]
	}

	return terms
}

type Term struct {
	Name        *string  `(@Ident`
	NumberFloat *float64 ` | @Float`
	NumberInt   *int     ` | @Int`
	Str         *string  ` | @String`
	Boolean     *Boolean ` | @("TRUE" | "FALSE")`
	Null        bool     ` | @"NULL" )`
}

func (t *Term) Val() interface{} {
	if t.Name != nil {
		return *t.Name
	}
	if t.NumberFloat != nil {
		return *t.NumberFloat
	}
	if t.NumberInt != nil {
		return *t.NumberInt
	}
	if t.Str != nil {
		return *t.Str
	}
	if t.Boolean != nil {
		return bool(*t.Boolean)
	}
	if t.Null {
		return nil
	}

	return nil
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}
