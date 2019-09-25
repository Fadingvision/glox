## GLOX

### The lox language interpreter implemented in go

## further reading

### Token

- Finite-state machine
- 乔姆斯基层次结构
- Lexical grammar: 决定了哪些字符被认为是词位的语法规则

正则语言又称正规语言是满足下述相互等价的一组条件的一类形式语言：

可以被确定有限状态自动机识别；
可以被非确定有限状态自动机识别；
可以被只读图灵机识别；
可以用正则表达式描述；
可以用正则文法生成。
可以用前缀文法生成。

lex 和 Flex 这种工具被用来自动生成词法分析器，扔给它们一些正则表达式，就可以自动生成完整的词法分析器。