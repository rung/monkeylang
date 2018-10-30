# Monkey language
- This is Monkey language interpreter and compiler from Thorsten Ball's book.
  - and original x64 compiler from Monkey bytecode to x64 assembly.

### Assembly Compiler(WIP) 
- I am writing compiler from monkey to x64 assembly now.
  - treat intel CPU as stack machine (This is same approach to Monkey VM compiler of the compiler book.)
  - now, it doesn't almost use registers. 

#### example
```bash
$ cat sample/function.mk 
let a = 1;

let fnA = fn() {
	let c = 3;
	return c;
}

let fnB = fn() {
	let c = 6;
	return c;
}

return a + fnA() + fnB()
```

<details>
<summary>x64 assembly code</summary>
<pre>
<code>

$ ./x64_gen sample/function.mk
.intel_syntax noprefix

.text
.global main
main:
	push rbp
	mov rbp, rsp
	sub rsp, 24
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
	mov rax, [rbp-8]
	push rax
	mov rax, [rbp-16]
	push rax
	pop rax
	call rax
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	mov rax, [rbp-24]
	push rax
	pop rax
	call rax
	push rax
	pop rbx
	pop rax
	add rax, rbx
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


$ 

</code>
</pre>
</details>



```bash
$ ./x64_gen sample/function.mk > /tmp/t.s; gcc /tmp/t.s; ./a.out
$ echo $?
10
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

