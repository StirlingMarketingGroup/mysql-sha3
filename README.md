# MySQL sha3

A small MySQL UDF library for creating sha3 hashes, written in Golang.

---

### `sha3`

Returns a sha3 hash, hex encoded.

```sql
`sha3` ( `message` , `bits` )
```
- `` `message` ``
  - The string to be hashed
- `` `bits` ``
  - The output size in bits

## Examples

```sql
select`sha3`('yeet', 224);
-- '243bbdbe21fc7f9b1bc3cd080435ebe99e2222684be94a391096a220'

select`sha3`('yeet', 256);
-- '5badd143cdd19064a3c72cd3f1af76c8e56b380f8619ff0fb39c1f877295cff1'

select`sha3`('yeet', 384);
-- '801d08645e906ac189dfe6051e060889c5d9eb1df33d6200a3a03284bf1ffe15fe9832e99c8485388c74af253a9fef7f'

select`sha3`('yeet', 512);
-- '60f639ce101c1f363429ea4e5ce5546cf6c08b45f25cccfce1eba794a41d3bc5dfa391d083fea2a5cbb823c7dfb57ca4b83b2ddb1d99e54b06100986ad5db6c3'

```
---
### `unhex_sha3`

Returns a sha3 hash, not encoded. In MySQL, doing `unhex(sha2(...` appears to be optomized internally, since it takes no time difference to do that vs `sha2(...`, however in our case, unhexing the output from `sha3` definitely adds overhead, I assume since we're encoding in the UDF and then decoding for no reason. 

**NOTE**: If/when the official sha3 functions come into MySQL, they will gracefully be used over our UDF for `sha3`, although that is not the case for this functions since there will be no replacement `unhex_sha3` function.

```sql
`unhex_sha3` ( `message` , `bits` )
```
- `` `message` ``
  - The string to be hashed
- `` `bits` ``
  - The output size in bits

## Examples

```sql
select`unhex_sha3`('yeet', 224);
-- 0x243BBDBE21FC7F9B1BC3CD080435EBE99E2222684BE94A391096A220

select`unhex_sha3`('yeet', 256);
-- 0x5BADD143CDD19064A3C72CD3F1AF76C8E56B380F8619FF0FB39C1F877295CFF1

select`unhex_sha3`('yeet', 384);
-- 0x801D08645E906AC189DFE6051E060889C5D9EB1DF33D6200A3A03284BF1FFE15FE9832E99C8485388C74AF253A9FEF7F

select`unhex_sha3`('yeet', 512);
-- 0x60F639CE101C1F363429EA4E5CE5546CF6C08B45F25CCCFCE1EBA794A41D3BC5DFA391D083FEA2A5CBB823C7DFB57CA4B83B2DDB1D99E54B06100986AD5DB6C3

```
---

## Dependencies

You will need Golang, which you can get from here https://golang.org/doc/install. You will also need the MySQL dev library.

Debian / Ubuntu
```shell
sudo apt update
sudo apt install libmysqlclient-dev
```
## Installing

You can find your MySQL plugin directory by running this MySQL query

```sql
select @@plugin_dir;
```

then replace `/usr/lib/mysql/plugin` below with your MySQL plugin directory.

```shell
cd ~ # or wherever you store your git projects
git clone https://github.com/StirlingMarketingGroup/mysql-sha3.git
cd mysql-sha3
go get -d ./...
go build -buildmode=c-shared -o sha3.so
sudo cp sha3.so /usr/lib/mysql/plugin/ # replace plugin dir here if needed
```

Enable the functions in MySQL by running this MySQL query

```sql
create function`Sha3`returns string soname'sha3.so';
create function`unhex_sha3`returns string soname'sha3.so';
```