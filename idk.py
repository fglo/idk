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
    while value != 'exit':
        try:
            tokenized_line = [tokenize_line(value, line_index)]
            ast = parse_interactive(tokenized_line)
            evaluate(ast)
        except LexerError as e:
            print('[ERROR]', e)
        except ParserError as e:
            print('[ERROR]', e)
        except SyntaxError as e:
            print('[ERROR]', e)
        except EvaluatorError as e:
            print('[ERROR]', e)
        value = input("$ ")
        line_index += 1
    
if __name__ == "__main__":
    if argv[1] == '-it':
        run_interactive()
    else:
        run(argv[1])