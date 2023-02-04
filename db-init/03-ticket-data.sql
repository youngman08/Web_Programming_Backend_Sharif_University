COPY aircraft_type FROM '/docker-entrypoint-initdb.d/aircraft_type.csv' DELIMITER ',' CSV HEADER;
COPY aircraft_layout FROM '/docker-entrypoint-initdb.d/aircraft_layout.csv' DELIMITER ',' CSV HEADER;
COPY aircraft FROM '/docker-entrypoint-initdb.d/aircraft.csv' DELIMITER ',' CSV HEADER;
COPY country FROM '/docker-entrypoint-initdb.d/country.csv' DELIMITER ',' CSV HEADER;
COPY city FROM '/docker-entrypoint-initdb.d/city.csv' DELIMITER ',' CSV HEADER;
COPY airport FROM '/docker-entrypoint-initdb.d/airport.csv' DELIMITER ',' CSV HEADER;
COPY flight FROM '/docker-entrypoint-initdb.d/flight.csv' DELIMITER ',' CSV HEADER;
