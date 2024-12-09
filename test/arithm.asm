arithm  START   0

first   LDA     x
        ADD     y
        STA     sum
        . ---------
        LDA     x
        SUB     y
        STA     diff
        . ---------
        LDA     x
        MUL     y
        STA     prod
        . ---------
        LDA     x
        DIV     y
        STA     quot
        . ---------
        MUL     y
        SUB     x
        MUL     minus
        STA     mod
halt    J       halt

. podatki
x       WORD    14
y       WORD    4
minus   WORD    -1
sum     WORD    0
diff    WORD    0
prod    WORD    0
quot    WORD    0
mod     WORD    0
        END     first