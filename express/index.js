import express from "express";
import cors from "cors";
import pg from "pg";
import bodyParser from "body-parser";
import fetch from "node-fetch";

const { Client } = pg;

const main = async () => {
  const app = express();
  app.use(cors());
  app.use(bodyParser.json());
  const client = new Client({
    user: "postgres",
    password: "postgres",
    database: "postgres",
    host: "db",
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
    console.log(body);
    
    const bodyBankTransaction = {receipt_id: 123, amount: 100, callback: "http://localhost:8000/TransactionCompleted.html?"};

    const response = await fetch('http://localhost:8002/transaction/', {
      method: 'post',
      body: JSON.stringify(bodyBankTransaction),
      headers: {'Content-Type': 'application/json'}
    });
    const data = await response.json();

    console.log(data);
    var url = "http://localhost:8002/payment/" + data.id;
    const out = {url: url};
    //res.status(200);
    res.send(out);
  });

  app.post("/getAllPurchases", async (req, res) => {
    const r = await client.query(
      `
      SELECT * FROM purchase WHERE corresponding_user_id == $1
      `,
      [body.corresponding_user_id]
    );

    res.send(JSON.stringify(r.rows));
  });

  
  app.post("/purchase", async (req, res) => {
    const { body } = req;

    const r = await client.query(
      `INSERT INTO purchase (corresponding_user_id, title, first_name, last_name, flight_serial, offer_price, offer_class, transaction_id, transaction_result)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, &9)`,
      [body.corresponding_user_id, body.title, body.first_name, body.last_name, body.flight_serial, body.offer_price, body.offer_class, body.transaction_id, body.transaction_result]
    );

    res.send(JSON.stringify(r.rows));
  });

  app.post("/cityList", async (req, res) => {
    const r = await client.query(
      "SELECT DISTINCT city FROM origin_destination"
    );
    res.send(JSON.stringify(r.rows));
  });

  app.listen(12345);
};


main();
