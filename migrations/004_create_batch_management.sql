-- Tabla de Lotes
CREATE TABLE batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL REFERENCES recipes(id),
    state VARCHAR(50) NOT NULL, -- 'Planned', 'Brewing', 'Fermenting', etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de historial (Snapshot de la receta al iniciar el lote)
CREATE TABLE batch_recipe_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID REFERENCES batches(id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL,
    snapshot_data JSONB NOT NULL, -- El JSON completo de la receta en ese momento
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
