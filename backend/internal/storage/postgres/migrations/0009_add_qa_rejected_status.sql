-- +goose Up
-- +goose StatementBegin
ALTER TABLE issues DROP CONSTRAINT IF EXISTS issues_status_check;
ALTER TABLE issues ADD CONSTRAINT issues_status_check
    CHECK (status IN ('Pending','Approved','InProgress','QA','Done','Rejected'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE issues SET status = 'InProgress' WHERE status = 'QA';
UPDATE issues SET status = 'Done' WHERE status = 'Rejected';
ALTER TABLE issues DROP CONSTRAINT IF EXISTS issues_status_check;
ALTER TABLE issues ADD CONSTRAINT issues_status_check
    CHECK (status IN ('Pending','Approved','InProgress','Done'));
-- +goose StatementEnd
