package eventfilter

import (
	"maps"
	"slices"
	"text/template/parse"

	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
)

// templateExternalDataRefs follows field and index expressions through template scopes.
// A nil path is unknown; an empty, non-nil path identifies the root template data.
func templateExternalDataRefs(tpl libtemplate.ParsedTemplate, refs []string) []string {
	if tpl.Tpl == nil || tpl.Tpl.Tree == nil {
		return nil
	}

	found := make(map[string]bool)
	mark := func(path []string) {
		if len(path) > 1 && path[0] == "ExternalData" {
			if path[1] == "*" {
				for _, ref := range refs {
					found[ref] = true
				}
			} else {
				found[path[1]] = true
			}
		}
	}
	appendPath := func(path, fields []string) []string {
		if path == nil {
			return nil
		}

		return append(slices.Clone(path), fields...)
	}
	var resolve func(parse.Node, []string, map[string][]string) []string
	resolve = func(node parse.Node, dot []string, vars map[string][]string) []string {
		var path []string
		switch n := node.(type) {
		case *parse.DotNode:
			path = dot
		case *parse.FieldNode:
			path = appendPath(dot, n.Ident)
		case *parse.VariableNode:
			path = appendPath(vars[n.Ident[0]], n.Ident[1:])
		case *parse.ChainNode:
			path = appendPath(resolve(n.Node, dot, vars), n.Field)
		case *parse.CommandNode:
			args := make([][]string, len(n.Args))
			for i, arg := range n.Args {
				args[i] = resolve(arg, dot, vars)
			}
			if len(n.Args) == 1 {
				path = args[0]
			} else if name, ok := n.Args[0].(*parse.IdentifierNode); ok && name.Ident == "index" {
				path = args[1]
				for _, arg := range n.Args[2:] {
					key, ok := arg.(*parse.StringNode)
					if !ok {
						break
					}
					path = appendPath(path, []string{key.Text})
				}
			}
		case *parse.PipeNode:
			if n == nil {
				return nil
			}
			for _, cmd := range n.Cmds {
				path = resolve(cmd, dot, vars)
			}
			for _, decl := range n.Decl {
				vars[decl.Ident[0]] = path
			}
		}
		mark(path)

		return path
	}

	visiting := make(map[string]bool)
	var walk func(parse.Node, []string, map[string][]string)
	walk = func(node parse.Node, dot []string, vars map[string][]string) {
		switch n := node.(type) {
		case *parse.ListNode:
			if n != nil {
				for _, child := range n.Nodes {
					walk(child, dot, vars)
				}
			}
		case *parse.ActionNode:
			resolve(n.Pipe, dot, vars)
		case *parse.IfNode:
			scope := maps.Clone(vars)
			resolve(n.Pipe, dot, scope)
			walk(n.List, dot, maps.Clone(scope))
			walk(n.ElseList, dot, scope)
		case *parse.WithNode:
			scope := maps.Clone(vars)
			path := resolve(n.Pipe, dot, scope)
			walk(n.List, path, maps.Clone(scope))
			walk(n.ElseList, dot, scope)
		case *parse.RangeNode:
			scope := maps.Clone(vars)
			path := resolve(n.Pipe, dot, scope)
			item := appendPath(path, []string{"*"})
			for i, decl := range n.Pipe.Decl {
				if len(n.Pipe.Decl) == 2 && i == 0 {
					scope[decl.Ident[0]] = nil
				} else {
					scope[decl.Ident[0]] = item
				}
			}
			walk(n.List, item, scope)
			walk(n.ElseList, dot, maps.Clone(vars))
		case *parse.TemplateNode:
			path := resolve(n.Pipe, dot, vars)
			if nested := tpl.Tpl.Lookup(n.Name); nested != nil && nested.Tree != nil && !visiting[n.Name] {
				visiting[n.Name] = true
				walk(nested.Tree.Root, path, map[string][]string{"$": path})
				delete(visiting, n.Name)
			}
		}
	}
	walk(tpl.Tpl.Tree.Root, []string{}, map[string][]string{"$": {}})

	var result []string
	for _, ref := range refs {
		if found[ref] {
			result = append(result, ref)
		}
	}

	return result
}
