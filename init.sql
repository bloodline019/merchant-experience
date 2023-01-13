CREATE TABLE products (
                          offer_id INTEGER NOT NULL,
                          name TEXT NOT NULL,
                          price NUMERIC NOT NULL,
                          quantity INTEGER NOT NULL,
                          available BOOLEAN NOT NULL,
                          seller_id INTEGER,
                          PRIMARY KEY (offer_id, seller_id)
);