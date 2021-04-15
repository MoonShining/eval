grammar cond_parser;

expr
   :  binary #BinaryExpr
   |  expr op=(AND|OR) expr #LogicExpr
   |  LPAREN expr RPAREN #ParenExpr
   ;


binary
   : ID op=(EQ|NEQ) atom #CompareBinary
   ;

atom
   : NULL
   | STRING
   | NUMBER
   ;

WS: [ \r\n\t] -> skip;
AND: '&&';
OR: '||';
LPAREN: '(';
RPAREN: ')';
EQ: '==';
NEQ: '!=';
NULL: 'null';

STRING
   : '"' (ESC | SAFECODEPOINT)* '"'
   ;
fragment ESC
   : '\\' (["\\/bfnrt] | UNICODE)
   ;
fragment UNICODE
   : 'u' HEX HEX HEX HEX
   ;
fragment HEX
   : [0-9a-fA-F]
   ;
fragment SAFECODEPOINT
   : ~ ["\\\u0000-\u001F]
   ;

ID   : [a-zA-Z]+;

NUMBER
   : '-'? INT ('.' [0-9] +)? EXP?
   ;
fragment INT
   : '0' | [1-9] [0-9]*
   ;
fragment EXP
   : [Ee] [+\-]? INT
   ;

