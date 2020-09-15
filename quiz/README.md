# Exercise 01: Quiz 

Opens a CSV file, configured by flag, with format "question,answer" and use its contents to display as an interactive game controled by timer (also configured by flag).

### CSV Example
``` csv
5+5,10
7+3,10
1+1,2
8+3,11
1+2,3
8+6,14
3+1,4
1+4,5
5+1,6
2+3,5
3+3,6
2+4,6
5+2,7
```

### Flags
```
Flag: -csv 
Description: CSV file with 'question,answer' format;
Default Value: "problems.csv"

Flag: -limit 
Description: Limit time to solve the quiz;
Default Value: 30

Flag: -h, -help 
Description: Shows all configurable flags
```