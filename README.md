# Monkey language
This is Monkey language interpreter and compiler from Thorsten Ball's book.

### Assembly Compiler(WIP) 
- I am writing compiler from monkey to x64 assembly now.
  - treat intel CPU as stack machine (This is same approach to Monkey VM compiler of the compiler book.)
  - now, it doesn't almost use registers. 

#### example
```
$ cat sample/sample.mk 
let a = 2;
let b = 3;
let c = a * b / 2 ;

return c;
$ 
$ ./x64_gen sample/sample.mk | head -20
.intel_syntax noprefix

.text
.global main
main:
	mov rbp, rsp
	push 2
	push 3
	mov rax, [rbp-8]
	push rax
	mov rax, [rbp-16]
	push rax
	pop rbx
	pop rax
	imul rbx
	push rax
	push 2
	pop rbx
	pop rax
	cdq
$ 
$ 
$ ./x64_gen sample/sample.mk > /tmp/t.s; gcc /tmp/t.s 
$ ./a.out 
$ echo $?
2
$ 

```

### Reference
 - "Writing An Interpreter In Go".
 - "Writing A Compiler In Go".

