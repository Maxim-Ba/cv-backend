ALTER TABLE work_history
    ALTER COLUMN logo_url TYPE TEXT USING CASE
        WHEN logo_url IS NULL THEN NULL
        WHEN jsonb_typeof(logo_url) = 'string' THEN logo_url #>> '{}'
        ELSE logo_url::text
    END;
