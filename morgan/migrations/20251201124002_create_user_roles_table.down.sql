DROP TRIGGER IF EXISTS trg_user_roles_tenant_consistency ON iam.user_roles;
DROP FUNCTION IF EXISTS iam.ensure_user_roles_tenant_consistency;
DROP TABLE IF EXISTS iam.user_roles;
