## GLOX

### The lox language interpreter implemented in go

## further reading

### 词法分析

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

### 语法分析

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

```js
expression → literal
           | unary
           | binary
           | grouping ;

literal    → NUMBER | STRING | "true" | "false" | "nil" ;
grouping   → "(" expression ")" ;
unary      → ( "-" | "!" ) expression ;
binary     → expression operator expression ;
operator   → "==" | "!=" | "<" | "<=" | ">" | ">="
           | "+"  | "-"  | "*" | "/" ;
```

### 抽象语法树(自动生成)
 
### 访问者模式

访问者模式是一种将算法与对象结构分离的软件设计模式。

这个模式的基本想法如下：首先我们拥有一个由许多对象构成的对象结构，这些对象的类都拥有一个accept方法用来接受访问者对象；访问者是一个接口，它拥有一个visit方法，这个方法对访问到的对象结构中不同类型的元素作出不同的反应；在对象结构的一次访问过程中，我们遍历整个对象结构，对每一个元素都实施accept方法，在每一个元素的accept方法中回调访问者的visit方法，从而使访问者得以处理对象结构的每一个元素。我们可以针对对象结构设计不同的实在的访问者类来完成不同的操作。




















