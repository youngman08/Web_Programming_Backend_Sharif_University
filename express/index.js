import express from "express";
import cors from "cors";
import pg from "pg";
import bodyParser from "body-parser";
import require from "request"

const { Client } = pg;

const main = async () => {
  const app = express();
  app.use(cors());
  app.use(bodyParser.json());
  const client = new Client({
    user: "postgres",
    password: "postgres",
    host: "127.0.0.1",
    port: "5432",
  });

  await client.connect();

  app.post("/search", async (req, res) => {
    const { body } = req;
    console.log(body);
    const r = await client.query(
      `
      SELECT f.*, c1.city AS source_city, c2.city AS dest_city
      FROM flight f
      JOIN origin_destination c1 ON f.origin = c1.iata
      JOIN origin_destination c2 ON f.destination = c2.iata
      WHERE c1.city = $1 AND c2.city = $2
      `,
      [body.from, body.to]
    );
    res.send(JSON.stringify(r.rows));
  });

  app.post("/transaction", async (req, res) => {
    const { body } = req;
    var request = require('request');

    request.post(
      //First parameter API to make post request
      'http://127.0.0.1:8000/transaction/',
  
      //Second parameter DATA which has to be sent to API
      { json: {
          receipt_id: 123,
          amount: body.amount
        } 
      },
      
      //Thrid parameter Callack function  
      function (error, response, body) {
          console.log(body);
          console.log(response);
          if (!error && response.statusCode == 201) {
              console.log(body);
          }
      }
    );

    res.status(200);
    res.send('Success');
  });


  app.post("/purchase", async (req, res) => {
    const r = await client.query(
      `INSERT INTO purchase (corresponding_user_id, title, first_name, last_name, flight_serial, offer_price, offer_class, transaction_id, transaction_result)
      VALUES ($1, $2, $3, $4, $5, $6, $7)`,
      [body.from, body.to]
    );
    res.send(JSON.stringify(r.rows));
  });

  app.post("/search", async (req, res) => {
    const { body } = req;
    console.log(body);
    const r = await client.query(
      `
      SELECT f.*, c1.city AS source_city, c2.city AS dest_city
      FROM flight f
      JOIN origin_destination c1 ON f.origin = c1.iata
      JOIN origin_destination c2 ON f.destination = c2.iata
      WHERE c1.city = $1 AND c2.city = $2
      `,
      [body.from, body.to]
    );
    res.send(JSON.stringify(r.rows));
  });

  app.listen(12345);
};

main();
