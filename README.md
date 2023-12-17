Spent a day to create this small app, it can read a math expression as input argument and get a calculation result also print expression AST created in the middle

Example:
```
./calc '6+2*(3+8)-4+27/3/3 + 37%9'
The generated math expression AST tree:
(L + R)
(L + R) (L % R)
(L - R) (L / R) ( 37 )  ( 9 )
(L + R) ( 4 )   (L / R) ( 3 )
( 6 )   (L * R) ( 27 )  ( 3 )
( 2 )   (L + R)
( 3 )   ( 8 )
The calculation result of math expression '6+2*(3+8)-4+27/3/3 + 37%9' is 31.00
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
