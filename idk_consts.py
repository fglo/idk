# ENUM GENERATOR:

current_enum_value = 0
def nextval():
    global current_enum_value
    value = current_enum_value
    current_enum_value += 1
    return value        

def NEW_ENUM(name):
    global current_enum_value
    current_enum_value = 0

# ENUMS:

NEW_ENUM('TOKENS')
TOKEN_INT = nextval()
TOKEN_CHAR = nextval()
TOKEN_BOOL = nextval()
TOKEN_WORD = nextval()
TOKEN_OPERATOR = nextval()
TOKEN_KEYWORD = nextval()
COUNT_TOKENS = nextval()

NEW_ENUM('OPERATORS')
OPERATOR_PLUS = nextval()
OPERATOR_MINUS = nextval()
OPERATOR_MULTIPLICATION = nextval()
OPERATOR_DIVISION = nextval()
OPERATOR_ASSIGMENT = nextval()
OPERATOR_EQ = nextval()
OPERATOR_GT = nextval()
OPERATOR_LT = nextval()
COUNT_OPERATORS = nextval()

NEW_ENUM('KEYWORDS')
KEYWORD_PRINT = nextval()
KEYWORD_IF = nextval()
KEYWORD_END = nextval()
COUNT_KEYWORDS = nextval()