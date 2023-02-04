import express from "express";
import pg from "pg";

const { Client } = pg;

const main = async () => {
  const app = express();
  const client = new Client({
    user: 'postgres',
    password: 'postgres',
    host: '127.0.0.1',
    port: '5432',
  });

  await client.connect();

  app.get("/search", async (req, res) => {
    const r = await client.query('SELECT * FROM flight');
    res.send(JSON.stringify(r));
  });

  app.listen(12345);
};

main();
