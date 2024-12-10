first   START   0

        JSUB    horner
           
        
halt    J       halt

. podatki
x       WORD    2
val     WORD    1
stack   RESW    1
in      WORD    0
        WORD    5
        WORD    42
lastin  EQU     *
len     EQU     lastin - in


horner  STA     stack
        LDA     val
        MUL     x
        ADD     #2
        MUL     x
        ADD     #3
        MUL     x
        ADD     #4
        MUL     x
        ADD     #5
        STA     val
        LDA     stack
        RSUB


        END     first
