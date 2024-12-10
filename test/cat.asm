prog    START 0

loop    RD #0
        SUB #48
        COMP #0
        JEQ halt
        ADD #48
        WD #1
        J loop


halt    J halt


        END prog