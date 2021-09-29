# read-and-forward
A program that reads JSON from a file, validates it and forwards it to an HTTP based API.

Requirements, thoughts and walkthrough of my thoughts included in Overview.pdf

Please note that the JSON file will be read from the root directory of where the program is run from.



Update - persistent.go
Performs the same tasks however this program waits for input on stdin stream
until stdin receives "done", where the program will finish
Main differences are since this has multiple runes it uses go routines, channels and selects!
