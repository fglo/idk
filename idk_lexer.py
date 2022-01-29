from idk_consts import * 

def try_parse_literal(token):
    if token.isnumeric():
        return TOKEN_INT
    if len(token) == 3 and token[0] == "'" and token[2] =="'":
        return TOKEN_CHAR
    return -1
    
def try_parse_keyword(token):
    if token == 'print':
        return KEYWORD_PRINT
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
    return -1

def interpret_line_tokens(tokens):
    interpreted_line = []
    for token in tokens:
        if not token:
            assert False, "Empty token!"
        literal = try_parse_literal(token)
        if literal > -1:
            if literal == TOKEN_INT:
                interpreted_line.append((TOKEN_INT, int(token)))
            elif literal == TOKEN_CHAR:
                interpreted_line.append((TOKEN_CHAR, ord(token[1])))
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

def tokenize_line(line) -> list:
    line = line.strip()
    tokens = line.split(' ')
    tokens = [token for token in tokens if token]
    return interpret_line_tokens(tokens)

def tokenize(program_file_path) -> list:
    tokenized_file_lines = []
    with open(program_file_path, "r") as pf:
        for line in pf:
            tokenized_line = tokenize_line(line)
            if len(tokenized_line) > 0:
                tokenized_file_lines.append(tokenized_line) 
    return tokenized_file_lines