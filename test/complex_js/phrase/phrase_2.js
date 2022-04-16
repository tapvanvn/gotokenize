//assign contain one line lambda
var a = (b)=> d 
var a = (b) => {d()}
//assign contain inline if
var b = f == 2 ? d : e
//inline if contain assign
var b = f == 3 ? e = f : g
//one line lambda contain assign
var d = (b) => { d == e }
//nest inline if
var b = f == 3 ? e ? d = f : g : h