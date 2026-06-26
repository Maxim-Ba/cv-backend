ALTER TABLE work_history
    ALTER COLUMN logo_url TYPE JSONB USING CASE
        WHEN logo_url IS NULL THEN NULL
        ELSE to_jsonb(logo_url)
    END;
