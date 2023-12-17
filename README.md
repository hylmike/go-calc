Spent a day to create this small app, it can read a math expression as input argument and get a calculation result also print expression AST created in the middle

Example:
```
./calc '-10+2*(3+81)-4*(9-3*(2+3))+27/3/3+137%9'
The generated math expression AST tree:
(L + R)
(L + R) (L % R)
(L - R) (L / R) ( 137 ) ( 9 )
(L + R) (L * R) (L / R) ( 3 )
( -10 ) (L * R) ( 4 )   (L - R) ( 27 )  ( 3 )
( 2 )   (L + R) ( 9 )   (L * R)
( 3 )   ( 81 )  ( 3 )   (L + R)
( 2 )   ( 3 )
The calculation result of math expression '-10+2*(3+81)-4*(9-3*(2+3))+27/3/3+137%9' is 200.00
```

If expression is invalid, will print error with error position

Example:
```
./calc '6+2*(3+)-4'
Invalid expression:  strconv.ParseFloat: parsing ")": invalid syntax
Should be '(' or '0-9' but get ')'
----------
6+2*(3+)-4
       ^
----------
```
