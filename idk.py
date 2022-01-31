from sys import argv
from idk_lexer import *
from idk_parser import *
from idk_syntax_checker import *
from idk_evaluator import *

def run(program_file_path):
    try:
        tokenized_file_lines = tokenize(program_file_path)
        ast = parse(tokenized_file_lines)
        check_syntax(ast)
        evaluate(ast)
    except LexerError as e:
        print('[ERROR]', e)
    except ParserError as e:
        print('[ERROR]', e)
    except SyntaxError as e:
        print('[ERROR]', e)
    except EvaluatorError as e:
        print('[ERROR]', e)

def run_interactive():
    print('Welcome to IDK interactive!')
    value = input("$ ")
    line_index = 1
    assignments = []
    tokenized_expr = []
    multilvl_expr = 0
    while value != 'exit':
        try:
            tokenized_line = tokenize_line(value, line_index)
            if tokenized_line[0][0] == TOKEN_KEYWORD and (tokenized_line[0][1] == KEYWORD_IF or tokenized_line[0][1] == KEYWORD_FOR):
                multilvl_expr += 1
            if tokenized_line[0][0] == TOKEN_KEYWORD and (tokenized_line[0][1] == KEYWORD_END):
                multilvl_expr -= 1
            
            tokenized_expr.append(tokenized_line)
            
            if multilvl_expr == 0:
                if len(tokenized_line) > 1 and tokenized_line[1][0] == TOKEN_OPERATOR and tokenized_line[0][1] == OPERATOR_ASSIGMENT:
                    assignments.append(tokenize_line)
                reset_ast()
                ast = parse(assignments + tokenized_expr)
                evaluate(ast)
                tokenized_expr = []
        except LexerError as e:
            print('[ERROR]', e)
        except ParserError as e:
            print('[ERROR]', e)
        except SyntaxError as e:
            print('[ERROR]', e)
        except EvaluatorError as e:
            print('[ERROR]', e)
        if multilvl_expr == 0:
            value = input("$ ")
        else:
            value = input("> ")
        line_index += 1
    
if __name__ == "__main__":
    if argv[1] == '-it':
        run_interactive()
    else:
        run(argv[1])