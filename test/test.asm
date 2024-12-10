test    START 0
        LDA #7

        COMP #8

        JGT out
halt    J halt

out     LDB #2
        J halt

