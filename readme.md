# JSON Configurator (jc)

## Description
Pass in two JSON files, a template and a file to compare to it.  Program will prompt user for any fields found in the template and not in the compared file, allowing the user to easily update their config for potential new flags or values.  

The idea is that you can have an example.json file that contains default config values for a software project.  If a new version of the software requires new config values, this program can assist by making the changes more obvious to the user and streamlining the file updates.

## Usage
Without building binary:
```
go run main.go -template=[pathtotemplate] -file=[pathtocomparefile]
```

Outputs the differences found to the console


## Next additions...
Program will also look at any potential changes to values that might not be reflected in the compared file (i.e. array value changes), and prompt the user to use the template values, input their own, or keep the current ones.  