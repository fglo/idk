from typing import Optional, Union
from idk_consts import *
from idk_types import Expression, KeywordExpression, OperatorExpression, Token, TokenExpression 

AST:list[Expression] = []

class ParserError(Exception):
    pass

def parser_error(line_index, message):
    raise ParserError('Parser error in line %d: %s' % (line_index, message))

def get_first_operator_index(tokens:list[Token]):
    operator_index = -1
    for index, token in enumerate(tokens):
        if token.type == TOKEN_OPERATOR:
            operator_index = index
            break
    return operator_index

def get_first_less_important_operator_index(tokens:list[Token]):
    found_operator_index:int = -1
    found_operator:Optional[Token] = None
    for index, token in enumerate(tokens):
        if token.type != TOKEN_OPERATOR:
            continue
        if token.value == OPERATOR_ASSIGMENT:
            found_operator_index = index
            break
        if found_operator is None:
            found_operator_index = index
            found_operator = token
            continue
        if ((found_operator.value == OPERATOR_ASTERISK or found_operator.value == OPERATOR_SLASH) and token.value >= OPERATOR_PLUS) \
            or ((found_operator.value == OPERATOR_PLUS or found_operator.value == OPERATOR_MINUS) and token.value >= OPERATOR_EQ) \
            or ((found_operator.value <= OPERATOR_LTE) and token.value >= OPERATOR_AND) \
            or ((found_operator.value == OPERATOR_AND) and (token.value == OPERATOR_OR or token.value == OPERATOR_XOR)):
            found_operator_index = index
            found_operator = token
            continue
    return found_operator_index
        
def handle_side_of_operator(expr:TokenExpression)-> Expression:
    local_ast:Expression
    if len(expr.tokens) == 1:
        local_ast = Expression(expr.tokens[0].type, expr.tokens[0].value, expr.line_index)
    else:
        local_ast = handle_operator(expr)
    return local_ast
        
def handle_operator_assignment(expr:TokenExpression, operator_index)-> OperatorExpression:
    left = Expression(expr.tokens[0].type, expr.tokens[0].value, expr.line_index)
    for expr_ast in AST:
        if expr_ast.value == OPERATOR_ASSIGMENT and expr_ast.args[1].value == left.value:
            parser_error(expr.line_index, "You cannot assign value to a variable that already has been assigned.")
    right = handle_side_of_operator(TokenExpression(expr.tokens[operator_index + 1:], expr.line_index, expr.nesting_lvl))
    expr_ast = OperatorExpression(TOKEN_OPERATOR, OPERATOR_ASSIGMENT, expr.line_index, args=[left, right])
    return expr_ast
    
def handle_operator_with_one_side_evaluation(expr:TokenExpression, operator_index)-> OperatorExpression:
    right = handle_side_of_operator(TokenExpression(expr.tokens[operator_index + 1:], expr.line_index, expr.nesting_lvl))
    token_operator = expr.tokens[operator_index].value
    expr_ast = OperatorExpression(TOKEN_OPERATOR, token_operator, expr.line_index, args=[right])
    return expr_ast
    
def handle_operator_with_two_side_evaluation(expr:TokenExpression, operator_index)-> OperatorExpression:
    left = handle_side_of_operator(TokenExpression(expr.tokens[0:operator_index], expr.line_index, expr.nesting_lvl))
    right = handle_side_of_operator(TokenExpression(expr.tokens[operator_index + 1:], expr.line_index, expr.nesting_lvl))
    token_operator = expr.tokens[operator_index].value
    expr_ast = OperatorExpression(TOKEN_OPERATOR, token_operator, expr.line_index, args=[left, right])
    return expr_ast
    
def handle_operator_in_with_range(expr:TokenExpression, operator_index)-> OperatorExpression:
    left = handle_side_of_operator(TokenExpression(expr.tokens[0:operator_index], expr.line_index, expr.nesting_lvl))
    right = handle_side_of_operator(TokenExpression(expr.tokens[operator_index + 1:], expr.line_index, expr.nesting_lvl))
    operator_type = expr.tokens[operator_index].type
    expr_ast = OperatorExpression(TOKEN_OPERATOR, operator_type, expr.line_index, args=[left, right])
    return expr_ast

def handle_operator(expr:TokenExpression)-> OperatorExpression:
    operator_index = get_first_less_important_operator_index(expr.tokens)
    token_operator = expr.tokens[operator_index].value
    if token_operator == OPERATOR_ASSIGMENT:
        expr_ast = handle_operator_assignment(expr, operator_index)
    elif token_operator == OPERATOR_PLUS:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_MINUS:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_ASTERISK:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_SLASH:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_EQ:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_GT:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_GTE:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_LT:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_LTE:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_XOR:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_OR:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_AND:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_NOT:
        expr_ast = handle_operator_with_one_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_IN:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_RANGE:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    elif token_operator == OPERATOR_RANGE_INCLUSIVE:
        expr_ast = handle_operator_with_two_side_evaluation(expr, operator_index)
    else:
        parser_error(expr.line_index, "Unknown operator")
    return expr_ast

def handle_keyword_print(expr:TokenExpression) -> KeywordExpression:
    operator_index = get_first_less_important_operator_index(expr.tokens)
    expr_ast:Expression
    if operator_index > -1:
        expr_ast = handle_operator(TokenExpression(expr.tokens[1:], expr.line_index, expr.nesting_lvl))
    elif expr.tokens[1].type <= TOKEN_WORD:
        expr_ast = Expression(expr.tokens[1].type, expr.tokens[1].value, expr.line_index)
    else:
        parser_error(expr.line_index, 'You can print only literals, variables or operations.')

    expr_ast = KeywordExpression(TOKEN_KEYWORD, KEYWORD_PRINT, expr.line_index, keyword_expr=expr_ast)
    return expr_ast

def handle_keyword_else(expr:TokenExpression) -> KeywordExpression:
    if len(expr.tokens) > 1 and expr.tokens[1].type == TOKEN_KEYWORD and expr.tokens[1].value == KEYWORD_IF:
        del expr.tokens[0]
        return handle_keyword_if(expr)
    
    exprs:list[TokenExpression] = []
    for internal_expr in expr.nested_exprs:
        token = internal_expr.tokens[0]
        if token.type == TOKEN_KEYWORD and token.value == KEYWORD_END:
            break
        exprs.append(internal_expr)
    
    expr_ast_prim:list[Expression] = []
    for e in exprs:
        expr_ast_prim.append(parse_expr(e))
    
    expr_ast = KeywordExpression(TOKEN_KEYWORD, KEYWORD_ELSE, expr.line_index, keyword_expr=Expression(TOKEN_BOOL, 0, expr.line_index), prim=expr_ast_prim)
    return expr_ast
    
def handle_keyword_if(expr:TokenExpression) -> KeywordExpression:
    if expr.tokens[1].type == TOKEN_KEYWORD and expr.tokens[1].value == KEYWORD_IF:
        expr.tokens = expr.tokens[1:]
            
    operator_index = get_first_less_important_operator_index(expr.tokens)
    expr_value:Expression
    if operator_index > -1:
        expr_value = handle_operator(TokenExpression(expr.tokens[1:], expr.line_index, expr.nesting_lvl))
    elif expr.tokens[1].type <= TOKEN_WORD:
        expr_value = Expression(expr.tokens[1].type, expr.tokens[1].value, expr.line_index)
    else:
        parser_error(expr.line_index, 'You can do operations only on literals, variables or other operations.')
        
    exprs_prim:list[TokenExpression] = []
    exprs_alt:TokenExpression = TokenExpression([Token(TOKEN_KEYWORD, KEYWORD_ELSE)], expr.line_index, expr.nesting_lvl)
    else_found = False
    for internal_expr in expr.nested_exprs:
        token = internal_expr.tokens[0]
        if token.type == TOKEN_KEYWORD and token.value == KEYWORD_END:
            break
        if token.type == TOKEN_KEYWORD and token.value == KEYWORD_ELSE:
            else_found = True
        
        if else_found:
            exprs_alt.nested_exprs.append(internal_expr)
        else:
            exprs_prim.append(internal_expr)
    
    expr_ast: Expression
    
    expr_ast_prim:list[Expression] = []
    for e in exprs_prim:
        expr_ast_prim.append(parse_expr(e))
    
    expr_ast_alt:Expression
    if else_found:
        expr_ast_alt = handle_keyword_else(exprs_alt)    
        expr_ast = KeywordExpression(TOKEN_KEYWORD, KEYWORD_IF, expr.line_index, keyword_expr=expr_value, prim=expr_ast_prim, alt=[expr_ast_alt])
    else:
        expr_ast = KeywordExpression(TOKEN_KEYWORD, KEYWORD_IF, expr.line_index, keyword_expr=expr_value, prim=expr_ast_prim)
    return expr_ast

def handle_keyword_for(expr:TokenExpression) -> KeywordExpression:
    operator_index = get_first_operator_index(expr.tokens)
    if operator_index > -1:
        expr_value = handle_operator(TokenExpression(expr.tokens[1:], expr.line_index, expr.nesting_lvl))
        
    if operator_index == -1 or (expr_value.type != OPERATOR_RANGE and expr_value.value != OPERATOR_RANGE_INCLUSIVE):
        parser_error(expr.line_index, "Only for range loops are currently supported.")
    
    exprs:list[TokenExpression]
    for internal_expr in expr.nested_exprs:
        token = internal_expr.tokens[0]
        if token.type == TOKEN_KEYWORD and token.value == KEYWORD_END:
            break
        exprs.append(internal_expr)
    
    expr_ast_prim = []
    for expr in exprs:
        expr_ast_prim.append(parse_expr(expr))
            
    expr_ast = KeywordExpression(TOKEN_KEYWORD, KEYWORD_FOR, expr.line_index, keyword_expr=expr_value, prim=expr_ast_prim)
    return expr_ast
    
def handle_keyword_action(expr:TokenExpression) -> Expression:
    expr_ast:Expression
    if expr.tokens[0].value == KEYWORD_PRINT:
        expr_ast = handle_keyword_print(expr)
    elif expr.tokens[0].value == KEYWORD_IF:
        expr_ast = handle_keyword_if(expr)
    elif expr.tokens[0].value == KEYWORD_FOR:
        expr_ast = handle_keyword_for(expr)
    elif expr.tokens[0].value == KEYWORD_ELSE:
        parser_error(expr.line_index, 'KEYWORD_ELSE was not expected.') 
    elif expr.tokens[0].value == KEYWORD_END:
        parser_error(expr.line_index, 'KEYWORD_END was not expected.') 
    else:
        parser_error(expr.line_index, 'Unknown keyword.') 
    return expr_ast

def parse_expr(expr:TokenExpression) -> Expression:
    expr_ast:Expression
    if expr.tokens[0].type == TOKEN_KEYWORD:
        expr_ast = handle_keyword_action(expr)    
    elif len(expr.tokens) > 1 and expr.tokens[0].type == TOKEN_WORD and expr.tokens[1].type == TOKEN_OPERATOR:
        expr_ast = handle_operator(expr)
    else:
        parser_error(expr.line_index, 'Not allowed word on the beginning of the line.')
    return expr_ast


def get_expr_list(tokenized_code_lines:list[list[Token]]):
    expr_list:list[TokenExpression] = []
    curr_expr:Optional[TokenExpression] = None
    multilvl_expr = 0
    for line_index, line_tokens in enumerate(tokenized_code_lines):
        if len(line_tokens) == 0:
            continue
    
        if len(line_tokens) >= 2 and line_tokens[0].type == TOKEN_WORD and line_tokens[1].type == TOKEN_WORD:
            parser_error(line_index, 'Unknown construction.')
        
        if len(line_tokens) >= 2 and line_tokens[0].type == TOKEN_WORD and line_tokens[1].type == TOKEN_KEYWORD:
            parser_error(line_index, 'Unknown construction.')
            
        first_token = line_tokens[0]
        
        if first_token.type == TOKEN_KEYWORD and first_token.value == KEYWORD_END:
            multilvl_expr -= 1
        
        if curr_expr is None:
            curr_expr = TokenExpression(line_tokens, line_index + 1, multilvl_expr)
        else:
            curr_expr.nested_exprs.append(TokenExpression(line_tokens, line_index + 1, multilvl_expr))
        
        if first_token.type == TOKEN_KEYWORD:
            if first_token.value == KEYWORD_IF \
                or first_token.value == KEYWORD_FOR \
                or first_token.value == KEYWORD_WHILE:
                multilvl_expr += 1
        
        if multilvl_expr == 0:
            expr_list.append(curr_expr)
            curr_expr = None

    return expr_list

def parse(tokenized_code_lines:list[list[Token]]):
    expr_list = get_expr_list(tokenized_code_lines)

    for expr in expr_list:
        expression_ast = parse_expr(expr)
        AST.append(expression_ast)
    return AST

def reset_ast():
    global AST
    AST = []