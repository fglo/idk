from idk_consts import * 

AST:list = []
EVALUATED_AST:list = []

class EvaluatorError(Exception):
    pass

def evaluator_error(line_index, message):
    raise EvaluatorError('Error in line %d: %s' % (line_index, message))

def evaluate_variable(exprvariable, line_index):
    for type, value, *exprs in EVALUATED_AST:
        if value == OPERATOR_ASSIGMENT and exprs[0][1] == exprvariable[1]:
            return exprs[1]
    evaluator_error(line_index, "Could not find assignment for variable %s" % exprvariable[1])
        
def evaluate_side_of_operator(expr, line_index):
    if expr[0] == TOKEN_INT or expr[0] == TOKEN_CHAR:
        return expr
    if expr[0] == TOKEN_WORD:
        return evaluate_variable(expr, line_index)
    elif expr[0] == TOKEN_OPERATOR:
        return evaluate_operation(expr, line_index)
    evaluator_error(line_index, "Only literals, words and operations are allowed as a side of a operations")
        
def evaluate_operation_assignment(expr, line_index):        
    left = expr[2] 
    right = evaluate_side_of_operator(expr[3], line_index)
    __, value, *__ = expr
    return (TOKEN_OPERATOR, value, left, right)
    
def evaluate_operation_plus(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_INT, left[1] + right[1])

def evaluate_operation_minus(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_INT, left[1] - right[1])

def evaluate_operation_multiplication(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_INT, left[1] * right[1])

def evaluate_operation_division(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_INT, int(left[1] / right[1]))

def evaluate_operation(expr, line_index):
    token_operator = expr[1]
    if token_operator == OPERATOR_ASSIGMENT:
        expr = evaluate_operation_assignment(expr, line_index)
    elif token_operator == OPERATOR_PLUS:
        expr = evaluate_operation_plus(expr, line_index)
    elif token_operator == OPERATOR_MINUS:
        expr = evaluate_operation_minus(expr, line_index)
    elif token_operator == OPERATOR_MULTIPLICATION:
        expr = evaluate_operation_multiplication(expr, line_index)
    elif token_operator == OPERATOR_DIVISION:
        expr = evaluate_operation_division(expr, line_index)
    else:
        evaluator_error(line_index, "Unknown operator")

    return expr

def evaluate_keyword_action_expression(expr, line_index):
    if expr[0] == TOKEN_INT or expr[0] == TOKEN_CHAR:
        expr = expr
    elif expr[0] == TOKEN_OPERATOR:
        expr = evaluate_operation(expr, line_index)
    elif expr[0] == TOKEN_WORD:
        expr = evaluate_variable(expr, line_index)
    else:
        evaluator_error(line_index, "Unknown token")
    return expr

def evaluate_keyword_print(expr, line_index):
    token_to_print = evaluate_keyword_action_expression(expr[2], line_index)
    if token_to_print[0] == TOKEN_INT:
        print(token_to_print[1])
    elif token_to_print[0] == TOKEN_CHAR:
        print(chr(token_to_print[1]))
    else:
        evaluator_error(line_index, "Unknown literal type")

def evaluate_keyword_action(expr, line_index):
    if expr[1] == KEYWORD_PRINT:
        evaluate_keyword_print(expr, line_index)
    else:
        evaluator_error(line_index, "Unknown keyword")
    return expr

def evaulate_expr(expr, line_index):
    if expr[0] == TOKEN_WORD and expr[1][0] == TOKEN_WORD:
        evaluator_error(line_index, 'Unknown construction.')
    
    if expr[0] == TOKEN_WORD and expr[1][0] == TOKEN_KEYWORD:
        evaluator_error(line_index, 'Unknown construction.')

    if expr[0] == TOKEN_KEYWORD:
        pass
        expr = evaluate_keyword_action(expr, line_index)    
    elif expr[0] == TOKEN_OPERATOR:
        expr = evaluate_operation(expr, line_index)
    else:
        evaluator_error(line_index, 'Not allowed word on the beginning of the line.')
    return expr

def evaluate(ast):
    global AST
    AST = ast
    line_index = 0
    for expr in AST:
        line_index += 1
        EVALUATED_AST.append(evaulate_expr(expr, line_index))
    return EVALUATED_AST