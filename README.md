# jsonschemaline
### 计划
将 validate_schema、output_schema 改成 key:value 格式(类似go tag)，提升可阅读性和维护效率
出参案例：
```
{"$schema":"http://json-schema.org/draft-07/schema#","$id":"execAPIXyxzBlacklistInfoPaginate","properties":{"items":{"type":"array","items":{"type":"object","properties":{"id":{"type":"string","src":"PaginateOut.#.Fid"},"openId":{"type":"string","src":"PaginateOut.#.Fopen_id"},"type":{"type":"string","src":"PaginateOut.#.Fopen_id_type"},"status":{"type":"string","src":"PaginateOut.#.Fstatus"}},"required":["id","openId","type","status"]}},"pageInfo":{"type":"object","properties":{"pageIndex":{"type":"string","src":"input.pageIndex"},"pageSize":{"type":"string","src":"input.pageSize"},"total":{"type":"string","src":"PaginateTotalOut"}},"required":["pageIndex","pageSize","total"]}},"type":"object","required":["items","pageInfo"]}
```
转换为：
```
fullname:items[].id,src:PaginateOut.#.Fid,type:string,required:true
fullname:items[].openId,src:PaginateOut.#.Fopen_id,type:string,required:true
fullname:items[].type,src:PaginateOut.#.Fopen_id_type,type:string,required:true
fullname:items[].status,src:PaginateOut.#.Fstatus,type:string,required:true
fullname:pageInfo.pageIndex,src:input.pageIndex,type:string,required:true
fullname:pageInfo.pageSize,src:input.pageSize,type:string,required:true
fullname:pageInfo.total,src:PaginateTotalOut,type:string,required:true
```

入参案例：
```
{"$schema":"http://json-schema.org/draft-07/schema#","$id":"execAPIXyxzBlacklistInfoInsert","properties":{"config":{"type":"object","properties":{"openId":{"type":"string","format":"DBValidate"},"type":{"type":"string","format":"number","enum":["1","2"]},"status":{"type":"string","format":"number","enum":["0","1"]}},"required":["openId","type","status"]}},"type":"object"}
```
转换为：
```
fullname:config.openId,dst:FopenID,format:DBValidate,type:string,required:true
fullname:config.type,dst:FopenIDType,enum:["1","2"],type:string,required:true
fullname:config.status,dst:Fstatus,enum:["0","1"],type:string,required:true

```
分页入参案例:
```
{"$schema":"http://json-schema.org/draft-07/schema#","$id":"execAPIInquiryScreenIdentifyPaginate","properties":{"pageSize":{"type":"string","format":"number"},"pageIndex":{"type":"string","format":"number"}},"type":"object","required":["pageSize","pageIndex"]}
```
转换为:
```
fullname:pageSize,dst:limit,format:number,type:string,required:true
fullname:pageIndex,dst:Offset,format:number,type:string,required:true,tpl:{{setValue . "Offset" (mul  (getValue .  "input.pageIndex")   (getValue . "input.pageSize"))}}
```
dst、tpl 字段会提炼出，动态生成 template内容
通过这种转换后，更易于书写和阅读，和计划中的文档格式更相似，同时dst、tpl等字段定义优化了值转换的管理为自动校验提供可行的机制，后续api 的 exec 字段可能被弃用

有偿功能实现,价格500左右,
概念:把
```
fullname:xxx,type:xxx,required:true ...
fullname:xxx,type:xxx,required:true ...
```
的格式称为jsonschema 行列式
一行代表一个jsonschema 元素的描述,其格式为key:value,... 其中key包含fullname、src、dst tpl 4个特殊值,其余的为标准的jsonschema 的属性(如 type,required、enum 等关键词) 每对key:value 之间使用英文","分割.顺序无要求,除了fullname、type 属性必须之外,其它属性可有可无,特殊属性描述:
fullname 必须存在,由[a-zA-z_]以及特殊字符串.和[]组成.其中"."代表json对象,"[]"代表数组 如"fullname:userList[].name" 代表 userList是数组格式,userList数组元素为对象,数组元素对象有个name属性
src、dst 两个属性无特殊意义,仅仅是给标准jsonschema增加的自定义属性
需要实现的功能:
标准jsonschema和jsonschema 行列式 互转


案例1:
jsonschema:
```
{"$schema":"http://json-schema.org/draft-07/schema#","properties":{"config":{"type":"object","properties":{"openId":{"type":"string","format":"DBValidate"},"type":{"type":"string","format":"number","enum":["1","2"]},"status":{"type":"string","format":"number","enum":["0","1"]}},"required":["openId","type","status"]}},"type":"object"}
```
对应 jsonschema 行列式

```
fullname:config.openId,dst:FopenID,format:DBValidate,type:string,required:true
fullname:config.type,dst:FopenIDType,enum:["1","2"],type:string,required:true
fullname:config.status,dst:Fstatus,enum:["0","1"],type:string,required:true

```

案例2：
jsonschema 行列式
```
fullname:config.openId,dst:FopenID,format:DBValidate,type:string,required:true
fullname:config.type,dst:FopenIDType,enum:["1","2"],type:string,required:true
fullname:config.status,dst:Fstatus,enum:["0","1"],type:string,required:true
```
jsonschema格式
```
{"$schema":"http://json-schema.org/draft-07/schema#","properties":{"config":{"type":"object","properties":{"openId":{"type":"string","format":"DBValidate"},"type":{"type":"string","format":"number","enum":["1","2"]},"status":{"type":"string","format":"number","enum":["0","1"]}},"required":["openId","type","status"]}},"type":"object"}
```
