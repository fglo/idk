import re
from idk_consts import * 

class LexerError(Exception):
    pass

def lexer_error(line_index, message):
    raise LexerError('Lexer error in line %d: %s' % (line_index, message))

def try_parse_literal(token):
    if token.lstrip("-").isnumeric():
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
    if token == 'for':
        return KEYWORD_FOR
    if token == 'while':
        return KEYWORD_WHILE
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
        return OPERATOR_ASTERISK
    if token == '/':
        return OPERATOR_SLASH
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
    if token == '..':
        return OPERATOR_RANGE
    if token == '..=':
        return OPERATOR_RANGE_INCLUSIVE
    if token == 'in':
        return OPERATOR_IN
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

#TODO: include other operators
def add_spaces_to_operators(line):
    line = line.replace('..=', ' ..= ')
    line = line.replace('..', ' .. ')
    line = line.replace('.. =', ' ..= ')
    # operators = [':=', '+', '-', '*', '/', '=', '>', '>=', '<', '<=', '..']
    # for op in operators:
    #     if re.match(f'[a-z0-9\s]{op}[a-z0-9\s]', line):
    #         line = line.replace(op, f' {op} ')
    return line
    
def tokenize_line(line, line_index) -> list:
    line = line.strip()
    if not line or (len(line) >= 2 and line[0:2] == '//'):
        return []
    
    tokens = add_spaces_to_operators(line).split(' ')
    
    cleared_tokens = []
    for token in tokens:
        if len(token) >= 2 and token[0:2] == '//':
            break
        if token:
            cleared_tokens.append(token)
    
    interpreted_tokens = interpret_line_tokens(cleared_tokens, line_index)
    return interpreted_tokens

def tokenize(program_file_path) -> list:
    tokenized_file_lines = []
    line_index = 0
    with open(program_file_path, "r") as pf:
        for line in pf:
            line_index += 1
            tokenized_line = tokenize_line(line, line_index)
            tokenized_file_lines.append(tokenized_line) 
    return tokenized_file_lines