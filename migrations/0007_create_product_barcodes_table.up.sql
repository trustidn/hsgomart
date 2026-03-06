CREATE TABLE product_barcodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    barcode VARCHAR(100) UNIQUE NOT NULL
);

CREATE INDEX idx_product_barcodes_product_id ON product_barcodes(product_id);
CREATE UNIQUE INDEX idx_product_barcodes_barcode ON product_barcodes(barcode);
