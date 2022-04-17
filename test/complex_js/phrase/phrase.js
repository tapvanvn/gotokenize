////simple
var a = b + c * d - (abc + f)
////inline if
var b = d ? e : f
////if follow by phrase
if( round && str.length == 1 ) str = '0'+ str
////if with else
if( round && str.length == 1 ) 
    str = '0'+ str 
else 
    str = '1' + str
////if with else if
if( round && str.length == 1 ) 
    str = '0'+ str 
else if (test)
    str = '1' + str
else 
    str = '1' + str

////lambda
var a = () => str
////lamda welform
var a = () => {
    str = "2"
}
////for with in
for (var b in c) a += 5
////for
for (let a = 0; a <= -10; i--) i *= Math.random()
////do phrase
do a+=5
    while (a==2)
////do welform
do {
    a += 5
} while (a == 2)
////while phrase
while(b == 3) 
    i -= 5
////while welform
while (b == 3) {
    i -= 5
}

////array
var a = [b, c, d]

////complex array
var a = [()=> d, function(){d += 5}, {a : b}]
////simple function
function a(b,c,d) {
    f = g
}
////class welform
class  b {
    constructor(a1,a2) {

    }
    add() {

    }
}
///class without name
var b = class {
    constructor(a1,a2) {
        a = [()=> d, function(){d += 5}, {a : b}]
    }
    add() {
        a = b + c * d - (abc + f)
    }
}
//// switch simple
switch (b) {
    case 1:
        break;
    case 2:
        break;
    default:
        break;
}
//// switch
switch (b) {
    case 1:
        a = [()=> d, function(){d += 5}, {a : b}]
        break;
    case 2:
    case 3:
        a = b + c * d - (abc + f)
        break;
    default:
        b = d ? e : f
        break;
}
////switch with block
switch (b) {
    case 1:{
            t=d
        }
    case 2:{}
        break;
    default:
        break;
}
////fullform try
try {
    a = b
} catch(ex) {
    c()
}finally{
    a = d
}
////try with finnaly
try {
    a = b
}finally{
    a = d
}
////try catch
try {
    a = b
} catch(ex) {
    c()
}

////function with return
function a(b) {
    return b + 1;
}

////label
function a(){
    d:
    if (b == c + e + f)  {
     break d;
    }
}