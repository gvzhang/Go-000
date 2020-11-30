# 学习笔记
- Packages that are reusable across many projects only return root error values.
- If the error is not going to be handled, wrap and return up the call stack.
- Once an error is handled, it is not alowed to be passed up the call stack any longer.

## 作业
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

### 思考
- dao层应该只负责读取数据，在dao层，数据不存在应该不算是错误。
- 数据是否存在，对业务的影响，应该由biz层做判断。
- dao层吞掉sql.ErrNoRows，其他层可减少对database/sql包的依赖。
