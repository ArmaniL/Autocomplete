const WebSocket = require('ws');
const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  });

const ws = new WebSocket('ws://localhost:8080/ws'); // replace with the WebSocket URL you want to test

ws.on('open', function open() {
  console.log('WebSocket connected');
  
  // send a test message to the server

  ws.send(JSON.stringify({prefix:"appl"}));
});

ws.on('message', function incoming(data) {
  console.log(`Received message: ${data}`);
});

ws.on('close', function close() {
  console.log('WebSocket disconnected');
});

ws.on('error', function error(error) {
  console.error(`WebSocket error: ${error}`);
});


rl.question('prefix? ', (answer) => {
    ws.send(JSON.stringify({prefix:answer}))
    rl.close();
  });
