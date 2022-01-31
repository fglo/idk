from idk_consts import * 

AST:list = []
EVALUATED_AST:list = []

class EvaluatorError(Exception):
    pass

def evaluator_error(line_index, message):
    raise EvaluatorError('Evaluator error in line %d: %s' % (line_index, message))

def try_evaluate_variable(token_var, scope):
    for type, value, *exprs in reversed(scope):
        if type == TOKEN_OPERATOR and value == OPERATOR_ASSIGMENT and exprs[0][1] == token_var[1]:
            return (True, exprs[1])
    return (False, -1)

def evaluate_variable(token_var, line_index, scope):
    success, value = try_evaluate_variable(token_var, scope)
    if success: 
        return value
    else:
        evaluator_error(line_index, "Could not find assignment for variable %s" % token_var[1])

def evaluate_global_variable(token_var, line_index):
    success, value = try_evaluate_variable(token_var, EVALUATED_AST)
    if success: 
        return value
    else:
        evaluator_error(line_index, "Could not find assignment for variable %s" % token_var[1])
        
def evaluate_side_of_operator(expr, line_index, scope):
    if expr[0] == TOKEN_INT or expr[0] == TOKEN_CHAR or expr[0] == TOKEN_BOOL:
        return expr
    if expr[0] == TOKEN_WORD:
        return evaluate_variable(expr, line_index, scope)
    elif expr[0] == TOKEN_OPERATOR:
        return evaluate_operator(expr, scope)
    evaluator_error(line_index, "Only literals, words and operations are allowed as a side of a operations")
        
def evaluate_operator_assignment(expr, scope):        
    left = expr[2] 
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    __, value, *__ = expr
    return (TOKEN_OPERATOR, value, left, right)
    
def evaluate_operator_plus(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_INT, left[1] + right[1])

def evaluate_operator_minus(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_INT, left[1] - right[1])

def evaluate_operator_multiplication(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_INT, left[1] * right[1])

def evaluate_operator_division(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_INT, int(left[1] / right[1]))

def evaluate_operator_eq(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] == right[1]))

def evaluate_operator_gt(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] > right[1]))

def evaluate_operator_gte(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] >= right[1]))

def evaluate_operator_lt(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] < right[1]))

def evaluate_operator_lte(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] <= right[1]))

def evaluate_operator_xor(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] != right[1]))

def evaluate_operator_or(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] or right[1]))

def evaluate_operator_and(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    return (TOKEN_BOOL, int(left[1] and right[1]))

def evaluate_operator_not(expr, scope):        
    right = evaluate_side_of_operator(expr[2], expr[-1], scope)
    return (TOKEN_BOOL, int(not right[1]))

def evaluate_operator_in(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    
    curr_value = right
    while left[1] != curr_value[1]:
        if len(curr_value) > 2:
            curr_value = curr_value[2]
        else:
            return (TOKEN_BOOL, 0)
    return (TOKEN_BOOL, 1)

def evaluate_operator_range(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    
    left_val = left[1]
    right_val = right[1]
    if left_val < right_val:
        if left_val < right_val - 1:
            return (TOKEN_ARRAY, left_val, evaluate_operator_range((expr[0], expr[1], (TOKEN_INT, left_val + 1), right), scope))
        if left_val == right_val - 1:
            return (TOKEN_INT, left_val)
    elif left_val > right_val:
        if left_val > right_val + 1:
            return (TOKEN_ARRAY, left_val, evaluate_operator_range((expr[0], expr[1], (TOKEN_INT, left_val - 1), right), scope))
        if left_val == right_val + 1:
            return (TOKEN_INT, left_val)

    evaluator_error(expr[-1], "Arguments of the exclusive range operator cannot be equal.")

def evaluate_operator_range_inclusive(expr, scope):        
    left = evaluate_side_of_operator(expr[2], expr[-1], scope)
    right = evaluate_side_of_operator(expr[3], expr[-1], scope)
    
    left_val = left[1]
    right_val = right[1]
    if left_val < right_val:
        return (TOKEN_ARRAY, left_val, evaluate_operator_range_inclusive((expr[0], expr[1], (TOKEN_INT, left_val + 1), right), scope))
    elif left_val > right_val:
        return (TOKEN_ARRAY, left_val, evaluate_operator_range_inclusive((expr[0], expr[1], (TOKEN_INT, left_val - 1), right), scope))
    else:
        return (TOKEN_INT, left_val)

def evaluate_operator(expr, scope):
    token_operator = expr[1]
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
        evaluator_error(expr[-1], "Unknown operator")

    return expr

def evaluate_keyword_action_expression(expr, line_index, scope):
    if expr[0] == TOKEN_INT or expr[0] == TOKEN_CHAR or expr[0] == TOKEN_BOOL:
        expr = expr
    elif expr[0] == TOKEN_OPERATOR:
        expr = evaluate_operator(expr, scope)
    elif expr[0] == TOKEN_WORD:
        expr = evaluate_variable(expr, line_index, scope)
    else:
        evaluator_error(expr[-1], "Unknown token")
    return expr

def evaluate_keyword_print(expr, scope):
    token_to_print = evaluate_keyword_action_expression(expr[2], expr[-1], scope)
    if token_to_print[0] == TOKEN_INT:
        print(token_to_print[1])
    elif token_to_print[0] == TOKEN_CHAR:
        print(chr(token_to_print[1]))
    elif token_to_print[0] == TOKEN_BOOL:
        print(f'{token_to_print[1] == 1}'.lower())
    elif token_to_print[0] == TOKEN_ARRAY:
        curr_value = token_to_print
        while curr_value[0] == TOKEN_ARRAY:
            print(curr_value[1], end = '')
            print(', ', end = '')
            if len(curr_value) > 2:
                curr_value = curr_value[2]
        print(curr_value[1])
    else:
        evaluator_error(expr[-1], "Unsupported variable type")

def evaluate_keyword_if(expr, scope):
    comparison_result = evaluate_keyword_action_expression(expr[2], expr[-1], scope)
    if comparison_result[1] == True or comparison_result[1] == 1:
        evaluate_expr_list(expr[3], scope)
    elif len(expr) > 4 and comparison_result[1] == False or comparison_result[1] == 0:
        evaluate_expr_list(expr[4], scope)

def evaluate_keyword_for(expr, scope):
    expr_result = evaluate_keyword_action_expression(expr[2], expr[-1], scope)
    curr_value = expr_result
    expr_list = [(TOKEN_OPERATOR, OPERATOR_ASSIGMENT, (TOKEN_WORD, '_it'), (TOKEN_INT, curr_value[1]))] + expr[3]
    while curr_value[0] == TOKEN_ARRAY:
        evaluate_expr_list(expr_list, scope)
        if len(curr_value) > 2:
            curr_value = curr_value[2]
        expr_list = [(TOKEN_OPERATOR, OPERATOR_ASSIGMENT, (TOKEN_WORD, '_it'), (TOKEN_INT, curr_value[1]))] + expr_list[1:]
    evaluate_expr_list(expr_list, scope)

def evaluate_keyword_action(expr, scope):
    if expr[1] == KEYWORD_PRINT:
        evaluate_keyword_print(expr, scope)
    elif expr[1] == KEYWORD_IF:
        evaluate_keyword_if(expr, scope)
    elif expr[1] == KEYWORD_FOR:
        evaluate_keyword_for(expr, scope)
    else:
        evaluator_error(expr[-1], "Unknown keyword")
    return expr

def evaluate_expr_list(expr_list, scope):
    if isinstance(expr_list, list):
        for internal_expr in expr_list:
            evaulate_single_expr(internal_expr, scope + expr_list)
    else:
        evaulate_single_expr(expr_list, scope)

def evaulate_single_expr(expr, scope):
    if isinstance(expr, list):
        expr = expr[0]
        
    if expr[0] == TOKEN_WORD and expr[1][0] == TOKEN_WORD:
        evaluator_error(expr[-1], 'Unknown construction.')
    
    if expr[0] == TOKEN_WORD and expr[1][0] == TOKEN_KEYWORD:
        evaluator_error(expr[-1], 'Unknown construction.')

    if expr[0] == TOKEN_KEYWORD:
        expr = evaluate_keyword_action(expr, scope)    
    elif expr[0] == TOKEN_OPERATOR:
        expr = evaluate_operator(expr, scope)
    else:
        evaluator_error(expr[-1], 'Not allowed word on the beginning of the line.')
    return expr

def evaluate(ast):
    global AST
    AST = ast
    for expr in AST:
        EVALUATED_AST.append(evaulate_single_expr(expr, EVALUATED_AST))
    return EVALUATED_AST