# Monkey language
- This is Monkey language interpreter and compiler from Thorsten Ball's book.
  - and original x64 compiler from Monkey bytecode to x64 assembly.

### Assembly Compiler(WIP) 
- I am writing compiler from monkey to x64 assembly now.
  - treat intel CPU as stack machine (This is same approach to Monkey VM compiler of the compiler book.)
  - now, it doesn't almost use registers. 

#### Sample
##### Monkey code
```bash
$ cat sample/sample.mk
let a = 1;

let fnA = fn() {
	let c = 3;
	return c;
}

let fnB = fn() {
	let c = 6;
	return c;
}

if (fnA() + fnB() + a != 10){
	let d = "dummy\n";
	puts(d);
} else {
	let d = "Hello World!\n";
	puts(d);
}

return fnA();

```

##### output x64 assembly
<details>
<summary>open</summary>
<pre>
<code>

$ ./x64_gen sample/sample.mk
.intel_syntax noprefix

.text
.section	.rodata
.STRGBL6:
	.string "dummy\n"
.STRGBL7:
	.string "Hello World!\n"

.text
.global main
main:
	push rbp
	mov rbp, rsp
	sub rsp, 40
	push 1
	pop rax
	mov [rbp-8] ,rax
	lea rax, function1[rip]
	push rax
	pop rax
	mov [rbp-16] ,rax
	lea rax, function2[rip]
	push rax
	pop rax
	mov [rbp-24] ,rax
	mov rax, [rbp-16]
	push rax
	mov rax, [rsp+0]
	call rax
	add rsp, 8
	push rax
	mov rax, [rbp-24]
	push rax
	mov rax, [rsp+0]
	call rax
	add rsp, 8
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	mov rax, [rbp-8]
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	push 10
	pop rax
	pop rbx
	cmp rax, rbx
	jne .LABEL0
	push 1
	jmp .LABEL1
.LABEL0:
	push 0
.LABEL1:
	pop rax
	cmp rax, 0
	jne .LABEL2
	lea rax, .STRGBL6[rip]
	push rax
	pop rax
	mov [rbp-32] ,rax
	lea rax, puts[rip]
	push rax
	mov rax, [rbp-32]
	push rax
	mov rax, [rsp+8]
	call rax
	add rsp, 16
	push rax
	jmp .LABEL3
.LABEL2:
	lea rax, .STRGBL7[rip]
	push rax
	pop rax
	mov [rbp-40] ,rax
	lea rax, puts[rip]
	push rax
	mov rax, [rbp-40]
	push rax
	mov rax, [rsp+8]
	call rax
	add rsp, 16
	push rax
.LABEL3:
	pop rax
	mov rax, [rbp-16]
	push rax
	mov rax, [rsp+0]
	call rax
	add rsp, 8
	push rax
	pop rax
	mov rsp, rbp
	pop rbp
	ret

.global function1
function1:
	push rbp
	mov rbp, rsp
	sub rsp, 8
	push 3
	pop rax
	mov [rbp-8] ,rax
	mov rax, [rbp-8]
	push rax
	pop rax
	mov rsp, rbp
	pop rbp
	ret

.global function2
function2:
	push rbp
	mov rbp, rsp
	sub rsp, 8
	push 6
	pop rax
	mov [rbp-8] ,rax
	mov rax, [rbp-8]
	push rax
	pop rax
	mov rsp, rbp
	pop rbp
	ret

.global puts
puts:
	#header
	push rbp
	mov rbp, rsp

	#strlen start
	xor rcx, rcx
	sub rcx, 1

	xor rax, rax
	mov rdx, [rbp+16]
	mov bl, [rdx]
	cmp bl, 0
	je .L1
.L0:
	add rdx, 1
	mov bl, [rdx]
	cmp bl, 0
	loopne .L0
.L1:
	not rcx
	mov rdx, rcx
	# strlen end

	# write(1, "string", strlen) // printf
	mov rax, 1
	mov rdi, 1
	mov rsi, [rbp+16]
	syscall
	push rax

	#footar
	mov rsp, rbp
	pop rbp
	ret


$

</code>
</pre>
</details>


##### Assemble(by gcc) and Execution

```bash
$ ./x64_gen sample/sample.mk > /tmp/t.s; gcc /tmp/t.s -o /tmp/t; /tmp/t
Hello World!
$ echo $?
3
$

```

#### support
- type: integer and string
  - let a = 1;
  - puts("Hello world!");
- local/global binding
  - let a = 1;
- calculate
  - 1 + 1 / 3 * 5 - 1
- function with arguments(definiction, call)
  - let f = fn(a, b, c){return a + b + c + 4;} return f(1, 2, 3);
- builtin function
  - puts("hello");
- if-then-else statement
  - if(1 > a){ return 1;}

#### unsupport

### Reference
 - "Writing An Interpreter In Go".
 - "Writing A Compiler In Go".

