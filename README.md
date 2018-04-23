# file encrypter
File encryption utility using AES implemented in Golang.

***

Build:
1. clone this repository  
2. go build  
3. binary will be generated.(following your repository name)  

Usage:  
EX:) [BINARY] -p ./ -e (Encrypt all files to have encrypted names starting from current directory)  
EX:) [BINARY] -p ./ -d (If filenames are encrypted via statement above, they will be decrypted by this)  
EX:) [BINARY] -p ./ -e -f (encrypt all files in current path)  
EX:) [BINARY] -p ./ -d -f (decrypt all files in current path)  
EX:) [BINARY] [-h] (Show help)  

Caution:  
You must test this to see how this works by using non-important bunch of directories and files.  
Otherwise, you may become desperate to send me an email.  

Dependency:  
github.com/fatih/color  

Purpose:  
If you don't want anyone else can access to your file, you may want to encrypt all the files to not have extension so that they look like not-runnable file. But you definetely don't wanna do that for each and every one of those files by yourself. Then you found the right place. But if you don't use this program with -f option, files can be opened. Using -f option will also encrypt the content of file.

***

Contact:  
kispi@naver.com
