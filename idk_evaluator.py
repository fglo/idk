from typing import Optional
from idk_consts import *
from idk_types import Expression, KeywordExpression, OperatorExpression 

AST:list[Expression] = []
EVALUATED_AST:list[Expression] = []

class EvaluatorError(Exception):
    pass

def evaluator_error(line_index, message):
    raise EvaluatorError('Evaluator error in line %d: %s' % (line_index, message))

def try_evaluate_variable(expr:Expression, scope):
    for expr_in_scope in reversed(scope):
        if expr_in_scope.type == TOKEN_OPERATOR and expr_in_scope.value == OPERATOR_ASSIGMENT and expr_in_scope.args[0].value == expr.value:
            return (True, expr_in_scope.args[1])
    return (False, None)

def evaluate_variable(expr:Expression, scope) -> Expression:
    success, value = try_evaluate_variable(expr, scope)
    if not success:
        evaluator_error(expr.line_index, "Could not find assignment for variable %s" % expr.value)
    expr = value 
    return expr

def evaluate_global_variable(expr:Expression) -> Expression:
    success, value = try_evaluate_variable(expr, EVALUATED_AST)
    if not success:
        evaluator_error(expr.line_index, "Could not find assignment for variable %s" % expr.value)
    expr = value 
    return expr
        
def evaluate_side_of_operator(expr:Expression, scope):
    if expr.type == TOKEN_INT or expr.type == TOKEN_CHAR or expr.type == TOKEN_BOOL:
        return expr
    if expr.type == TOKEN_WORD:
        return evaluate_variable(expr, scope)
    elif expr.type == TOKEN_OPERATOR:
        return evaluate_operator(expr, scope)
    evaluator_error(expr.line_index, "Only literals, words and operations are allowed as a side of a operations")
        
def evaluate_left_and_right_operator_args(args:list[Expression], scope) -> tuple[Expression, Expression]:       
    left = evaluate_side_of_operator(args[0], scope)
    right = evaluate_side_of_operator(args[1], scope)
    return left, right      

def evaluate_operator_assignment(expr:OperatorExpression, scope):
    left = expr.args[0]
    right = evaluate_side_of_operator(expr.args[1], scope)
    return OperatorExpression(TOKEN_OPERATOR, expr.value, expr.line_index, [left, right])
    
def evaluate_operator_plus(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_INT, left.value + right.value, expr.line_index)

def evaluate_operator_minus(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_INT, left.value - right.value, expr.line_index)

def evaluate_operator_multiplication(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_INT, left.value * right.value, expr.line_index)

def evaluate_operator_division(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_INT, int(left.value / right.value), expr.line_index)

def evaluate_operator_eq(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value == right.value), expr.line_index)

def evaluate_operator_gt(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value > right.value), expr.line_index)

def evaluate_operator_gte(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value >= right.value), expr.line_index)

def evaluate_operator_lt(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value < right.value), expr.line_index)

def evaluate_operator_lte(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value <= right.value), expr.line_index)

def evaluate_operator_xor(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value != right.value), expr.line_index)

def evaluate_operator_or(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value or right.value), expr.line_index)

def evaluate_operator_and(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    return Expression(TOKEN_BOOL, int(left.value and right.value), expr.line_index)

def evaluate_operator_not(expr:OperatorExpression, scope):        
    right = evaluate_side_of_operator(expr.args[0], scope)
    return Expression(TOKEN_BOOL, int(not right.value), expr.line_index)

def evaluate_operator_in(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    
    curr_value = right
    while left.value != curr_value.value:
        if curr_value.__class__ == OperatorExpression:
            curr_value = curr_value.args[0]
        else:
            return Expression(TOKEN_BOOL, 0, expr.line_index)
    return Expression(TOKEN_BOOL, 1, expr.line_index)

def evaluate_operator_range(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)
    
    if left.value < right.value:
        return Expression(TOKEN_ARRAY, list(range(left.value, right.value)), expr.line_index)
    elif left.value > right.value:
        return Expression(TOKEN_ARRAY, list(reversed(list(range(right.value + 1, left.value + 1)))), expr.line_index)

    evaluator_error(expr.line_index, "Arguments of the exclusive range operator cannot be equal.")

def evaluate_operator_range_inclusive(expr:OperatorExpression, scope):
    left, right = evaluate_left_and_right_operator_args(expr.args, scope)

    if left.value < right.value:
        return Expression(TOKEN_ARRAY, list(range(left.value, right.value + 1)), expr.line_index)
    elif left.value > right.value:
        return Expression(TOKEN_ARRAY, list(reversed(list(range(right.value, left.value + 1)))), expr.line_index)
    else:
        return Expression(TOKEN_ARRAY, [left.value], expr.line_index)

def evaluate_operator(expr:OperatorExpression, scope):
    token_operator = expr.value
    if token_operator == OPERATOR_ASSIGMENT:
        expr = evaluate_operator_assignment(expr, scope)
    elif token_operator == OPERATOR_PLUS:
        expr = evaluate_operator_plus(expr, scope)
    elif token_operator == OPERATOR_MINUS:
        expr = evaluate_operator_minus(expr, scope)
    elif token_operator == OPERATOR_ASTERISK:
        expr = evaluate_operator_multiplication(expr, scope)
    elif token_operator == OPERATOR_SLASH:
        expr = evaluate_operator_division(expr, scope)
    elif token_operator == OPERATOR_EQ:
        expr = evaluate_operator_eq(expr, scope)
    elif token_operator == OPERATOR_GT:
        expr = evaluate_operator_gt(expr, scope)
    elif token_operator == OPERATOR_GTE:
        expr = evaluate_operator_gte(expr, scope)
    elif token_operator == OPERATOR_LT:
        expr = evaluate_operator_lt(expr, scope)
    elif token_operator == OPERATOR_LTE:
        expr = evaluate_operator_lte(expr, scope)
    elif token_operator == OPERATOR_XOR:
        expr = evaluate_operator_xor(expr, scope)
    elif token_operator == OPERATOR_OR:
        expr = evaluate_operator_or(expr, scope)
    elif token_operator == OPERATOR_AND:
        expr = evaluate_operator_and(expr, scope)
    elif token_operator == OPERATOR_NOT:
        expr = evaluate_operator_not(expr, scope)
    elif token_operator == OPERATOR_IN:
        expr = evaluate_operator_in(expr, scope)
    elif token_operator == OPERATOR_RANGE:
        expr = evaluate_operator_range(expr, scope)
    elif token_operator == OPERATOR_RANGE_INCLUSIVE:
        expr = evaluate_operator_range_inclusive(expr, scope)
    else:
        evaluator_error(expr.line_index, "Unknown operator")

    return expr

def evaluate_keyword_action_expression(expr:Expression, scope) -> Expression:
    if expr.type == TOKEN_INT or expr.type == TOKEN_CHAR or expr.type == TOKEN_BOOL:
        expr = expr
    elif expr.type == TOKEN_OPERATOR:
        expr = evaluate_operator(expr, scope)
    elif expr.type == TOKEN_WORD:
        expr = evaluate_variable(expr, scope)
    else:
        evaluator_error(expr.line_index, "Unknown token")
    return expr

def evaluate_keyword_print(expr:KeywordExpression, scope):
    token_to_print = evaluate_keyword_action_expression(expr.keyword_expr, scope)
    if token_to_print.type == TOKEN_INT:
        print(token_to_print.value)
    elif token_to_print.type == TOKEN_CHAR:
        print(chr(token_to_print.value))
    elif token_to_print.type == TOKEN_BOOL:
        print(f'{token_to_print.value == 1}'.lower())
    elif token_to_print.type == TOKEN_ARRAY:
        curr_value = token_to_print
        print(f'{curr_value.value}'.replace('[', '').replace(']', ''))
    else:
        evaluator_error(expr.line_index, "Unsupported variable type")

def evaluate_keyword_if(expr:KeywordExpression, scope):
    comparison_result = evaluate_keyword_action_expression(expr.keyword_expr, scope)
    if comparison_result.value == True or comparison_result.value == 1:
        evaluate_expr_list(expr.prim, scope)
    elif len(expr.alt) > 0 and comparison_result.value == False or comparison_result.value == 0:
        evaluate_expr_list(expr.alt, scope)

def evaluate_keyword_for(expr:KeywordExpression, scope):
    expr_result = evaluate_keyword_action_expression(expr.keyword_expr, scope)
    curr_value = expr_result
    expr_list = [OperatorExpression(TOKEN_OPERATOR, OPERATOR_ASSIGMENT, expr.line_index, [Expression(TOKEN_WORD, '_it', expr.line_index), Expression(TOKEN_INT, curr_value.value, expr.line_index)])] + expr.args[1]
    while curr_value[0] == TOKEN_ARRAY:
        evaluate_expr_list(expr_list, scope)
        if len(curr_value) > 2:
            curr_value = curr_value[2]
    expr_list = [OperatorExpression(TOKEN_OPERATOR, OPERATOR_ASSIGMENT, expr.line_index, [Expression(TOKEN_WORD, '_it', expr.line_index), Expression(TOKEN_INT, curr_value.value, expr.line_index)])] + expr.args[1:]
    evaluate_expr_list(expr_list, scope)

def evaluate_keyword_action(expr:KeywordExpression, scope):
    if expr.value == KEYWORD_PRINT:
        evaluate_keyword_print(expr, scope)
    elif expr.value == KEYWORD_IF:
        evaluate_keyword_if(expr, scope)
    elif expr.value == KEYWORD_FOR:
        evaluate_keyword_for(expr, scope)
    else:
        evaluator_error(expr.line_index, "Unknown keyword")
    return expr

def evaluate_expr_list(expr_list, scope):
    if isinstance(expr_list, list):
        for internal_expr in expr_list:
            evaulate_single_expr(internal_expr, scope + expr_list)
    else:
        evaulate_single_expr(expr_list, scope)

def evaulate_single_expr(expr:Expression, scope) -> Expression:
    if expr.type == TOKEN_KEYWORD:
        expr = evaluate_keyword_action(expr, scope)    
    elif expr.type == TOKEN_OPERATOR:
        expr = evaluate_operator(expr, scope)
    else:
        evaluator_error(expr.line_index, 'Not allowed word on the beginning of the line.')
    return expr

def evaluate(ast:list[Expression]):
    global AST
    global EVALUATED_AST
    EVALUATED_AST = []
    AST = ast
    for expr in AST:
        EVALUATED_AST.append(evaulate_single_expr(expr, EVALUATED_AST))
    return EVALUATED_AST