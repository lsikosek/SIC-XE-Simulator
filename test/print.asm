puts    START 0
        CLEAR X 
loop    LDCH txt, X . nalozi en bajt v spodnji bajt registra a
        JSUB putc             .WD #1 . write device 1

        TIX #len
        JLT loop

halt    J halt

. rutina, ki izpise znak v registru A
putc    WD #170
        RSUB

txt     BYTE C'SIC/XE'
lastin  EQU *
len     EQU lastin - txt

        END puts