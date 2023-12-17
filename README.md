Spent a day to create this small app, it can read a math expression as input argument and get a calculation result also print expression AST created in the middle

`Example:

 ./calc '34 + 56 * (3 + 8) + 39/3 + 57 % 12'
 
 Expression AST: {Operator:+ Lhs:{Operator:+ Lhs:{Operator:+ Lhs:{Val:34} Rhs:{Operator:* Lhs:{Val:56} Rhs:{Operator:+ Lhs:{Val:3} Rhs:{Val:8}}}} Rhs:{Operator:/ Lhs:{Val:39} Rhs:{Val:3}}} Rhs:{Operator:% Lhs:{Val:57} Rhs:{Val:12}}}
 
The calculation result of math expression '34 + 56 * (3 + 8) + 39/3 + 57 % 12' is 667.00
`

