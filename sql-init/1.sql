use pubsub;

CREATE TABLE service_logs (
   service_name VARCHAR(100) NOT NULL,
   payload VARCHAR(2048) NOT NULL,
   severity ENUM("debug", "info", "warn", "error", "fatal") NOT NULL,
   timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE service_severity (
   service_name VARCHAR(100) NOT NULL,
   severity ENUM("debug", "info", "warn", "error", "fatal") NOT NULL,
   count INT(4) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
   PRIMARY KEY (service_name, severity)
);
