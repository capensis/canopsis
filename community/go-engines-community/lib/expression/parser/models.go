//nolint:govet
package parser

import (
	"fmt"
	"reflect"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics/schema"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/bson"
)

const DefaultArgsLen = 6

type Expression struct {
	Or []*OrCondition `@@ { "OR" @@ }`
}

func (e *Expression) PostgresQuery(prefix string) (string, pgx.NamedArgs) {
	orQueries := make([]string, len(e.Or))
	orArgs := make(map[string]any)

	var args map[string]any
	for i, cond := range e.Or {
		orQueries[i], args = cond.PostgresQuery(prefix)

		for k, v := range args {
			orArgs[k] = v
		}
	}

	if len(orQueries) == 1 {
		return "WHERE " + orQueries[0], orArgs
	}

	return "WHERE (" + strings.Join(orQueries, " OR ") + ")", orArgs
}

func (e *Expression) MongoQuery() bson.M {
	or := make([]bson.M, len(e.Or))
	for i, v := range e.Or {
		or[i] = v.MongoQuery()
	}

	if len(or) == 1 {
		return or[0]
	}

	return bson.M{"$or": or}
}

func (e *Expression) MongoExprQuery() bson.M {
	or := make([]bson.M, len(e.Or))
	for i, v := range e.Or {
		or[i] = v.MongoExprQuery()
	}

	if len(or) == 1 {
		return or[0]
	}

	return bson.M{"$or": or}
}

func (e *Expression) GetFields() []string {
	fields := make([]string, 0, len(e.Or))
	for _, v := range e.Or {
		fields = append(fields, v.GetFields()...)
	}

	return fields
}

type OrCondition struct {
	And []*Condition `@@ { "AND" @@ }`
}

func (c *OrCondition) PostgresQuery(prefix string) (string, pgx.NamedArgs) {
	andQueries := make([]string, len(c.And))
	andArgs := make(map[string]any)

	var args map[string]any
	for i, v := range c.And {
		andQueries[i], args = v.PostgresQuery(prefix)

		for k, v := range args {
			andArgs[k] = v
		}
	}

	if len(andQueries) == 1 {
		return andQueries[0], andArgs
	}

	return strings.Join(andQueries, " AND "), andArgs
}

func (c *OrCondition) MongoQuery() bson.M {
	and := make([]bson.M, len(c.And))
	for i, v := range c.And {
		and[i] = v.MongoQuery()
	}

	if len(and) == 1 {
		return and[0]
	}

	return bson.M{"$and": and}
}

func (c *OrCondition) MongoExprQuery() bson.M {
	and := make([]bson.M, len(c.And))
	for i, v := range c.And {
		and[i] = v.MongoExprQuery()
	}

	if len(and) == 1 {
		return and[0]
	}

	return bson.M{"$and": and}
}

func (c *OrCondition) GetFields() []string {
	fields := make([]string, 0, len(c.And))
	for _, v := range c.And {
		fields = append(fields, v.GetFields()...)
	}

	return fields
}

type Condition struct {
	Operand *ConditionOperand `  @@`
	Not     *Condition        `| "NOT" @@`
}

func (c *Condition) PostgresQuery(prefix string) (string, pgx.NamedArgs) {
	if c.Operand != nil {
		return c.Operand.PostgresQuery(prefix)
	}
	if c.Not != nil {
		return c.Not.Operand.NotPostgresQuery(prefix)
	}

	return "", nil
}

func (c *Condition) MongoQuery() bson.M {
	if c.Operand != nil {
		return c.Operand.MongoQuery()
	}
	if c.Not != nil {
		return c.Not.Operand.NotMongoQuery()
	}

	return nil
}

func (c *Condition) MongoExprQuery() bson.M {
	if c.Operand != nil {
		return c.Operand.MongoExprQuery()
	}
	if c.Not != nil {
		return c.Not.Operand.NotExprQuery()
	}

	return nil
}

func (c *Condition) GetFields() []string {
	if c.Operand != nil {
		return c.Operand.GetFields()
	}

	if c.Not != nil {
		return c.Not.Operand.GetFields()
	}

	return nil
}

type ConditionOperand struct {
	Operand      *Operand      `@@`
	ConditionRHS *ConditionRHS `[ @@ ]`
}

func (o *ConditionOperand) PostgresQuery(prefix string) (string, pgx.NamedArgs) {
	left := ""
	var right string
	var args pgx.NamedArgs
	if o.Operand != nil {
		left, _ = o.Operand.Val().(string)
	}

	left = schema.TransformToEntityMetaField(left, prefix)

	if o.ConditionRHS != nil {
		right, args = o.ConditionRHS.PostgresQuery()
		if args != nil && o.ConditionRHS.Compare != nil && o.ConditionRHS.Compare.Operator == "!=" || o.ConditionRHS.NotLike != nil {
			return "(" + left + " IS NULL OR " + left + " " + right + ")", args
		} else if o.ConditionRHS.NotContains != nil {
			return "(" + left + " IS NULL OR NOT " + left + " " + right + ")", args
		}
	}

	return left + " " + right, args
}

func (o *ConditionOperand) MongoQuery() bson.M {
	left := ""
	var right any
	if o.Operand != nil {
		left, _ = o.Operand.Val().(string)
	}
	if o.ConditionRHS != nil {
		right = o.ConditionRHS.MongoQuery()
	}

	return bson.M{left: right}
}

func (o *ConditionOperand) MongoExprQuery() bson.M {
	op := ""
	var res bson.M
	if o.Operand != nil {
		op, _ = o.Operand.Val().(string)
	}
	if o.ConditionRHS != nil {
		res = o.ConditionRHS.ExprQuery(op)
	}

	return res
}

func (o *ConditionOperand) GetFields() []string {
	if o.Operand != nil {
		field, _ := o.Operand.Val().(string)
		return []string{field}
	}

	return nil
}

func (o *ConditionOperand) NotPostgresQuery(prefix string) (string, pgx.NamedArgs) {
	left := ""
	var right string
	var args pgx.NamedArgs
	if o.Operand != nil {
		left, _ = o.Operand.Val().(string)
	}

	left = schema.TransformToEntityMetaField(left, prefix)

	if o.ConditionRHS != nil {
		right, args = o.ConditionRHS.NotPostgresQuery()
		if args != nil && o.ConditionRHS.Compare != nil && o.ConditionRHS.Compare.Operator == "!=" || o.ConditionRHS.Like != nil {
			return "(" + left + " IS NULL OR " + left + " " + right + ")", args
		} else if o.ConditionRHS.Contains != nil {
			return "(" + left + " IS NULL OR NOT " + left + " " + right + ")", args
		}
	}

	return left + " " + right, args
}

func (o *ConditionOperand) NotMongoQuery() bson.M {
	left := ""
	var right any
	if o.Operand != nil {
		left, _ = o.Operand.Val().(string)
	}
	if o.ConditionRHS != nil {
		right = o.ConditionRHS.NotMongoQuery()
	}

	return bson.M{left: right}
}

func (o *ConditionOperand) NotExprQuery() bson.M {
	op := ""
	var res bson.M
	if o.Operand != nil {
		op, _ = o.Operand.Val().(string)
	}
	if o.ConditionRHS != nil {
		res = o.ConditionRHS.NotExprQuery(op)
	}

	return res
}

type ConditionRHS struct {
	Compare     *Compare     `  @@`
	Like        *Like        `| "LIKE" @@`
	NotLike     *NotLike     `| "NOT" "LIKE" @@`
	Contains    *Contains    `| "CONTAINS" @@`
	NotContains *NotContains `| "NOT" "CONTAINS" @@`
}

func (r *ConditionRHS) PostgresQuery() (string, pgx.NamedArgs) {
	if r.Compare != nil {
		return r.Compare.PostgresQuery()
	}
	if r.Like != nil {
		return r.Like.PostgresQuery()
	}
	if r.NotLike != nil {
		return r.NotLike.PostgresQuery()
	}
	if r.Contains != nil {
		return r.Contains.PostgresQuery()
	}
	if r.NotContains != nil {
		return r.NotContains.PostgresQuery()
	}

	return "", nil
}

func (r *ConditionRHS) NotPostgresQuery() (string, pgx.NamedArgs) {
	if r.Compare != nil {
		return r.Compare.NotPostgresQuery()
	}
	if r.Like != nil {
		return r.Like.NotPostgresQuery()
	}
	if r.NotLike != nil {
		return r.NotLike.NotPostgresQuery()
	}
	// There are no NotPostgresQuery() functions for contains and not contains.
	if r.Contains != nil {
		return r.Contains.PostgresQuery()
	}
	if r.NotContains != nil {
		return r.NotContains.PostgresQuery()
	}

	return "", nil
}

func (r *ConditionRHS) MongoQuery() bson.M {
	if r.Compare != nil {
		return r.Compare.MongoQuery()
	}
	if r.Like != nil {
		return r.Like.MongoQuery()
	}
	if r.NotLike != nil {
		return r.NotLike.MongoQuery()
	}
	if r.Contains != nil {
		return r.Contains.MongoQuery()
	}
	if r.NotContains != nil {
		return r.NotContains.MongoQuery()
	}

	return nil
}

func (r *ConditionRHS) ExprQuery(op string) bson.M {
	if r.Compare != nil {
		return r.Compare.ExprQuery(op)
	}
	if r.Like != nil {
		return r.Like.ExprQuery(op)
	}
	if r.NotLike != nil {
		return r.NotLike.ExprQuery(op)
	}
	if r.Contains != nil {
		return r.Contains.ExprQuery(op)
	}
	if r.NotContains != nil {
		return r.NotContains.ExprQuery(op)
	}

	return nil
}

func (r *ConditionRHS) NotMongoQuery() bson.M {
	q := r.MongoQuery()
	if len(q) == 0 {
		return q
	}

	return bson.M{"$not": q}
}

func (r *ConditionRHS) NotExprQuery(op string) bson.M {
	q := r.ExprQuery(op)
	if len(q) == 0 {
		return q
	}

	return bson.M{"$not": q}
}

type Compare struct {
	Operator string   `@( "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  *Operand `(  @@`
	Select   *Operand ` | @@ )`
}

func (c *Compare) NotPostgresQuery() (string, pgx.NamedArgs) {
	// invert operators for NOT statement
	switch c.Operator {
	case "<=":
		c.Operator = ">"
	case "<":
		c.Operator = ">="
	case "=":
		c.Operator = "!="
	case "!=":
		c.Operator = "="
	case ">=":
		c.Operator = "<"
	case ">":
		c.Operator = "<="
	}

	return c.PostgresQuery()
}

func (c *Compare) PostgresQuery() (string, pgx.NamedArgs) {
	var operand any
	if c.Operand != nil {
		operand = c.Operand.Val()
	}
	if c.Select != nil {
		operand = c.Select.Val()
	}

	if operand == nil {
		switch c.Operator {
		case "=":
			return "IS NULL", nil
		case "!=":
			return "IS NOT NULL", nil
		}

		// no default after switch, other operators don't cause any errors with null value.
		// so keep them as is.
	}

	argName := randArgumentName()

	return c.Operator + " @" + argName, pgx.NamedArgs{argName: operand}
}

func (c *Compare) MongoQuery() bson.M {
	var operand any
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

func (c *Compare) ExprQuery(op string) bson.M {
	var operand any
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

	return bson.M{mapOperator[c.Operator]: bson.A{"$" + op, operand}}
}

type Like struct {
	Operand *Operand `@@`
}

func (l *Like) PostgresQuery() (string, pgx.NamedArgs) {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	argName := randArgumentName()

	return "~ @" + argName, pgx.NamedArgs{argName: fmt.Sprintf("%v", operand)}
}

func (l *Like) NotPostgresQuery() (string, pgx.NamedArgs) {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	argName := randArgumentName()

	return "!~ @" + argName, pgx.NamedArgs{argName: fmt.Sprintf("%v", operand)}
}

func (l *Like) MongoQuery() bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	return bson.M{"$regex": fmt.Sprintf("%v", operand)}
}

func (l *Like) ExprQuery(op string) bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	return bson.M{"$regexMatch": bson.M{
		"input": "$" + op,
		"regex": fmt.Sprintf("%v", operand),
	}}
}

type NotLike struct {
	Operand *Operand `@@`
}

func (l *NotLike) NotPostgresQuery() (string, pgx.NamedArgs) {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	argName := randArgumentName()

	return "~ @" + argName, pgx.NamedArgs{argName: fmt.Sprintf("%v", operand)}
}

func (l *NotLike) PostgresQuery() (string, pgx.NamedArgs) {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	argName := randArgumentName()

	return "!~ @" + argName, pgx.NamedArgs{argName: fmt.Sprintf("%v", operand)}
}

func (l *NotLike) MongoQuery() bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	return bson.M{"$not": bson.M{"$regex": fmt.Sprintf("%v", operand)}}
}

func (l *NotLike) ExprQuery(op string) bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
	}

	return bson.M{"$not": bson.M{"$regexMatch": bson.M{
		"input": "$" + op,
		"regex": fmt.Sprintf("%v", operand),
	}}}
}

type Contains struct {
	Operand *Operand `@@`
}

func (l *Contains) MongoQuery() bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
		if reflect.TypeOf(operand).Kind() != reflect.Slice {
			operand = []any{operand}
		}
	}

	return bson.M{"$in": operand}
}

func (l *Contains) PostgresQuery() (string, pgx.NamedArgs) {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
		operandSlice, ok := operand.([]any)
		if !ok {
			operand = []any{fmt.Sprintf("%v", operand)}
		} else {
			for i, v := range operandSlice {
				operandSlice[i] = fmt.Sprintf("%v", v)
			}
		}
	}

	argName := randArgumentName()

	return "= ANY (@" + argName + ")", pgx.NamedArgs{argName: operand}
}

func (l *Contains) ExprQuery(op string) bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
		if reflect.TypeOf(operand).Kind() != reflect.Slice {
			operand = []any{operand}
		}
	}

	return bson.M{"$in": bson.A{"$" + op, operand}}
}

type NotContains struct {
	Operand *Operand `@@`
}

func (l *NotContains) MongoQuery() bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
		if reflect.TypeOf(operand).Kind() != reflect.Slice {
			operand = []any{operand}
		}
	}

	return bson.M{"$nin": operand}
}

func (l *NotContains) PostgresQuery() (string, pgx.NamedArgs) {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
		operandSlice, ok := operand.([]any)
		if !ok {
			operand = []any{fmt.Sprintf("%v", operand)}
		} else {
			for i, v := range operandSlice {
				operandSlice[i] = fmt.Sprintf("%v", v)
			}
		}
	}

	argName := randArgumentName()

	// it's the same as Contains query, it's negated in the caller function above.
	return "= ANY (@" + argName + ")", pgx.NamedArgs{argName: operand}
}

func (l *NotContains) ExprQuery(op string) bson.M {
	var operand any
	if l.Operand != nil {
		operand = l.Operand.Val()
		if reflect.TypeOf(operand).Kind() != reflect.Slice {
			operand = []any{operand}
		}
	}

	return bson.M{"$not": bson.M{"$in": bson.A{"$" + op, operand}}}
}

type Operand struct {
	Terms []*Term `@@ { @@ }`
}

func (o *Operand) Val() any {
	terms := make([]any, len(o.Terms))
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

func (t *Term) Val() any {
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

func randArgumentName() string {
	// a little hack here, postgres named arguments can't start from number, so just add some letter prefix
	return "a" + utils.RandString(DefaultArgsLen)
}
