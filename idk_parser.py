from idk_consts import * 

AST:list = []

class ParserError(Exception):
    pass

def parser_error(line_index, message):
    raise ParserError('Error in line %d: %s' % (line_index, message))

def get_first_operator_index(tokens):
    operator_index = -1
    for index, token in enumerate(tokens):
        if token[0] == TOKEN_OPERATOR:
            operator_index = index
            break
    return operator_index

def get_first_less_important_operator_index(tokens):
    found_operator_index = -1
    found_operator = None
    for index, token in enumerate(tokens):
        if token[0] == TOKEN_OPERATOR:
            if token[1] == OPERATOR_ASSIGMENT:
                found_operator_index = index
                break
            if token[1] == OPERATOR_EQ or token[1] == OPERATOR_GT or token[1] == OPERATOR_LT:
                found_operator_index = index
                break
            if found_operator is None:
                found_operator_index = index
                found_operator = token
            elif found_operator[1] >= OPERATOR_MULTIPLICATION and token[1] <= OPERATOR_MINUS:
                found_operator_index = index
                found_operator = token
                break
    return found_operator_index
        
def handle_side_of_operator(side, line_index, nesting_lvl):
    if len(side) == 1:
        side = side[0]
    else:
        side = handle_operator((side, line_index, nesting_lvl))
    return side
        
def handle_operator_assignment(expr, operator_index):
    if isinstance(expr, list):
        tokens, line_index, nesting_lvl = expr[0]
    else:
        tokens, line_index, nesting_lvl = expr

    left = tokens[0]
    for expr in AST:
        if expr[1] == OPERATOR_ASSIGMENT and expr[2][1] == left[1]:
            parser_error(line_index, "You cannot assign value to a variable that already has been assigned.")
    right = handle_side_of_operator(tokens[operator_index + 1:], line_index, nesting_lvl)
    return (TOKEN_OPERATOR, OPERATOR_ASSIGMENT, left, right)  
    
def handle_operator_with_two_side_evaluation(expr, operator_index):
    tokens, line_index, nesting_lvl = expr        
    left = handle_side_of_operator(tokens[0:operator_index], line_index, nesting_lvl)
    right = handle_side_of_operator(tokens[operator_index + 1:], line_index, nesting_lvl)
    __, operator_type = tokens[operator_index]
    return (TOKEN_OPERATOR, operator_type, left, right)

def handle_operator(expr):
    tokens, line_index, nesting_lvl = expr
    operator_index = get_first_less_important_operator_index(tokens)
    token_operator = tokens[operator_index][1]
    if token_operator == OPERATOR_ASSIGMENT:
        expr = handle_operator_assignment(expr, operator_index)
    elif token_operator == OPERATOR_PLUS:
        expr = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_MINUS:
        expr = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_MULTIPLICATION:
        expr = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_DIVISION:
        expr = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_EQ:
        expr = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_GT:
        expr = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_LT:
        expr = handle_operator_with_two_side_evaluation(expr, operator_index)
    else:
        parser_error(line_index, "Unknown operator")
    return expr
    
def handle_keyword_print(expr):
    if isinstance(expr, list):
        tokens, line_index, nesting_lvl = expr[0]
    else:
        tokens, line_index, nesting_lvl = expr
        
    operator_index = get_first_operator_index(tokens)
    if operator_index > -1:
        expr_value = handle_operator((tokens[1:], line_index, nesting_lvl))
    elif tokens[1][0] <= TOKEN_WORD:
        expr_value = tokens[1]
    else:
        parser_error(line_index, 'You can print only literals, variables or operations.')
    return (TOKEN_KEYWORD, KEYWORD_PRINT, expr_value)

def handle_keyword_else(expr_list):
    tokens, __, nesting_lvl = expr_list[0]
    if len(tokens) > 1 and tokens[1][0] == TOKEN_KEYWORD and tokens[1][1] == KEYWORD_IF:
        del expr_list[0][0][0]
        for line_expr in expr_list:
            line_expr = (line_expr[0], line_expr[1], line_expr[2] + 1)
        return handle_keyword_if(expr_list)
    
    expr_ast = []
    nested_found = 0
    nested_statement = []
    for internal_expr in expr_list[1:]:
        token = internal_expr[0][0]
        if token[0] == TOKEN_KEYWORD and token[1] == KEYWORD_END and internal_expr[2] == nesting_lvl:
            break
        if token[0] == TOKEN_KEYWORD and internal_expr[0][0][1] == KEYWORD_IF and internal_expr[2] == nesting_lvl + 1:
            nested_found += 1
            nested_statement = []
        if token[0] == TOKEN_KEYWORD and token[1] == KEYWORD_END and internal_expr[2] == nesting_lvl + 1:
            nested_found -= 1
            expr_ast.append(nested_statement)
            continue
        if nested_found <= 0:
            expr_ast.append(internal_expr)
        else:
            nested_statement.append(internal_expr)
    
    expr_ast_prim = []
    for expr in expr_ast:
        expr_ast_prim.append(parse_expr(expr))
    
    return expr_ast_prim
    
def handle_keyword_if(expr_list):
    first_line_tokens, line_index, nesting_lvl = expr_list[0]
            
    if first_line_tokens[1][0] == TOKEN_KEYWORD and first_line_tokens[1][1] == KEYWORD_IF:
        first_line_tokens = first_line_tokens[1:]
            
    operator_index = get_first_operator_index(first_line_tokens)
    if operator_index > -1:
        expr_value = handle_operator((first_line_tokens[1:], line_index, nesting_lvl))
    elif first_line_tokens[1][0] <= TOKEN_WORD:
        expr_value = first_line_tokens[1]
    else:
        parser_error(line_index, 'You can do operations only on literals, variables or other operations.')
        
    expr_ast = []
    expr_ast_alt = []
    else_found = False
    nested_found = 0
    nested_statement = []
    for internal_expr_line in expr_list[1:]:
        token = internal_expr_line[0][0]
        if token[0] == TOKEN_KEYWORD and token[1] == KEYWORD_END and internal_expr_line[2] == nesting_lvl:
            break
        if token[0] == TOKEN_KEYWORD and token[1] == KEYWORD_ELSE and internal_expr_line[2] == nesting_lvl:
            else_found = True
        
        if else_found:
            expr_ast_alt.append(internal_expr_line)
        else:
            if token[0] == TOKEN_KEYWORD and token[1] == KEYWORD_IF and internal_expr_line[2] == nesting_lvl + 1:
                nested_found += 1
                nested_statement = []
            if token[0] == TOKEN_KEYWORD and token[1] == KEYWORD_END and internal_expr_line[2] == nesting_lvl + 1:
                nested_found -= 1
                expr_ast.append(nested_statement)
                continue
            if nested_found <= 0:
                expr_ast.append(internal_expr_line)
            else:
                nested_statement.append(internal_expr_line)
    
    expr_ast_prim = []
    for expr in expr_ast:
        expr_ast_prim.append(parse_expr(expr))
    
    if else_found:
        expr_ast_alt = handle_keyword_else(expr_ast_alt)
            
    expr_ast = (TOKEN_KEYWORD, KEYWORD_IF, expr_value, expr_ast_prim, expr_ast_alt)
    return expr_ast
    
def handle_keyword_action(expr):
    if isinstance(expr, list):
        tokens, line_index, __ = expr[0]
    else:
        tokens, line_index, __ = expr
        
    if tokens[0][1] == KEYWORD_PRINT:
        expr_ast = handle_keyword_print(expr)
    elif tokens[0][1] == KEYWORD_IF:
        expr_ast = handle_keyword_if(expr)
    elif tokens[0][1] == KEYWORD_ELSE:
        parser_error(line_index, 'KEYWORD_ELSE is not expected.') 
        expr_ast = handle_keyword_else(expr)
    elif tokens[0][1] == KEYWORD_END:
        parser_error(line_index, 'KEYWORD_END is not expected.') 
    else:
        parser_error(line_index, 'Unknown keyword.') 
    return expr_ast

def parse_expr(expr):
    if isinstance(expr, list):
        first_line_of_expr, line_index, __ = expr[0]
    else:
        first_line_of_expr, line_index, __ = expr

    if first_line_of_expr[0][0] == TOKEN_KEYWORD:
        expr = handle_keyword_action(expr)    
    elif first_line_of_expr[0][0] == TOKEN_WORD and len(first_line_of_expr) > 1 and first_line_of_expr[1][0] == TOKEN_OPERATOR:
        expr = handle_operator(expr[0])
    else:
        parser_error(line_index, 'Not allowed word on the beginning of the line.')
    return expr

def parse(tokenized_code_lines):
    expr_list = []
    curr_expr = []
    multilvl_expr = 0
    for line_index, line_tokens in enumerate(tokenized_code_lines):
        if len(line_tokens) == 0:
            continue
    
        if line_tokens[0][0] == TOKEN_WORD and line_tokens[1][0] == TOKEN_WORD:
            parser_error(line_index, 'Unknown construction.')
        
        if line_tokens[0][0] == TOKEN_WORD and line_tokens[1][0] == TOKEN_KEYWORD:
            parser_error(line_index, 'Unknown construction.')
        
        if line_tokens[0][1] == KEYWORD_END:
            multilvl_expr -= 1
        
        if line_tokens[0][1] == KEYWORD_ELSE:
            multilvl_expr -= 1
        
        curr_expr.append((line_tokens, line_index + 1, multilvl_expr))
        
        if line_tokens[0][1] == KEYWORD_IF:
            multilvl_expr += 1
        
        if line_tokens[0][1] == KEYWORD_ELSE:
            multilvl_expr += 1
        
        if multilvl_expr == 0:
            expr_list.append(curr_expr)
            curr_expr = []

    for expr in expr_list:
        expression_ast = parse_expr(expr)
        AST.append(expression_ast)
    return AST