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
            if found_operator is None:
                found_operator_index = index
                found_operator = token
            elif found_operator[1] >= OPERATOR_MULTIPLICATION and token[1] <= OPERATOR_MINUS:
                found_operator_index = index
                found_operator = token
                break
    return found_operator_index
        
def handle_side_of_operator(side, line_index):
    if len(side) == 1:
        side = side[0]
    else:
        side = handle_operator(side, line_index)
    return side
        
def handle_operator_assignment(tokens, line_index, operator_index):
    left = tokens[0]
    for expr in AST:
        if expr[1] == OPERATOR_ASSIGMENT and expr[2][1] == left[1]:
            parser_error(line_index, "You cannot assign value to a variable that already has been assigned.")
    right = handle_side_of_operator(tokens[operator_index + 1:], line_index)
    return (TOKEN_OPERATOR, OPERATOR_ASSIGMENT, left, right)  
    
def handle_operator_plus(tokens, line_index, operator_index):        
    left = handle_side_of_operator(tokens[0:operator_index], line_index)
    right = handle_side_of_operator(tokens[operator_index + 1:], line_index)
    return (TOKEN_OPERATOR, OPERATOR_PLUS, left, right)  
    
def handle_operator_minus(tokens, line_index, operator_index):        
    left = handle_side_of_operator(tokens[0:operator_index], line_index)
    right = handle_side_of_operator(tokens[operator_index + 1:], line_index)
    return (TOKEN_OPERATOR, OPERATOR_MINUS, left, right)  
    
def handle_operator_multiplication(tokens, line_index, operator_index):        
    left = handle_side_of_operator(tokens[0:operator_index], line_index)
    right = handle_side_of_operator(tokens[operator_index + 1:], line_index)
    return (TOKEN_OPERATOR, OPERATOR_MULTIPLICATION, left, right)  
    
def handle_operator_division(tokens, line_index, operator_index):        
    left = handle_side_of_operator(tokens[0:operator_index], line_index)
    right = handle_side_of_operator(tokens[operator_index + 1:], line_index)
    return (TOKEN_OPERATOR, OPERATOR_DIVISION, left, right)

def handle_operator(tokens, line_index):
    operator_index = get_first_less_important_operator_index(tokens)
    token_operator = tokens[operator_index][1]
    if token_operator == OPERATOR_ASSIGMENT:
        expr = handle_operator_assignment(tokens, line_index, operator_index)
    elif token_operator == OPERATOR_PLUS:
        expr = handle_operator_plus(tokens, line_index, operator_index)
    elif token_operator == OPERATOR_MINUS:
        expr = handle_operator_minus(tokens, line_index, operator_index)
    elif token_operator == OPERATOR_MULTIPLICATION:
        expr = handle_operator_multiplication(tokens, line_index, operator_index)
    elif token_operator == OPERATOR_DIVISION:
        expr = handle_operator_division(tokens, line_index, operator_index)
    else:
        parser_error(line_index, "Unknown operator")
    return expr
    
def handle_keyword_print(tokens, line_index):
    operator_index = get_first_operator_index(tokens)
    if operator_index > -1:
        expr_value = handle_operator(tokens[1:], line_index)
    elif tokens[1][0] <= TOKEN_WORD:
        expr_value = tokens[1]
    else:
        parser_error(line_index, 'You can print only literals, variables or operations.')
    return (TOKEN_KEYWORD, KEYWORD_PRINT, expr_value)
    
def handle_keyword_action(tokens, line_index):  
    if tokens[0][1] == KEYWORD_PRINT:
        expr = handle_keyword_print(tokens, line_index)
    else:
        parser_error(line_index, 'Unknown keyword.') 
    return expr

def parse_line(tokens, line_index):
    if tokens[0][0] == TOKEN_WORD and tokens[1][0] == TOKEN_WORD:
        parser_error(line_index, 'Unknown construction.')
    
    if tokens[0][0] == TOKEN_WORD and tokens[1][0] == TOKEN_KEYWORD:
        parser_error(line_index, 'Unknown construction.')

    if tokens[0][0] == TOKEN_KEYWORD:
        expr = handle_keyword_action(tokens, line_index)    
    elif tokens[0][0] == TOKEN_WORD and tokens[1][0] == TOKEN_OPERATOR:
        expr = handle_operator(tokens, line_index)
    else:
        parser_error(line_index, 'Not allowed word on the beginning of the line.')
    return expr

def parse(tokenized_lines):
    line_index = 0
    for line in tokenized_lines:
        line_index += 1
        AST.append(parse_line(line, line_index))
    return AST