prog    START 0

        LDA #stack
        STA sp

        LDB @sp

halt    J halt








sp WORD 0

stack WORD 7