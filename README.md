## GLOX

### The lox language interpreter implemented in go


### Get started

```
go build
./lox index.lox
```

and then you can see our lox program begin to fly.

-----

## Notes on [Crafting interpreters Part II](http://www.craftinginterpreters.com/contents.html):

## 词法分析

- Finite-state machine
- 乔姆斯基层次结构
- Lexical grammar: 决定了字符是如何被组成词语的。

正则语言又称正规语言是满足下述相互等价的一组条件的一类形式语言：

可以被确定有限状态自动机识别；
可以被非确定有限状态自动机识别；
可以被只读图灵机识别；
可以用正则表达式描述；
可以用正则文法生成。
可以用前缀文法生成。

lex 和 Flex 这种工具被用来自动生成词法分析器，扔给它们一些正则表达式，就可以自动生成完整的词法分析器。

## 语法分析

在词法分析器生成已归类单词后，语法分析器的任务是判断单词流表示的输入程序在该语言中是否是一个有效的句子。


### Context-Free Grammar 上下文无关语法

有很多种可以表达CFG的方法：传统的例如(Backus-Naur Form), BNF;

这里我们使用简单的语法：

- 产生式(规则)：CFG的每一条规则都叫做一个产生式，因为他们产生最后的语法串
- 非终结符：语法产生式中使用的语法变量：例如expr, paren0, etc.
- 终结符：出现在语句的单词。也就是词法分析中得到的token词素，会组成最终的字符串。
- `→`：表示推导；

用通配符简化的产生式：

```js
expr → expr ( "(" ( expr ( "," expr )* )? ")" | "." IDENTIFIER )*
     | IDENTIFIER
     | NUMBER
```

相当于：

原始产生式：

```js
expr → expr ( "(" ( expr ( "," expr )* )? ")" | "." IDENTIFIER )*
// expr → expr ( "(" ( expr ( "," expr )* )? ")" )*
expr → expr
expr → expr paren0
expr → expr paren1
expr → expr paren2

// ( "(" ")" )*
paren0 → "(" ")"
paren0 → "(" ")" paren0

// ( "(" expr ( "," expr )* ")" )*
paren1 → "(" expr comma ")"
paren1 → "(" expr comma ")" paren1

// paren → ( "(" expr ")" )*
paren2 → "(" expr ")"
paren2 → "(" expr ")" paren2

// ( "," expr )*
comma → "," expr
comma → "," expr comma

// expr → expr ("." IDENTIFIER )*
expr → expr
expr → expr "." IDENTIFIER

//  | IDENTIFIER
expr → IDENTIFIER
//  | NUMBER
expr → NUMBER
```

上述产生式可能表达的语法串：

- `123`
- `asd`
- `abc()`
- `abc()()`
- `abc(234, eee, ddd)`
- `abc(234)`
- `asd.asd`
- `asd.asd()`


### Lox的表达式语法规则：

__带有运算优先级的表达式CFG表示：__

参考C语言运算优先级，从上往下优先级依次递增：

```js
expression → sequence ;

sequence       → assignment ("," assignment)*
assignment → ( call "." )? IDENTIFIER "=" assignment
           | condition ;
condition    	 → logic_or ("?" condition ":" condition)?

logic_or   → logic_and ( "or" logic_and )* ;
logic_and  → equality ( "and" equality )* ;

equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | call ;
call           → primary ( "(" sequence? ")" | "." IDENTIFIER )* ;
primary        → NUMBER | IDENTIFIER | STRING | "this" | "false" | "true" | "nil" | func | "(" expression ")" | ("super" "." IDENTIFIER) ;
func → "fun" IDENTIFIER? "(" parameters? ")" block ;
parameters → IDENTIFIER ( "," IDENTIFIER )* ;
```

>下表列出 C 运算符的优先级和结合性。运算符从顶到底以降序列出。

[C 运算符的优先级](https://en.cppreference.com/w/c/language/operator_precedence)


语句：


```js
// 我们的程序就是无数个语句构成的。
program     → declaration* EOF ;

declaration → varDecl
            | classDecl
            | funDecl
            | statement ;

statement   → exprStmt
            | printStmt
            | ifStmt
            | whileStmt
            | returnStmt
            | forStmt
            | blockStmt ;

varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;
classDecl   → "class" IDENTIFIER ( "<" IDENTIFIER )? "{" "static"? function* "}" ;
funDecl  → "fun" function ;

function → IDENTIFIER "(" parameters? ")" block ;
parameters → IDENTIFIER ( "," IDENTIFIER )* ;

exprStmt  → expression ";" ;
printStmt → "print" expression ";" ;
blockStmt → "{" declaration* "}" ;

returnStmt → "return" expression? ";" ;

ifStmt    → "if" "(" expression ")" statement ( "else" statement )? ;
whileStmt → "while" "(" expression ")" statement ;
forStmt   → "for" "(" ( varDecl | exprStmt | ";" )
                      expression? ";"
                      expression? ")" statement ;
```

### 类和继承

类的实现方式通常有三种：

- classes
- prototypes
- multimethods(多分派)

类的主要功能：

- 暴露一个构造函数用来创建和初始化新的实例
- 提供一个方法用来存储和读取实例上的字段
- 定义被所有实例共享的一系列方法用来操控实例的状态

### 递归下降

递归下降是构造一个健壮的解析器的最简单的方法，不需要使用复杂的解析器生成器例如Yacc, Bison或者ANTLR。


### 错误处理

1. 尽可能的报告更多的错误信息。
2. 最小化连续的错误。


### 抽象语法树(自动生成)

将推导过程表示为图的树称为语法分析树。
 
### 访问者模式

访问者模式是一种将算法与对象结构分离的软件设计模式。

这个模式的基本想法如下：首先我们拥有一个由许多对象构成的对象结构，这些对象的类都拥有一个accept方法用来接受访问者对象；访问者是一个接口，它拥有一个visit方法，这个方法对访问到的对象结构中不同类型的元素作出不同的反应；在对象结构的一次访问过程中，我们遍历整个对象结构，对每一个元素都实施accept方法，在每一个元素的accept方法中回调访问者的visit方法，从而使访问者得以处理对象结构的每一个元素。我们可以针对对象结构设计不同的实在的访问者类来完成不同的操作。

访问者模式的使用场景
- 对象结构比较稳定，但经常需要在此对象结构上定义新的操作。
- 需要对一个对象结构中的对象进行很多不同的并且不相关的操作，而需要避免这些操作“污染”这些对象的类，也不希望在增加新操作时修改这些类。

角色介绍

* Visitor：接口或者抽象类，定义了对每个 Element 访问的行为，它的参数就是被访问的元素，它的方法个数理论上与元素的个数是一样的，因此，访问者模式要求元素的类型要稳定，如果经常添加、移除元素类，必然会导致频繁地修改 Visitor 接口，如果出现这种情况，则说明不适合使用访问者模式。
* ConcreteVisitor：具体的访问者，它需要给出对每一个元素类访问时所产生的具体行为。
* Element：元素接口或者抽象类，它定义了一个接受访问者（accept）的方法，其意义是指每一个元素都要可以被访问者访问。
* ElementA、ElementB：具体的元素类，它提供接受访问的具体实现，而这个具体的实现，通常情况下是使用访问者提供的访问该元素类的方法。
* ObjectStructure：定义当中所提到的对象结构，对象结构是一个抽象表述，它内部管理了元素集合，并且可以迭代这些元素提供访问者访问。

访问者模式的优点：
1. 各角色职责分离，符合单一职责原则
2. 具有优秀的扩展性
3. 如果需要增加新的访问者，增加实现类 ConcreteVisitor 就可以快速扩展。
使得数据结构和作用于结构上的操作解耦，使得操作集合可以独立变化
4. 灵活性

访问者模式的缺点：
1. 具体元素对访问者公布细节，违反了迪米特原则
2. 具体元素变更时导致修改成本大
3. 违反了依赖倒置原则，为了达到“区别对待”而依赖了具体类，没有以来抽象
访问者 visit 方法中，依赖了具体员工的具体方法。




















