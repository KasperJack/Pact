from lark import Lark

grammar = r"""
start: pair+

pair: NAME "=" value

?value: STRING
      | DATE
      | NUMBER

NAME: /[a-zA-Z_][a-zA-Z0-9_]*/

STRING: ESCAPED_STRING
DATE: /\d{4}-\d{2}-\d{2}/
NUMBER: /\d+(\.\d+)?/

%import common.ESCAPED_STRING
%import common.WS
%ignore WS
"""

parser = Lark(grammar, parser="lalr", propagate_positions=True)





text = """
id = "official_16"
source = "official"
released = 2024-01-15
version = "1.6"
"""

tree = parser.parse(text)
print(tree.pretty())