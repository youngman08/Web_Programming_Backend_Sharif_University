// import Parse from 'parse/dist/parse.min.js';

const Parse = require('parse/node');

Parse.initialize("myAppId", "", "myMasterKey");
Parse.serverURL = 'http://localhost:1337/'

const GameScore = Parse.Object.extend("testClass");
const query = new Parse.Query(GameScore);
query.equalTo("name", "reza");
const results = await query.find();
alert("Successfully retrieved " + results.length + " scores.");
// Do something with the returned Parse.Object values
for (let i = 0; i < results.length; i++) {
  const object = results[i];
  alert(object.id + ' - ' + object.get('playerName'));
}