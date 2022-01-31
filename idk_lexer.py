from idk_consts import * 

class LexerError(Exception):
    pass

def lexer_error(line_index, message):
    raise LexerError('Error in line %d: %s' % (line_index, message))

def try_parse_literal(token):
    if token.isnumeric():
        return TOKEN_INT
    if len(token) == 3 and token[0] == "'" and token[2] =="'":
        return TOKEN_CHAR
    if token == "true" or token == "false":
        return TOKEN_BOOL
    return -1
    
def try_parse_keyword(token):
    if token == 'print':
        return KEYWORD_PRINT
    if token == 'if':
        return KEYWORD_IF
    if token == 'else':
        return KEYWORD_ELSE
    if token == 'end':
        return KEYWORD_END
    return -1
    
def try_parse_operator(token):
    if token == ':=':
        return OPERATOR_ASSIGMENT
    if token == '+':
        return OPERATOR_PLUS
    if token == '-':
        return OPERATOR_MINUS
    if token == '*':
        return OPERATOR_MULTIPLICATION
    if token == '/':
        return OPERATOR_DIVISION
    if token == '=':
        return OPERATOR_EQ
    if token == '>':
        return OPERATOR_GT
    if token == '>=':
        return OPERATOR_GTE
    if token == '<':
        return OPERATOR_LT
    if token == '<=':
        return OPERATOR_LTE
    if token == 'not':
        return OPERATOR_NOT
    if token == 'and':
        return OPERATOR_AND
    if token == 'or':
        return OPERATOR_OR
    if token == 'xor':
        return OPERATOR_XOR
    return -1

def interpret_line_tokens(tokens, line_index):
    interpreted_line = []
    for token in tokens:
        if not token:
            lexer_error(line_index, "Empty token!")
        literal = try_parse_literal(token)
        if literal > -1:
            if literal == TOKEN_INT:
                interpreted_line.append((TOKEN_INT, int(token)))
            elif literal == TOKEN_CHAR:
                interpreted_line.append((TOKEN_CHAR, ord(token[1])))
            elif literal == TOKEN_BOOL:
                interpreted_line.append((TOKEN_BOOL, int(token == "true")))
            continue
        keyword = try_parse_keyword(token)
        if keyword > -1:
            interpreted_line.append((TOKEN_KEYWORD, keyword))
            continue
        operator = try_parse_operator(token)
        if operator > -1:
            interpreted_line.append((TOKEN_OPERATOR, operator))
            continue
        interpreted_line.append((TOKEN_WORD, token))
    return interpreted_line

 #TODO: tokenize also using operators
def tokenize_line(line, line_index) -> list:
    line = line.strip()
    if not line or line[0] == '/' and line[1] == '/':
        return []
    
    tokens = line.split(' ')
    
    cleared_tokens = []
    for token in tokens:
        if len(token) >= 2 and token[0] == '/' and token[1] == '/':
            break
        if token:
            cleared_tokens.append(token)
    
    # tokens = [token for token in tokens if token]
    return interpret_line_tokens(cleared_tokens, line_index)

def tokenize(program_file_path) -> list:
    tokenized_file_lines = []
    line_index = 0
    with open(program_file_path, "r") as pf:
        for line in pf:
            line_index += 1
            tokenized_line = tokenize_line(line, line_index)
            tokenized_file_lines.append(tokenized_line) 
    return tokenized_file_lines