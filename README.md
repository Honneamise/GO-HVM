# HVM : THE HONNY VIRTUAL MACHINE

---
## ARCHITECTURE
---
<br>

```
STATUS : 8 bit status flag

ACC : 8 bit accumulator

FL : 8 bit flag

R[0..F] : 16 registers, 8 bit each

MEM : 65535 bytes, first 2048 bytes reserved

PC : 16 bit program counter

SP : 8 bit stack pointer

BUS : 8 bit bi-directional buffered bus
```

---
## INSTRUCTIONS SET
---
<br>

8 bit opcodes

```
OP  : XX
h   : x-
l   : -x
```
flag
```
ZERO
POS
NEG
```
status
```
ST_RUN
ST_SUSPENDED
ST_HALTED
```
---
## REGISTERS OPERATIONS
---
<br>

## [LDI] 0x
Load Immediate
```
R[l] <- PC; PC++
```

## [LDM] 1x
Load From Memory
```
R[l] <- MEM[PC,PC+1]; PC+=2
```

## [STR] 2x
Store To Memory
```
MEM[PC,PC+1] <- R[l]; PC+=2
```

---
## COMBINED OPERATIONS
---
<br>

## [MOV] 3x
Move to Accumulator
```
ACC <- R[l]
```

## [SAV] 4x
Save Accumulator to Register
```
R[l] <- ACC
```

## [ADD] 5x
Add Register to Accumulator
```
ACC <- ACC + R[l]
```

## [SUB] 6x
Subtract Register from Accumulator
```
ACC <- ACC - R[l]
```

## [MUL] 7x
Multiply Accumulator with Register
```
ACC <- ACC * R[l]
```

## [DIV] 8x
Divide Accumulator by Register
```
ACC <- ACC / R[l]
```

## [AND] 9x
Accumulator logical AND with Register
```
ACC <- ACC & R[l]
```

## [OR] Ax
Accumulator logical OR with Register
```
ACC <- ACC | R[l]
```

---
## ACC OPERATIONS
---
<br>

## [INC] B0
Increase Accumulator by 1
```
ACC++
```

## [DEC] B1
Decrease Accumulator by 1
```
ACC--
```

## [LSH] B2
Logical left shift on Accumulator
```
ACC <- ACC<<1
```

## [RSH] B3
Logical right shift on Accumulator
```
ACC <- ACC>>1
```

## [NEG] B4
Negate the sign of Accumulator
```
ACC <- -ACC
```

## [BNEG] B5
Bitwise negate on Accumulator
```
ACC <- ~ACC
```

---
## FLOW CONTROL
---
<br>

## [BZ] C0
Branch if Zero
```
if FL==FL_ZERO then PC++
```

## [BP] C1
Branch if Positive
```
if FL==FL_POS then PC++
```

## [BN] C2
Branch if Negative
```
if FL==FL_NEG then PC++
```

## [CALL] C3
Call Routine 
```
MEM[SP] <- PC[high byte]; MEM[SP+1] <- PC[low byte]; SP+=2
PC[high byte] <- MEM[PC]; PC[low byte] <- MEM[PC+1]
```

## [RET] C4
Return from Routine
```
PC[high byte] <- SP-2; PC[low byte] <- SP-1; SP-=2; PC+=2
```

## [JMP] C5
Jump
```
PC[high byte] <- MEM[PC]; PC[low byte] <- MEM[PC+1]
```

---
## IO OPERATIONS
---
<br>

## [BUSR] D0
Bus Read
```
ACC <- BUS
```

## [BUSW] D1
Bus Write
```
BUS <- ACC
```

## [BUSCLR] D2
Bus Clear
```
reset the bus pending operations
```

---
## OTHER OPERATIONS
---
<br>

## [NOP] E0
No Operation
```
No op
```

## [HALT] E1
Halt
```
STATUS = ST_HALTED
```

---
## ASSEMBLY
---

**#**  is a comment, everything on the same line after this is ignored

**LABEL:** a label must be declared using a name followed by **:**

**REGISTER** can be any of the available registers R0,R1,R3,...

**VALUE** must be an hex value ( 1 or 2 bytes according to the instruction )

| Instruction | Bytes | Example |
| :--- | :----: | :--- |
| LDI | 2 | LDI REGISTER VALUE (1byte)
| LDM | 3 | LDI VALUE (2bytes)
| STR | 3 | STR VALUE (2bytes)
| MOV | 1 | MOV REGISTER
| SAV | 1 | SAV REGISTER
| ADD | 1 | ADD REGISTER
| SUB | 1 | SUB REGISTER
| MUL | 1 | MUL REGISTER
| DIV | 1 | DIV REGISTER
| AND | 1 | AND REGISTER
| OR | 1 | OR REGISTER
| INC | 1 | INC
| DEC | 1 | DEC
| LSH | 1 | LSH
| RSH | 1 | RSH
| NEG | 1 | NEG
| BNEG | 1 | BNEG
| BZ | 1 | BZ
| BP | 1 | BP
| BN | 1 | BN
| CALL | 3 | CALL LABEL
| RET | 3 | RET
| JMP | 1 | JUMP LABEL
| BUSR | 1 | BUSR
| BUSW | 1 | BUSW
| BUSCLR | 1 | BUSCLR
| NOP | 1 | NOP
| HALT | 1 | HALT


