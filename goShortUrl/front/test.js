var http = require("http");

var server = http.createServer(function (request, response) {
	response.end("hello")
});

server.listen("9050");