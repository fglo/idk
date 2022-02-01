from dataclasses import dataclass, field
from typing import Optional, Union

@dataclass
class Token:
    type: int
    value: int

@dataclass
class TokenExpression:
    tokens: list[Token]
    line_index: int
    nesting_lvl: int
    nested_exprs: list['TokenExpression'] = field(default_factory=list)

@dataclass
class Expression:
    type: int
    value: Union[int, list]
    line_index: int
    
@dataclass
class OperatorExpression(Expression):
    args: list['Expression']

@dataclass
class KeywordExpression(Expression):
    keyword_expr: 'Expression'
    prim: list['Expression'] = field(default_factory=list)
    alt: list['Expression'] = field(default_factory=list)
    