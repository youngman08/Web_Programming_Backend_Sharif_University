import express from "express";
import cors from "cors";
import pg from "pg";
import bodyParser from "body-parser";

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

  app.post("/cityList", async (req, res) => {
    const r = await client.query(
      "SELECT DISTINCT city FROM origin_destination"
    );
    res.send(JSON.stringify(r.rows));
  });

  app.listen(12345);
};

main();
