#include "textflag.h"

//func U64toByte_64_asm(a uint64, b []byte) 
TEXT ·U64toByte_64_asm(SB), NOSPLIT, $0 
    MOVQ a+0(FP), AX 
    MOVQ b+8(FP), SI
    
    BSWAPQ AX   
    MOVQ AX, (8*0)(SI) 
    
RET

//func BytetoU64_64(a []byte)(b uint64)
TEXT ·BytetoU64_64(SB), NOSPLIT, $0 
    MOVQ a+0(FP), SI  
    MOVQ (8*0)(SI), AX
    
    BSWAPQ AX  
    MOVQ AX, b+24(FP)        
RET

//func U64toByte_256_asm(a [4]uint64, b []byte) 
TEXT ·U64toByte_256_asm(SB), NOSPLIT, $0 
    MOVQ a+0(FP), R8
    MOVQ a+8(FP), R9
    MOVQ a+16(FP), R10
    MOVQ a+24(FP), R11 
    MOVQ b+32(FP), SI
    
    BSWAPQ R8
    BSWAPQ R9
    BSWAPQ R10
    BSWAPQ R11
       
    MOVQ R11, (8*0)(SI) 
    MOVQ R10, (8*1)(SI)
    MOVQ R9, (8*2)(SI)
    MOVQ R8, (8*3)(SI)
    
RET

//BytetoU64_256(x []byte) (y [4]uint64)
TEXT ·BytetoU64_256(SB), NOSPLIT, $0 
    MOVQ x+0(FP), SI

    MOVQ (8*0)(SI), R8
    MOVQ (8*1)(SI), R9
    MOVQ (8*2)(SI), R10
    MOVQ (8*3)(SI), R11 
        
    BSWAPQ R8
    BSWAPQ R9
    BSWAPQ R10
    BSWAPQ R11
       
    MOVQ R11, y+24(FP) 
    MOVQ R10, y+32(FP)
    MOVQ R9,  y+40(FP)
    MOVQ R8,  y+48(FP)
    
RET

//U32toByte_256_asm(a [8]uint32, b []byte)
TEXT ·U32toByte_256_asm(SB), NOSPLIT, $0 
    MOVQ a+0(FP), R8
    MOVQ a+8(FP), R9
    MOVQ a+16(FP), R10
    MOVQ a+24(FP), R11 
    MOVQ b+32(FP), SI
    
    BSWAPQ R8
    BSWAPQ R9
    BSWAPQ R10
    BSWAPQ R11
    
    MOVL R8, (4*1)(SI)
    SHRQ $32, R8
    MOVL R8, (4*0)(SI)   
    MOVL R9, (4*3)(SI)
    SHRQ $32, R9
    MOVL R9, (4*2)(SI)
    MOVL R10, (4*5)(SI)
    SHRQ $32, R10
    MOVL R10, (4*4)(SI)
    MOVL R11, (4*7)(SI)
    SHRQ $32, R11
    MOVL R11, (4*6)(SI)
    
RET

//BytetoU32_256(x []byte) (y [8]uint32)
TEXT ·BytetoU32_256(SB), NOSPLIT, $0 
    MOVQ x+0(FP), SI

    MOVQ (8*0)(SI), R8
    MOVQ (8*1)(SI), R9
    MOVQ (8*2)(SI), R10
    MOVQ (8*3)(SI), R11 
        
    BSWAPQ R8
    BSWAPQ R9
    BSWAPQ R10
    BSWAPQ R11
       
    MOVL R8, y+28(FP)
    SHRQ $32, R8
    MOVL R8, y+24(FP) 
    MOVL R9, y+36(FP)
    SHRQ $32, R9
    MOVL R9, y+32(FP) 
    MOVL R10, y+44(FP)
    SHRQ $32, R10
    MOVL R10, y+40(FP)
    MOVL R11, y+52(FP)
    SHRQ $32, R11
    MOVL R11, y+48(FP)
    
RET

//U32toByte_128_asm(a [4]uint32, b []byte)
TEXT ·U32toByte_128_asm(SB), NOSPLIT, $0 
    MOVQ a+0(FP), R8
    MOVQ a+8(FP), R9
    MOVQ b+16(FP), SI
    
    BSWAPQ R8
    BSWAPQ R9
    
    MOVL R8, (4*1)(SI)
    SHRQ $32, R8
    MOVL R8, (4*0)(SI)   
    MOVL R9, (4*3)(SI)
    SHRQ $32, R9
    MOVL R9, (4*2)(SI)
    
RET

//BytetoU32_128(x []byte) (y [4]uint32)
TEXT ·BytetoU32_128(SB), NOSPLIT, $0 
    MOVQ x+0(FP), SI

    MOVQ (8*0)(SI), R8
    MOVQ (8*1)(SI), R9
        
    BSWAPQ R8
    BSWAPQ R9
       
    MOVL R8, y+28(FP)
    SHRQ $32, R8
    MOVL R8, y+24(FP) 
    MOVL R9, y+36(FP)
    SHRQ $32, R9
    MOVL R9, y+32(FP) 
    
RET

//U64toU32_256_asm(a [4]uint64) (b [8]uint32) 
TEXT ·U64toU32_256(SB), NOSPLIT, $0 
    MOVQ x+0(FP), AX
    MOVL AX, b+60(FP)
    SHRQ $32, AX
    MOVL AX, b+56(FP)   
    
    MOVQ x+8(FP), AX
    MOVL AX, b+52(FP)
    SHRQ $32, AX
    MOVL AX, b+48(FP) 
    
    MOVQ x+16(FP), AX
    MOVL AX, b+44(FP)
    SHRQ $32, AX
    MOVL AX, b+40(FP) 
    
    MOVQ x+24(FP), AX
    MOVL AX, b+36(FP)
    SHRQ $32, AX
    MOVL AX, b+32(FP) 
    
RET

//U32toU64_256_asm(a [8]uint32) (b [4]uint64) 
TEXT ·U32toU64_256(SB), NOSPLIT, $0 
    MOVL a+0(FP), AX
    MOVL AX, b+60(FP)
    MOVL a+4(FP), AX
    MOVL AX, b+56(FP)
    
    MOVL a+8(FP), AX
    MOVL AX, b+52(FP)
    MOVL a+12(FP), AX
    MOVL AX, b+48(FP)
    
    MOVL a+16(FP), AX
    MOVL AX, b+44(FP)
    MOVL a+20(FP), AX
    MOVL AX, b+40(FP)
    
    MOVL a+24(FP), AX
    MOVL AX, b+36(FP)
    MOVL a+28(FP), AX
    MOVL AX, b+32(FP)   
    
RET

//BytetoU32_512(x []byte) (y [16]uint32)
TEXT ·BytetoU32_512(SB), NOSPLIT, $0 
    MOVQ x+0(FP), SI

    MOVQ (8*0)(SI), R8
    MOVQ (8*1)(SI), R9
    MOVQ (8*2)(SI), R10
    MOVQ (8*3)(SI), R11 
    MOVQ (8*4)(SI), R12
    MOVQ (8*5)(SI), R13
    MOVQ (8*6)(SI), R14
    MOVQ (8*7)(SI), R15   
        
    BSWAPQ R8
    BSWAPQ R9
    BSWAPQ R10
    BSWAPQ R11
    BSWAPQ R12
    BSWAPQ R13
    BSWAPQ R14
    BSWAPQ R15
       
    MOVL R8, y+28(FP)
    SHRQ $32, R8
    MOVL R8, y+24(FP) 
    MOVL R9, y+36(FP)
    SHRQ $32, R9
    MOVL R9, y+32(FP) 
    MOVL R10, y+44(FP)
    SHRQ $32, R10
    MOVL R10, y+40(FP)
    MOVL R11, y+52(FP)
    SHRQ $32, R11
    MOVL R11, y+48(FP)
    MOVL R12, y+60(FP)
    SHRQ $32, R12
    MOVL R12, y+56(FP)
    MOVL R13, y+68(FP)
    SHRQ $32, R13
    MOVL R13, y+64(FP)
    MOVL R14, y+76(FP)
    SHRQ $32, R14
    MOVL R14, y+72(FP)
    MOVL R15, y+84(FP)
    SHRQ $32, R15
    MOVL R15, y+80(FP)
RET
























