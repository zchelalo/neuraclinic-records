UPDATE attachments
SET upload_status = 'FILE_STATUS_AVAILABLE',
    updated_at = now()
WHERE deleted_at IS NULL;
