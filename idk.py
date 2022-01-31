from sys import argv
from idk_lexer import *
from idk_parser import *
from idk_evaluator import *

def run(program_file_path):
    tokenized_file_lines = tokenize(program_file_path)
    ast = parse(tokenized_file_lines)
    evaluate(ast)

def run_interactive():
    print('Welcome to IDK interactive!')
    value = input("$ ")
    line_index = 1
    while value != 'exit':
        try:
            tokenized_line = [tokenize_line(value, line_index)]
            ast = parse_interactive(tokenized_line)
            evaluate(ast)
        except Exception as e:
            print(e)
        value = input("$ ")
        line_index += 1
    
if __name__ == "__main__":
    if argv[1] == '-it':
        run_interactive()
    else:
        run(argv[1])