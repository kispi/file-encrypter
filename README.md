# file encrypter
Renaming utility implemented in Golang.

***

Build:
1. clone this repository  
2. go build  
3. 'file encrypter' binary will be generated.  

Usage:  
EX I:) file encrypter -p ./ -e (Encrypt all files to have encrypted names starting from current directory.)  
EX II:) file encrypter -p ./ -d (If filenames are encrypted via statement above, they will be decrypted by this.)  
EX III:) file encrypter [-h] (Show help)  

Caution:  
You must test this to see how this works by using non-important bunch of directories and files.  
Otherwise, you may become desperate to send me an email.  

Dependency:  
github.com/fatih/color  

Purpose:  
If you don't want anyone else can access to your file, you may want to encrypt all the files to not have extension so that they look like not-runnable file. But you definetely don't wanna do that for each and every one of those files by yourself. Then you found the right place. But if somehow they open the file, it'll be opened.

***

Contact:  
kispi@naver.com
