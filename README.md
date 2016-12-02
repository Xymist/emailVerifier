# Email Verifier

This package provides two functions: 

 - FindEmail, which takes a first and last name and a company name, generates possible domain names, looks up their MX records and if it gets a hit tries various permutations of the name to see if it can find an email address. This is very much in development.
 - VerifyEmail, which takes an email address, checks that there is a valid MX record for that domain, and then checks whether that email exists at that domain.

 Just include and call.
