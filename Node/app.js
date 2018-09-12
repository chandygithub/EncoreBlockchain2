var express = require('express');
var app = express();
var invoke = require('./routes/invoke')
var query = require('./routes/query')

app.use('/invoke', invoke);
app.use('/query', query);

app.listen(3000);
console.log("app listening at localhost:3000");