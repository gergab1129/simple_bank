-- +goose Up
-- +goose StatementBegin

ALTER TABLE IF EXISTS public.transfers 
RENAME amount_id TO amount;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
