from idk_consts import * 

AST:list = []
EVALUATED_AST:list = []

class EvaluatorError(Exception):
    pass

def evaluator_error(line_index, message):
    raise EvaluatorError('Error in line %d: %s' % (line_index, message))

def evaluate_global_variable(expr_var, line_index):
    for type, value, *exprs in EVALUATED_AST:
        if value == OPERATOR_ASSIGMENT and exprs[0][1] == expr_var[1]:
            return exprs[1]
    evaluator_error(line_index, "Could not find assignment for variable %s" % expr_var[1])
        
def evaluate_side_of_operator(expr, line_index):
    if expr[0] == TOKEN_INT or expr[0] == TOKEN_CHAR or expr[0] == TOKEN_BOOL:
        return expr
    if expr[0] == TOKEN_WORD:
        return evaluate_global_variable(expr, line_index)
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

def evaluate_operation_eq(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] == right[1]))

def evaluate_operation_gt(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] > right[1]))

def evaluate_operation_gte(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] >= right[1]))

def evaluate_operation_lt(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] < right[1]))

def evaluate_operation_lte(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] <= right[1]))

def evaluate_operation_xor(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] != right[1]))

def evaluate_operation_or(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] or right[1]))

def evaluate_operation_and(expr, line_index):        
    left = evaluate_side_of_operator(expr[2], line_index)
    right = evaluate_side_of_operator(expr[3], line_index)
    return (TOKEN_BOOL, int(left[1] and right[1]))

def evaluate_operation_not(expr, line_index):        
    right = evaluate_side_of_operator(expr[2], line_index)
    return (TOKEN_BOOL, int(not right[1]))

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
    elif token_operator == OPERATOR_EQ:
        expr = evaluate_operation_eq(expr, line_index)
    elif token_operator == OPERATOR_GT:
        expr = evaluate_operation_gt(expr, line_index)
    elif token_operator == OPERATOR_GTE:
        expr = evaluate_operation_gte(expr, line_index)
    elif token_operator == OPERATOR_LT:
        expr = evaluate_operation_lt(expr, line_index)
    elif token_operator == OPERATOR_LTE:
        expr = evaluate_operation_lte(expr, line_index)
    elif token_operator == OPERATOR_XOR:
        expr = evaluate_operation_xor(expr, line_index)
    elif token_operator == OPERATOR_OR:
        expr = evaluate_operation_or(expr, line_index)
    elif token_operator == OPERATOR_AND:
        expr = evaluate_operation_and(expr, line_index)
    elif token_operator == OPERATOR_NOT:
        expr = evaluate_operation_not(expr, line_index)
    else:
        evaluator_error(line_index, "Unknown operator")

    return expr

def evaluate_keyword_action_expression(expr, line_index):
    if expr[0] == TOKEN_INT or expr[0] == TOKEN_CHAR or expr[0] == TOKEN_BOOL:
        expr = expr
    elif expr[0] == TOKEN_OPERATOR:
        expr = evaluate_operation(expr, line_index)
    elif expr[0] == TOKEN_WORD:
        expr = evaluate_global_variable(expr, line_index)
    else:
        evaluator_error(line_index, "Unknown token")
    return expr

def evaluate_keyword_print(expr, line_index):
    token_to_print = evaluate_keyword_action_expression(expr[2], line_index)
    if token_to_print[0] == TOKEN_INT:
        print(token_to_print[1])
    elif token_to_print[0] == TOKEN_CHAR:
        print(chr(token_to_print[1]))
    elif token_to_print[0] == TOKEN_BOOL:
        print(f'{token_to_print[1] == 1}'.lower())
    else:
        evaluator_error(line_index, "Unknown literal type")

def evaluate_keyword_if(expr, line_index):
    comparison_result = evaluate_keyword_action_expression(expr[2], line_index)
    if comparison_result[1] == True or comparison_result[1] == 1:
        evaluate_expr_list(expr[3], line_index)
    elif len(expr) > 4 and comparison_result[1] == False or comparison_result[1] == 0:
        evaluate_expr_list(expr[4], line_index)

def evaluate_keyword_action(expr, line_index):
    if expr[1] == KEYWORD_PRINT:
        evaluate_keyword_print(expr, line_index)
    elif expr[1] == KEYWORD_IF:
        evaluate_keyword_if(expr, line_index)
    else:
        evaluator_error(line_index, "Unknown keyword")
    return expr

def evaluate_expr_list(expr_list, line_index):
    if isinstance(expr_list, list):
        for internal_expr in expr_list:
            evaulate_single_expr(internal_expr, line_index)
    else:
        evaulate_single_expr(expr_list, line_index)

def evaulate_single_expr(expr, line_index):
    if isinstance(expr, list):
        expr = expr[0]
        
    if expr[0] == TOKEN_WORD and expr[1][0] == TOKEN_WORD:
        evaluator_error(line_index, 'Unknown construction.')
    
    if expr[0] == TOKEN_WORD and expr[1][0] == TOKEN_KEYWORD:
        evaluator_error(line_index, 'Unknown construction.')

    if expr[0] == TOKEN_KEYWORD:
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
        EVALUATED_AST.append(evaulate_single_expr(expr, line_index))
    return EVALUATED_AST