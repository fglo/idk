from idk_consts import *
from idk_types import Expression, KeywordExpression, OperatorExpression 

class SyntaxError(Exception):
    pass

def syntax_error(line_index, message):
    raise SyntaxError('Syntax error in line %d: %s' % (line_index, message))

def get_operator_definition(op:int):
    return SYNTAX_RULES[TOKEN_OPERATOR][op]

def get_keyword_definition(kw:int):
    return SYNTAX_RULES[TOKEN_KEYWORD][kw]

def get_returned_types(exprs:list[Expression]) -> list:
    types:list[int] = []
    for expr in exprs:
        if expr is None:
            continue
        if expr.type == TOKEN_INT:
            types.append(TOKEN_INT)
        elif expr.type == TOKEN_CHAR:
            types.append(TOKEN_CHAR)
        elif expr.type == TOKEN_BOOL:
            types.append(TOKEN_BOOL)
        elif expr.type == TOKEN_ARRAY:
            types.append(TOKEN_ARRAY)
        elif expr.type == TOKEN_WORD:
            types.append(TOKEN_WORD)
        elif expr.type == TOKEN_OPERATOR:
            rules = get_operator_definition(expr.value)
            types.append(rules['returns'])
    return types

def check_expression(expr:Expression):
    if expr.type != TOKEN_OPERATOR and expr.type != TOKEN_KEYWORD:
        return
    
    exprs:list[Expression] = []
    
    if expr.type == TOKEN_OPERATOR:
        rules = get_operator_definition(expr.value)
        args_return_types = get_returned_types(expr.args)
        if (expr.value != OPERATOR_NOT and not (args_return_types[0], args_return_types[1]) in rules['arguments']) \
            or expr.value == OPERATOR_NOT and not args_return_types[0] in rules['arguments']:
            syntax_error(expr.line_index, f"Wrong argument types for operator {rules['operator']}")
        exprs += expr.args
    elif expr.type == TOKEN_KEYWORD:
        rules = get_keyword_definition(expr.value)
        args_return_types = get_returned_types([expr.keyword_expr])
        if not args_return_types[0] in rules['arguments']:
            syntax_error(expr.line_index, f"Wrong operation return type for keyword {rules['keyword']}")
        exprs.append(expr.keyword_expr)
        exprs += expr.prim
        exprs += expr.alt
    
    for expr in exprs:
        if expr is None:
            continue
        check_expression(expr)
    
def check_syntax(ast):
    for expr in ast:
        check_expression(expr)
        
        
        
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
        OPERATOR_RANGE_INCLUSIVE: {
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
