package evaluator

import (
	"monkey_cc/ast"
	"monkey_cc/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Exp)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.PrefixExpression:
		return evalPrefixExp(node.Operator, Eval(node.Right))
	case *ast.Boolean:
		if node.Value {
			return object.TRUE
		} else {
			return object.FALSE
		}
	}
	return nil
}

// for an array of statements, the overall value is
// the value of last statement
// if a statement is not expression statement, it's value
// is Null
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func evalPrefixExp(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return object.NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
