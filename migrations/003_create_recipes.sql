-- 1. La entidad Receta (Header)
-- Soporta versionado mediante (name + version)
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    style VARCHAR(100) NOT NULL,
    version INT NOT NULL DEFAULT 1,
    og NUMERIC(5, 3),
    fg NUMERIC(5, 3),
    abv NUMERIC(4, 2),
    ibu INT,
    srm INT,
    batch_size_liters FLOAT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, version)
);

-- 2. Las Etapas de la receta
CREATE TABLE recipe_stages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    sequence_order INT NOT NULL,
    duration_minutes INT,
    temperature_celsius FLOAT,
    instructions TEXT
);

-- 3. Ingredientes de la receta
CREATE TABLE recipe_ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_stage_id UUID REFERENCES recipe_stages(id) ON DELETE CASCADE,
    inventory_product_id UUID NOT NULL,
    quantity FLOAT NOT NULL,
    unit VARCHAR(20) NOT NULL,
    addition_time_minutes INT DEFAULT 0
);
