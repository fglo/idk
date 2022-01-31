from idk_consts import * 

class SyntaxError(Exception):
    pass

def syntax_error(line_index, message):
    raise SyntaxError('Syntax error in line %d: %s' % (line_index, message))

def get_operator_definition(op):
    return SYNTAX_RULES[TOKEN_OPERATOR][op]

def get_keyword_definition(kw):
    return SYNTAX_RULES[TOKEN_KEYWORD][kw]

def get_returned_types(arg_exprs:list) -> list:
    types = []
    for expr in arg_exprs:
        if len(expr) == 0:
            continue
        if expr[0] == TOKEN_INT:
            types.append(TOKEN_INT)
        elif expr[0] == TOKEN_CHAR:
            types.append(TOKEN_CHAR)
        elif expr[0] == TOKEN_BOOL:
            types.append(TOKEN_BOOL)
        elif expr[0] == TOKEN_ARRAY:
            types.append(TOKEN_ARRAY)
        elif expr[0] == TOKEN_WORD:
            types.append(TOKEN_WORD)
        elif expr[0] == TOKEN_OPERATOR:
            rules = get_operator_definition(expr[1])
            types.append(rules['returns'])
    return types

def check_token(expr):
    if expr[0] != TOKEN_OPERATOR and expr[0] != TOKEN_KEYWORD:
        return
    
    token_type, token_value, *arg_exprs, line_index = expr
    if token_type == TOKEN_OPERATOR:
        rules = get_operator_definition(token_value)
        args_return_types = get_returned_types(arg_exprs)
        if (token_value != OPERATOR_NOT and not (args_return_types[0], args_return_types[1]) in rules['arguments']) \
            or token_value == OPERATOR_NOT and not args_return_types[0] in rules['arguments']:
            syntax_error(line_index, f"Wrong argument types for operator {rules['operator']}")
    elif token_type == TOKEN_KEYWORD:
        rules = get_keyword_definition(token_value)
        args_return_types = get_returned_types(arg_exprs)
        if not args_return_types[0] in rules['arguments']:
            syntax_error(line_index, f"Wrong operation return type for keyword {rules['keyword']}")
    
    for expr in arg_exprs:
        if len(expr) == 0:
            continue
        check_token(expr)
    
def check_syntax(ast):
    for expr in ast:
        check_token(expr)
        
        
        
SYNTAX_RULES = {
    TOKEN_OPERATOR: {
        OPERATOR_ASSIGMENT: {
            'operator': 'OPERATOR_ASSIGMENT',
            'arguments' : [
                (TOKEN_WORD, TOKEN_WORD),
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)      
            ]
        },
        OPERATOR_PLUS: {
            'operator': 'OPERATOR_PLUS', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_INT
        },
        OPERATOR_MINUS: {
            'operator': 'OPERATOR_MINUS', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_INT
        },
        OPERATOR_ASTERISK: {
            'operator': 'OPERATOR_ASTERISK', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_INT
        },
        OPERATOR_SLASH: {
            'operator': 'OPERATOR_SLASH', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_INT
        },
        OPERATOR_EQ: {
            'operator': 'OPERATOR_EQ', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_BOOL
        },
        OPERATOR_GT: {
            'operator': 'OPERATOR_GT', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_BOOL
        },
        OPERATOR_GTE: {
            'operator': 'OPERATOR_GTE', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_BOOL
        },
        OPERATOR_LT: {
            'operator': 'OPERATOR_LT', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_BOOL
        },
        OPERATOR_LTE: {
            'operator': 'OPERATOR_LTE', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_BOOL
        },
        OPERATOR_IN: {
            'operator': 'OPERATOR_IN', 
            'arguments' : [
                (TOKEN_INT, TOKEN_ARRAY),        
                (TOKEN_CHAR, TOKEN_ARRAY),        
                (TOKEN_BOOL, TOKEN_ARRAY),        
                (TOKEN_WORD, TOKEN_ARRAY),           
            ],
            'returns': TOKEN_BOOL
        },
        OPERATOR_RANGE: {
            'operator': 'OPERATOR_RANGE', 
            'arguments' : [
                (TOKEN_INT, TOKEN_INT),        
                (TOKEN_INT, TOKEN_BOOL),        
                (TOKEN_INT, TOKEN_CHAR),        
                (TOKEN_INT, TOKEN_WORD),        
                (TOKEN_CHAR, TOKEN_CHAR),        
                (TOKEN_CHAR, TOKEN_INT),        
                (TOKEN_CHAR, TOKEN_BOOL),        
                (TOKEN_CHAR, TOKEN_WORD),
                (TOKEN_BOOL, TOKEN_BOOL),        
                (TOKEN_BOOL, TOKEN_INT),        
                (TOKEN_BOOL, TOKEN_CHAR),        
                (TOKEN_BOOL, TOKEN_WORD), 
                (TOKEN_WORD, TOKEN_WORD),           
                (TOKEN_WORD, TOKEN_INT),        
                (TOKEN_WORD, TOKEN_BOOL),        
                (TOKEN_WORD, TOKEN_CHAR)       
            ],
            'returns': TOKEN_ARRAY
        },
        OPERATOR_AND: {
            'operator': 'OPERATOR_AND',
            'arguments' : [
                (TOKEN_BOOL, TOKEN_BOOL),
            ],
            'returns' : TOKEN_BOOL
        },
        OPERATOR_OR: {
            'operator': 'OPERATOR_OR',
            'arguments' : [
                (TOKEN_BOOL, TOKEN_BOOL),
            ],
            'returns' : TOKEN_BOOL
        },
        OPERATOR_XOR: {
            'operator': 'OPERATOR_XOR',
            'arguments' : [
                (TOKEN_BOOL, TOKEN_BOOL),
            ],
            'returns' : TOKEN_BOOL
        },
        OPERATOR_NOT: {
            'operator': 'OPERATOR_NOT', 
            'arguments' : [
                TOKEN_BOOL      
            ],
            'returns' : TOKEN_BOOL
        },
    },
    TOKEN_KEYWORD: {
        KEYWORD_IF: {
            'keyword': 'KEYWORD_IF',
            'arguments': [TOKEN_BOOL, TOKEN_WORD]
        },
        KEYWORD_FOR: {
            'keyword': 'KEYWORD_FOR',
            'arguments': [TOKEN_BOOL, TOKEN_WORD, TOKEN_ARRAY]
        },
        KEYWORD_PRINT: {
            'keyword': 'KEYWORD_PRINT',
            'arguments': [TOKEN_INT, TOKEN_CHAR, TOKEN_BOOL, TOKEN_WORD, TOKEN_ARRAY]
        }
    }
}
