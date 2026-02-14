-- ============================================================================
-- Morgan Backend - Initial Seed Data
-- ============================================================================
-- This script populates the database with initial data for development
-- Updated: 2026-02-13
-- ============================================================================

-- ============================================================================
-- 1. INSTITUTIONS
-- ============================================================================
INSERT INTO auth.institutions (id, code, name, description, is_active) VALUES
    ('550e8400-e29b-41d4-a716-446655440001'::UUID, 'UP', 'Universitas Pertamina', 'Universitas Pertamina Jakarta', true),
    ('550e8400-e29b-41d4-a716-446655440002'::UUID, 'ITB', 'Institut Teknologi Bandung', 'ITB Bandung', true),
    ('550e8400-e29b-41d4-a716-446655440003'::UUID, 'UI', 'Universitas Indonesia', 'UI Jakarta', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 2. GROUPS (needed for user_roles)
-- ============================================================================
INSERT INTO iam.groups (id, institution_id, name, group_type, description, level, is_active) VALUES
    ('550e8400-e29b-41d4-a716-446655440010'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'IT Department', 'department', 'Information Technology Department', 0, true),
    ('550e8400-e29b-41d4-a716-446655440011'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'Academic Department', 'department', 'Academic Department', 0, true),
    ('550e8400-e29b-41d4-a716-446655440012'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'Operations Department', 'department', 'Operations Department', 0, true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 3. PERMISSIONS
-- ============================================================================
INSERT INTO iam.permissions (id, institution_id, code, description, module, sub_module, page, action, scope_type, is_system) VALUES
    -- User Management
    ('550e8400-e29b-41d4-a716-446655440101'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'users.iam.users.view', 'View Users', 'users', 'iam', 'users', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440102'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'users.iam.users.create', 'Create User', 'users', 'iam', 'users', 'create', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440103'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'users.iam.users.edit', 'Edit User', 'users', 'iam', 'users', 'edit', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440104'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'users.iam.users.delete', 'Delete User', 'users', 'iam', 'users', 'delete', 'api', true),

    -- Role Management
    ('550e8400-e29b-41d4-a716-446655440105'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'roles.iam.roles.view', 'View Roles', 'roles', 'iam', 'roles', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440106'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'roles.iam.roles.create', 'Create Role', 'roles', 'iam', 'roles', 'create', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440107'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'roles.iam.roles.edit', 'Edit Role', 'roles', 'iam', 'roles', 'edit', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440108'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'roles.iam.roles.delete', 'Delete Role', 'roles', 'iam', 'roles', 'delete', 'api', true),

    -- Permissions Management
    ('550e8400-e29b-41d4-a716-446655440109'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'permissions.iam.permissions.view', 'View Permissions', 'permissions', 'iam', 'permissions', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440110'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'permissions.iam.permissions.manage', 'Manage Permissions', 'permissions', 'iam', 'permissions', 'manage', 'both', true),

    -- Group Management
    ('550e8400-e29b-41d4-a716-446655440111'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'groups.iam.groups.view', 'View Groups', 'groups', 'iam', 'groups', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440112'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'groups.iam.groups.create', 'Create Group', 'groups', 'iam', 'groups', 'create', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440113'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'groups.iam.groups.edit', 'Edit Group', 'groups', 'iam', 'groups', 'edit', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440114'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'groups.iam.groups.delete', 'Delete Group', 'groups', 'iam', 'groups', 'delete', 'api', true),

    -- Shift Sessions
    ('550e8400-e29b-41d4-a716-446655440115'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'shift_sessions.schedule.shift_sessions.view', 'View Shift Sessions', 'shift_sessions', 'schedule', 'shift_sessions', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440116'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'shift_sessions.schedule.shift_sessions.create', 'Create Shift Session', 'shift_sessions', 'schedule', 'shift_sessions', 'create', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440117'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'shift_sessions.schedule.shift_sessions.edit', 'Edit Shift Session', 'shift_sessions', 'schedule', 'shift_sessions', 'edit', 'both', true),

    -- Severity Levels
    ('550e8400-e29b-41d4-a716-446655440118'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'severity_levels.hr.severity_levels.view', 'View Severity Levels', 'severity_levels', 'hr', 'severity_levels', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440119'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'severity_levels.hr.severity_levels.create', 'Create Severity Level', 'severity_levels', 'hr', 'severity_levels', 'create', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440120'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'severity_levels.hr.severity_levels.edit', 'Edit Severity Level', 'severity_levels', 'hr', 'severity_levels', 'edit', 'both', true),

    -- Shift Groups
    ('550e8400-e29b-41d4-a716-446655440121'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'shift_groups.hr.shift_groups.view', 'View Shift Groups', 'shift_groups', 'hr', 'shift_groups', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440122'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'shift_groups.hr.shift_groups.create', 'Create Shift Group', 'shift_groups', 'hr', 'shift_groups', 'create', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440123'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'shift_groups.hr.shift_groups.edit', 'Edit Shift Group', 'shift_groups', 'hr', 'shift_groups', 'edit', 'both', true),

    -- Domains
    ('550e8400-e29b-41d4-a716-446655440124'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'domains.organization.domains.view', 'View Domains', 'domains', 'organization', 'domains', 'view', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440125'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'domains.organization.domains.create', 'Create Domain', 'domains', 'organization', 'domains', 'create', 'both', true),
    ('550e8400-e29b-41d4-a716-446655440126'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'domains.organization.domains.edit', 'Edit Domain', 'domains', 'organization', 'domains', 'edit', 'both', true),

    -- System Admin
    ('550e8400-e29b-41d4-a716-446655440127'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'system.admin.admin.admin', 'System Admin', 'system', 'admin', 'admin', 'admin', 'both', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 4. ROLES
-- ============================================================================
INSERT INTO iam.roles (id, institution_id, name, display_name, description, is_system_role, is_active) VALUES
    ('550e8400-e29b-41d4-a716-446655440201'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'System Administrator', 'System Admin', 'Full system access', true, true),
    ('550e8400-e29b-41d4-a716-446655440202'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'Administrator', 'Admin', 'Administrative access', true, true),
    ('550e8400-e29b-41d4-a716-446655440203'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'Manager', 'Manager', 'Manager access', true, true),
    ('550e8400-e29b-41d4-a716-446655440204'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'Operator', 'Operator', 'Operator access', true, true),
    ('550e8400-e29b-41d4-a716-446655440205'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'User', 'User', 'Regular user access', true, true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 5. ROLE_PERMISSIONS (All permissions for System Administrator)
-- ============================================================================
INSERT INTO iam.role_permissions (id, role_id, permission_id) 
SELECT 
    gen_random_uuid()::UUID,
    '550e8400-e29b-41d4-a716-446655440201'::UUID,
    id
FROM iam.permissions
WHERE institution_id = '550e8400-e29b-41d4-a716-446655440001'::UUID
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Basic permissions for Administrator
INSERT INTO iam.role_permissions (id, role_id, permission_id) 
SELECT 
    gen_random_uuid()::UUID,
    '550e8400-e29b-41d4-a716-446655440202'::UUID,
    id
FROM iam.permissions
WHERE institution_id = '550e8400-e29b-41d4-a716-446655440001'::UUID
AND code IN (
    'users.iam.users.view',
    'users.iam.users.create',
    'users.iam.users.edit',
    'roles.iam.roles.view',
    'roles.iam.roles.edit',
    'groups.iam.groups.view',
    'groups.iam.groups.create',
    'groups.iam.groups.edit',
    'domains.organization.domains.view',
    'domains.organization.domains.create',
    'domains.organization.domains.edit',
    'shift_groups.hr.shift_groups.view',
    'shift_groups.hr.shift_groups.create',
    'shift_groups.hr.shift_groups.edit',
    'shift_sessions.schedule.shift_sessions.view',
    'shift_sessions.schedule.shift_sessions.create',
    'shift_sessions.schedule.shift_sessions.edit'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Limited permissions for Manager
INSERT INTO iam.role_permissions (id, role_id, permission_id) 
SELECT 
    gen_random_uuid()::UUID,
    '550e8400-e29b-41d4-a716-446655440203'::UUID,
    id
FROM iam.permissions
WHERE institution_id = '550e8400-e29b-41d4-a716-446655440001'::UUID
AND code IN (
    'users.iam.users.view',
    'roles.iam.roles.view',
    'groups.iam.groups.view',
    'domains.organization.domains.view',
    'shift_groups.hr.shift_groups.view',
    'shift_sessions.schedule.shift_sessions.view',
    'shift_sessions.schedule.shift_sessions.edit',
    'severity_levels.hr.severity_levels.view'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Limited permissions for Operator
INSERT INTO iam.role_permissions (id, role_id, permission_id) 
SELECT 
    gen_random_uuid()::UUID,
    '550e8400-e29b-41d4-a716-446655440204'::UUID,
    id
FROM iam.permissions
WHERE institution_id = '550e8400-e29b-41d4-a716-446655440001'::UUID
AND code IN (
    'users.iam.users.view',
    'roles.iam.roles.view',
    'groups.iam.groups.view',
    'domains.organization.domains.view',
    'shift_groups.hr.shift_groups.view',
    'shift_sessions.schedule.shift_sessions.view',
    'severity_levels.hr.severity_levels.view'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Minimal permissions for User
INSERT INTO iam.role_permissions (id, role_id, permission_id) 
SELECT 
    gen_random_uuid()::UUID,
    '550e8400-e29b-41d4-a716-446655440205'::UUID,
    id
FROM iam.permissions
WHERE institution_id = '550e8400-e29b-41d4-a716-446655440001'::UUID
AND code IN (
    'users.iam.users.view',
    'roles.iam.roles.view',
    'domains.organization.domains.view'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================================================
-- 6. USERS
-- ============================================================================
INSERT INTO auth.users (id, institution_id, external_subject, identity_provider, metadata, status) VALUES
    ('550e8400-e29b-41d4-a716-446655440301'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'admin@university.edu', 'central', 
        jsonb_build_object('full_name', 'Administrator User', 'email', 'admin@university.edu', 'department', 'IT', 'is_verified', true),
        'active'),
    ('550e8400-e29b-41d4-a716-446655440302'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'manager@university.edu', 'central',
        jsonb_build_object('full_name', 'Manager User', 'email', 'manager@university.edu', 'department', 'Academic', 'is_verified', true),
        'active'),
    ('550e8400-e29b-41d4-a716-446655440303'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'operator@university.edu', 'central',
        jsonb_build_object('full_name', 'Operator User', 'email', 'operator@university.edu', 'department', 'Operations', 'is_verified', true),
        'active'),
    ('550e8400-e29b-41d4-a716-446655440304'::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, 'user@university.edu', 'central',
        jsonb_build_object('full_name', 'Regular User', 'email', 'user@university.edu', 'department', 'General', 'is_verified', true),
        'active')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 7. USER_ROLES
-- ============================================================================
INSERT INTO iam.user_roles (id, institution_id, user_id, role_id, group_id, is_active) VALUES
    -- Admin user gets System Administrator role
    (gen_random_uuid()::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, '550e8400-e29b-41d4-a716-446655440301'::UUID, '550e8400-e29b-41d4-a716-446655440201'::UUID, '550e8400-e29b-41d4-a716-446655440010'::UUID, true),
    -- Manager gets Manager role
    (gen_random_uuid()::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, '550e8400-e29b-41d4-a716-446655440302'::UUID, '550e8400-e29b-41d4-a716-446655440203'::UUID, '550e8400-e29b-41d4-a716-446655440011'::UUID, true),
    -- Operator gets Operator role
    (gen_random_uuid()::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, '550e8400-e29b-41d4-a716-446655440303'::UUID, '550e8400-e29b-41d4-a716-446655440204'::UUID, '550e8400-e29b-41d4-a716-446655440012'::UUID, true),
    -- Regular user gets User role
    (gen_random_uuid()::UUID, '550e8400-e29b-41d4-a716-446655440001'::UUID, '550e8400-e29b-41d4-a716-446655440304'::UUID, '550e8400-e29b-41d4-a716-446655440205'::UUID, '550e8400-e29b-41d4-a716-446655440010'::UUID, true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 8. SHIFT_SESSIONS (in hr schema)
-- ============================================================================
INSERT INTO hr.shift_sessions (id, name, start, "end", status) VALUES
    ('550e8400-e29b-41d4-a716-446655440401'::UUID, 'Morning Shift', '06:00:00', '14:00:00', true),
    ('550e8400-e29b-41d4-a716-446655440402'::UUID, 'Afternoon Shift', '14:00:00', '22:00:00', true),
    ('550e8400-e29b-41d4-a716-446655440403'::UUID, 'Night Shift', '22:00:00', '06:00:00', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 9. SEVERITY_LEVELS (in master schema)
-- ============================================================================
INSERT INTO master.severity_levels (id, name, status) VALUES
    ('550e8400-e29b-41d4-a716-446655440501'::UUID, 'Critical', true),
    ('550e8400-e29b-41d4-a716-446655440502'::UUID, 'High', true),
    ('550e8400-e29b-41d4-a716-446655440503'::UUID, 'Medium', true),
    ('550e8400-e29b-41d4-a716-446655440504'::UUID, 'Low', true),
    ('550e8400-e29b-41d4-a716-446655440505'::UUID, 'Info', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 10. SHIFT_GROUPS (in hr schema)
-- ============================================================================
INSERT INTO hr.shift_groups (id, name, status) VALUES
    ('550e8400-e29b-41d4-a716-446655440601'::UUID, 'Group A', true),
    ('550e8400-e29b-41d4-a716-446655440602'::UUID, 'Group B', true),
    ('550e8400-e29b-41d4-a716-446655440603'::UUID, 'Group C', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 11. DOMAINS (in master schema)
-- ============================================================================
INSERT INTO master.domains (id, name, status) VALUES
    ('550e8400-e29b-41d4-a716-446655440701'::UUID, 'main.edu', true),
    ('550e8400-e29b-41d4-a716-446655440702'::UUID, 'api.edu', true),
    ('550e8400-e29b-41d4-a716-446655440703'::UUID, 'portal.edu', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 12. SESSIONS (test sessions for manual API testing)
-- ============================================================================
-- Session for Admin User with System Administrator role (all permissions)
WITH admin_roles AS (
    SELECT 
        json_build_object(
            'role_id', '550e8400-e29b-41d4-a716-446655440201'::TEXT,
            'role_name', 'System Administrator',
            'groups', ARRAY[]::TEXT[],
            'permissions', COALESCE(
                (
                    SELECT ARRAY_AGG(DISTINCT p.code)
                    FROM iam.role_permissions rp
                    JOIN iam.permissions p ON rp.permission_id = p.id
                    WHERE rp.role_id = '550e8400-e29b-41d4-a716-446655440201'::UUID
                ),
                ARRAY[]::TEXT[]
            )
        ) as role_json
)
INSERT INTO auth.sessions (session_id, institution_id, user_id, external_subject, roles, access_token, expires_at) 
SELECT 
    'test-admin-session-01' as session_id,
    '550e8400-e29b-41d4-a716-446655440001'::UUID as institution_id,
    '550e8400-e29b-41d4-a716-446655440301'::TEXT as user_id,
    'admin@university.edu' as external_subject,
    JSONB_BUILD_ARRAY(role_json::JSONB) as roles,
    'test-admin-token-12345678901234567890' as access_token,
    (NOW() + INTERVAL '30 days') as expires_at
FROM admin_roles
ON CONFLICT DO NOTHING;

-- Session for Manager User with Manager role
WITH manager_roles AS (
    SELECT 
        json_build_object(
            'role_id', '550e8400-e29b-41d4-a716-446655440203'::TEXT,
            'role_name', 'Manager',
            'groups', ARRAY[]::TEXT[],
            'permissions', COALESCE(
                (
                    SELECT ARRAY_AGG(DISTINCT p.code)
                    FROM iam.role_permissions rp
                    JOIN iam.permissions p ON rp.permission_id = p.id
                    WHERE rp.role_id = '550e8400-e29b-41d4-a716-446655440203'::UUID
                ),
                ARRAY[]::TEXT[]
            )
        ) as role_json
)
INSERT INTO auth.sessions (session_id, institution_id, user_id, external_subject, roles, access_token, expires_at) 
SELECT 
    'test-manager-session-01' as session_id,
    '550e8400-e29b-41d4-a716-446655440001'::UUID as institution_id,
    '550e8400-e29b-41d4-a716-446655440302'::TEXT as user_id,
    'manager@university.edu' as external_subject,
    JSONB_BUILD_ARRAY(role_json::JSONB) as roles,
    'test-manager-token-1234567890123456789011' as access_token,
    (NOW() + INTERVAL '30 days') as expires_at
FROM manager_roles
ON CONFLICT DO NOTHING;

-- ============================================================================
-- End of Seed Data
-- ============================================================================
