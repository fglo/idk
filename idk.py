from sys import argv
from idk_lexer import *
from idk_parser import *
from idk_evaluator import *

def run(program_file_path):
    tokenized_file_lines = tokenize(program_file_path)
    ast = parse(tokenized_file_lines)
    evaluate(ast)

if __name__ == "__main__":
    run(argv[1])