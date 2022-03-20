# Direct Connection in Backend package

## Code Explanation

### The first step

> The first step is to send the initial handshake packet from MariaDB to Gaea.

There are some details about the initial handshake packet in [the official document](https://mariadb.com/kb/en/connection/), and please see the details below.

<img src="./assets/image-20220315221559157.png" alt="image-20220315221559157" style="zoom:100%;" /> 

The actual packet demonstrates how this handshake works, and please see details below.

| packet                          | exmaple                                                      |
| ------------------------------- | ------------------------------------------------------------ |
| int<1> protocol version         | Protocol Version 10                                          |
| string<NUL> server version      | MariaDB version is <br /><br />[]uint8{<br />53, 46, 53, 46, 53,<br />45, 49, 48, 46, 53,<br />46, 49, 50, 45, 77,<br />97, 114, 105, 97, 68,<br />66, 45, 108, 111, 103<br />}<br /><br />Converting the array to ASCII, the result is "5.5.5-10.5.12-MariaDB-log". |
| int<4> connection id            | Connection ID is []uint8{16, 0, 0, 0}.<br /><br />After reversing the array, it becomes []uint8{0, 0, 0, 16} that equals to []uint32{16}. |
| string<8> scramble 1st part     | The first part of the scramble:<br /><br />MariaDB utilizes the scramble for secure password authentication.<br /><br />The scramble is 20 bytes of data; the first part occupies 8 bytes, []uint8{81, 64, 43, 85, 76, 90, 97, 91}. |
| string<1> reserved byte         | It occupies 1 byte, []uint8{0}.                              |
| int<2> server capabilities      | The first part of the capability occupies 2 bytes,  []uint8{254, 247}. |
| int<1> server default collation | The charset of MariaDB in the current exameple is 33.<br /><br />After checking<br />[character-sets-and-collations](https://mariadb.com/kb/en/supported-character-sets-and-collations/)<br />or<br />using a command "SHOW CHARACTER SET LIKE 'utf8'",<br />finding out that number 33 means "utf8_general_ci". |
| int<2> status flags             | The status of MariaDB in the current exameple is []uint8{2, 0}.<br/><br />Reversing from the status flags to []uint8{0, 2} and then converting them to binary, []uint{0b000000000, 0b00000010}.<br /><br />After referring to "Gaea/mysql/constants.go", the result means "Autocommit." |
| int<2> server capabilities      | The second part of the capability occupies 2 bytes,  []uint8{255, 129}. |

Calculate the whole capability

```
Gathering two parts of the capability and combining them, the result is []uint8{254, 247, 255, 129}.

After Converting the result to binary, the result becomes []uint8{0b10000001, 0b11111111, 0b11110111, 0b11111110}.

After that, refer to https://mariadb.com/kb/en/connection/ and ensure some details without difficulty.

For example, the first element of the capability is 0, which means the packet came from MariaDB to Gaea.
```

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & PLUGIN_AUTH)<br/>        int<1> plugin data length <br/>    else<br/>        int<1> 0x00 |
| example | skip 1 byte                                                  |

The next table follows on from the previous one.

| item    | value            |
| ------- | ---------------- |
| packet  | string<6> filler |
| example | skip 6 bytes     |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & CLIENT_MYSQL)<br/>        string<4> filler <br/>    else<br/>        int<4> server capabilities 3rd part .<br />        MariaDB specific flags /* MariaDB 10.2 or later */ |
| example | skip 4 bytes                                                 |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & CLIENT_SECURE_CONNECTION)<br/>        string<n> scramble 2nd part . Length = max(12, plugin data length - 9)<br/>        string<1> reserved byte |
| example | The scramble is 20 bytes of data; the second part occupies 12 (20-8=12) bytes, []uint8{34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}. |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & PLUGIN_AUTH)<br/>        string<NUL> authentication plugin name |
| example | Gaea discards the rest of the data in the packet because there is no use for "PLUGIN_AUTH". |

combine the whole data of the scramble:

```
The first part of the scramble is []uint8{81, 64, 43, 85, 76, 90, 97, 91}
The second part of the scramble is []uint8{34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}

After combining them, the final result is []uint8{81, 64, 43, 85, 76, 90, 97, 91, 34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}.
```

### The second step

> The second step is to calculate the auth base on the scramble, combined with two parts of the scramble.

There are some details about the auth formula in [the official document](https://dev.mysql.com/doc/internals/en/secure-password-authentication.html).

```
some formulas for the auth:

SHA1( password ) XOR SHA1( "20-bytes random data from server" <concat> SHA1( SHA1( password ) ) )
    stage1 = SHA1( password )
    stage1Hash = SHA1( stage1 ) = SHA1( SHA1( password ) )
    scramble = SHA1( scramble <concat> SHA1( stage1Hash ) ) // the first new scramble
    scramble = stage1 XOR scramble // the second new scramble
```

Assume

- The password for a secure login process in MariaDB is 12345.
- The auth base on the scramble, combined with two parts of the scramble, is []uint8{81, 64, 43, 85, 76, 90, 97, 91, 34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}.
  The result that converted from decimal to hexadecimal is []uint8{51, 40, 2B, 55, 4c, 5a, 61, 5b, 22, 35, 24, 55, 5d, 56,  75,  69, 31, 57, 41,  7d}.
  It is the same as  51402B554c5A615b223524555d5675693157417d.

Regarding the stage1 formula, Linux Bash calculates the result and compares it.

```bash
# stage1 = SHA1( password )

# calculate stage1
$ echo -n 12345 | sha1sum | head -c 40 # convert password 12345 to stage1
8cb2237d0679ca88db6464eac60da96345513964 # stage1
```

As regards the stage1Hash formula, Linux Bash calculates the result and compares it.

```bash
# stage1Hash = SHA1( stage1 ) = SHA1( SHA1( password ) )

$ echo -n 12345 | sha1sum | xxd -r -p | sha1sum | head -c 40
00a51f3f48415c7d4e8908980d443c29c69b60c9 # stage1hash

$ echo -n 8cb2237d0679ca88db6464eac60da96345513964 | xxd -r -p | sha1sum | head -c 40
00a51f3f48415c7d4e8908980d443c29c69b60c9 # stage1hash
```

Linux Bash concatenates the scramble and stage1Hash into the string concat.

```bash
# scramble is 51402B554c5A615b223524555d5675693157417d, the first half part.
# stage1Hash is 00a51f3f48415c7d4e8908980d443c29c69b60c9, the second half part.

# calculate "20-bytes random data from server" <concatenate> SHA1( SHA1( password ) )
$ echo -n 51402B554c5A615b223524555d5675693157417d 00a51f3f48415c7d4e8908980d443c29c69b60c9 |  sed "s/ //g"
51402B554c5A615b223524555d5675693157417d00a51f3f48415c7d4e8908980d443c29c69b60c9 # concat
```

In terms of the first new scramble, Linux Bash calculates the result and compares it.

```bash
# scramble = SHA1( concat ) = SHA1( scramble <concatenate> SHA1( stage1Hash ) )

$ echo -n 51402B554c5A615b223524555d5675693157417d00a51f3f48415c7d4e8908980d443c29c69b60c9 | xxd -r -p | sha1sum | head -c 40
0ca0f764a59d1cdb10a87f0155d61aa54be1c71a # The first new scramble
```

In the case of the second new scramble, Linux Bash calculates the result and compares it.

```bash
# scramble = stage1 XOR scramble

$ stage1=0x8cb2237d0679ca88db6464eac60da96345513964 # stage1
$ scramble=0x0ca0f764a59d1cdb10a87f0155d61aa54be1c71a # The first new scramble
$ echo $(( $stage1^$scramble ))
-7792437067003134338 # insufficient precision

$ stage1=0x8cb2237d0679ca88db6464eac60da96345513964
$ scramble=0x0ca0f764a59d1cdb10a87f0155d61aa54be1c71a
$ printf "0x%X" $(( (($stage1>>40)) ^ (($scrambleFirst>>40)) ))
0xFFFFFFFFFF93DBB3 # insufficient precision

# Linux Bash divides stage1 and the first new scramble into four parts and individually makes four bitwise XOR operations.
$ printf "0x%X" $(( ((0x8cb2237d06)) ^ ((0x0ca0f764a5)) ))
$ printf "%X" $(( ((0x79ca88db64)) ^ ((0x9d1cdb10a8)) ))
$ printf "%X" $(( ((0x64eac60da9)) ^ ((0x7f0155d61a)) ))
$ printf "%X" $(( ((0x6345513964)) ^ ((0xa54be1c71a)) ))
0x8012D419A3E4D653CBCC1BEB93DBB3C60EB0FE7E # correct

# scramble is []uint8{ 80, 12,  D4, 19,  A3,  E4,  D6, 53,  CB,  CC, 1B,  EB,  93,  DB,  B3,  C6, 0E,  B0,  FE,  7E} // hexadecimal
# decimal
# scramble 为 []uint8{128, 18, 212, 25, 163, 228, 214, 83, 203, 204, 27, 235, 147, 219, 179, 198, 14, 176, 254, 126} // the same as the result in Gaea
```

The correct result, auth, is the same as Gaea's.

<img src="./assets/image-20220318183833245.png" alt="image-20220318183833245" style="zoom:70%;" /> 

### The third step

> The second step is to reply to MariaDB after receiving the initial handshake packet.

There are some details about the response packet in [the official document](https://mariadb.com/kb/en/connection/), and please see the details below.

<img src="./assets/image-20220318083633693.png" alt="image-20220318083633693" style="zoom:100%;" /> 

The actual packet demonstrates how this response works, and please see details below.

| capability                   | binary             | decimal |
| ---------------------------- | ------------------ | ------- |
| mysql.ClientProtocol41       | 0b0000001000000000 | 512     |
| mysql.ClientSecureConnection | 0b1000000000000000 | 32768   |
| mysql.ClientLongPassword     | 0b0000000000000001 | 1       |
| mysql.ClientTransactions     | 0b0010000000000000 | 8192    |
| mysql.ClientLongFlag         | 0b0000000000000100 | 4       |
|                              |                    |         |
| sum                          |                    |         |
| Gaea's capability            | 0b1010001000000101 | 41477   |

Calculate the new capability that coordinates with MariaDB's.

```
With regard to the first step, the dc object's capability is 0b10000001111111111111011111111110, which equals 2181036030 in decimal.

Obviously, it does not support "mysql.ClientLongPassword".

SpeedCrunch calculator takes both Gaea's capability and the capability to make a bitwise AND operation. The result is the mutual capability.
Gaea's capability & dc.capability = []uint32{41477} & []uint32{2181036030} = []uint32{41476}
```

<img src="./assets/image-20220319002738908.png" alt="image-20220319002738908" style="zoom:70%;" /> 

| packet                            | exmaple                                                      |
| --------------------------------- | ------------------------------------------------------------ |
| int<4> client capabilities        | Gaea reverses from the mutual capability []uint32{41476} to []uint8{4, 162, 0, 0} when sending the packet to MariaDB.<br /><img src="./assets/image-20220319113026919.png" alt="image-20220319113026919" style="zoom:50%;" /> |
| int<4> max packet size            | It occupies 4 bytes, []uint8{0, 0, 0, 0}.                    |
| int<1> client character collation | After checking [character-sets-and-collations](https://mariadb.com/kb/en/supported-character-sets-and-collations/), finding out that number 46 means "utf8mb4_bin". |
| string<19> reserved               | It occupies 19 bytes,  []uint8{<br />                                                       0, 0, 0, 0, 0,<br />                                                       0, 0, 0, 0, 0,<br />                                                       0, 0, 0, 0, 0,<br />                                                       0, 0, 0, 0,<br />                                                   } |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if not (server_capabilities & CLIENT_MYSQL)<br/>    int<4> extended client capabilities <br/>else<br/>    string<4> reserved |
| example | **CLIENT_MYSQL** means this packet belongs to Gaea, True.<br /><br />It occupies 4 bytes, []uint8{0, 0, 0, 0}, because the formula, <br />**not (server_capabilities & CLIENT_MYSQL)**<br />, is False. |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | string<NUL> username                                         |
| example | xiaomi is an account to log in MariaDB. However, it needs one byte to terminate it.<br />It will occupy 7 bytes, []uint8{120, 105, 97, 111, 109, 105, 0} |

With respect to the second new scramble, Linux Bash calculates the result and compares it.

```bash
$ echo -n xiaomi | od -td1
0000000  120  105   97  111  109  105
0000006
```

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & PLUGIN_AUTH_LENENC_CLIENT_DATA)<br/>    string<lenenc> authentication data <br/>else if (server_capabilities & CLIENT_SECURE_CONNECTION)<br/>    int<1> length of authentication response<br/>    string<fix> authentication response (length is indicated by previous field) <br/>else<br/>    string<NUL> authentication response null ended |
| example | Gaea supports CLIENT_SECURE_CONNECTION in this case.<br /><br />In the second step, the auth is []uint8{128, 18, 212, 25, 163, 228, 214, 83, 203, 204, 27, 235, 147, 219, 179, 198, 14, 176, 254, 126}.<br />However, Gaea has to send the length of the auth first.<br /><br />It will occupy 21 bytes, []uint8{20, 128, 18, 212, 25, 163, 228, 214, 83, 203, 204, 27, 235, 147, 219, 179, 198, 14, 176, 254, 126} |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & CLIENT_CONNECT_WITH_DB)<br/>    string<NUL> default database name |
| example | Gaea supports capabilities that are<br />mysql.ClientProtocol41,<br/>mysql.ClientSecureConnection,<br/>mysql.ClientTransactions,<br/>and mysql.ClientLongFlag.<br /><br />The example skips this one because Gaea doesn't support CLIENT_CONNECT_WITH_DB. |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & CLIENT_PLUGIN_AUTH)<br/>    string<NUL> authentication plugin name |
| example | Gaea supports capabilities that are<br />mysql.ClientProtocol41,<br/>mysql.ClientSecureConnection,<br/>mysql.ClientTransactions,<br/>and mysql.ClientLongFlag.<br /><br />The example skips this one because Gaea doesn't support CLIENT_PLUGIN_AUTH. |

The next table follows on from the previous one.

| item    | value                                                        |
| ------- | ------------------------------------------------------------ |
| packet  | if (server_capabilities & CLIENT_CONNECT_ATTRS)<br/>    int<lenenc> size of connection attributes<br/>    while packet has remaining data<br/>        string<lenenc> key<br/>        string<lenenc> value |
| example | Gaea supports capabilities that are<br />mysql.ClientProtocol41,<br/>mysql.ClientSecureConnection,<br/>mysql.ClientTransactions,<br/>and mysql.ClientLongFlag.<br /><br />The example skips this one because Gaea doesn't support CLIENT_CONNECT_ATTRS. |

## Testing

> There is an explanation of what I consider about in Unit Test.

### Considering about Anonymous Function

The function of the unit test whose name is "Response after Handshake," containing an anonymous function.

The variables in the anonymous function inside the test will take the address of other variables and bring them inside the function.

I consider about it again and again. It seems correct.

The problem that I am concerned about is that the anonymous function takes the wrong values.

```go
	// 交握第二步 Step2
	t.Run("Response after Handshake", func(t *testing.T) {
		var connForSengingMsgToMariadb = mysql.NewConn(mockGaea.GetConnWrite())
		dc.conn = connForSengingMsgToMariadb
		dc.conn.StartWriterBuffering()
        
		customFunc := func() {
			err := dc.writeHandshakeResponse41()
			require.Equal(t, err, nil)
			err = dc.conn.Flush()
			require.Equal(t, err, nil)
			err = mockGaea.GetConnWrite().Close()
			require.Equal(t, err, nil)
		}

		fmt.Println(mockGaea.CustomSend(customFunc).ArrivedMsg(mockMariaDB))
	})
```

## Check Result

The result is hexadecimal when I use the Linux command and tools on other websites to calculate sha1sum.

However, the IDE tool produces the result in decimal.

### Linux command

Linux Bash generates the Sha1shum.

<img src="./assets/image-20220314214316673.png" alt="image-20220314214316673" style="zoom:80%;" /> 

### Website

The tool on the website https://coding.tools/tw/sha1 calculates the sha1shum.

<img src="./assets/image-20220314215924425.png" alt="image-20220314215924425" style="zoom:80%;" /> 

### Broken Point

The IDE uses broken point to take a look at the Stage1 variable.

<img src="./assets/image-20220314220921338.png" alt="image-20220314220921338" style="zoom:100%;" /> 

### Comparison

The table below compares the two results.

One comes from taking broken point, and the other comes from tools on the website.

Thus I  am sure the result is correct.

| Position |  Binary  | Decimal | Hexadecimal |
| :------: | :------: | :-----: | :---------: |
|    0     | 10001100 |   140   |     8c      |
|    1     | 10110010 |   178   |     b2      |
|    2     | 00100011 |   35    |     23      |
|    3     | 01111101 |   125   |     7d      |
|    4     | 00000110 |    6    |     06      |
|    5     | 01111001 |   121   |     79      |
|    6     | 11001010 |   202   |     ca      |
|    7     | 10001000 |   136   |     88      |
|    8     | 11011011 |   219   |     db      |
|    9     | 01100100 |   100   |     64      |
|    10    | 01100100 |   100   |     64      |
|    11    | 11101010 |   234   |     ea      |
|    12    | 11000110 |   198   |     c6      |
|    13    | 00001101 |   13    |     0d      |
|    14    | 10101001 |   169   |     a9      |
|    15    | 01100011 |   99    |     63      |
|    16    | 01000101 |   69    |     45      |
|    17    | 01010001 |   81    |     51      |
|    18    | 00111001 |   57    |     39      |
|    19    | 01100100 |   100   |     64      |
