#### go_web_api.user 
用户表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id |  | bigint unsigned | PRI | NO | auto_increment |  |
| 2 | username | 用户名 | varchar(100) | UNI | NO |  |  |
| 3 | password |  | varchar(100) |  | NO |  |  |
| 4 | mobile |  | varchar(20) |  | NO |  |  |
| 5 | is_deleted |  | tinyint |  | NO |  | 0 |
| 6 | created_time |  | timestamp |  | NO | DEFAULT_GENERATED | CURRENT_TIMESTAMP |
| 7 | update_time |  | timestamp |  | NO | DEFAULT_GENERATED | CURRENT_TIMESTAMP |
