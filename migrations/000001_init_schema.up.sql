-- ============================================
-- USERS
-- ============================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- JOB CATEGORIES
-- ============================================
CREATE TABLE job_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(150) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- COMPANIES
-- ============================================
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(150) NOT NULL,
    description TEXT,
    website VARCHAR(255),
    logo_url TEXT,
    email VARCHAR(150),
    phone VARCHAR(50),
    address TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- JOBS
-- ============================================
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    company_id UUID NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    category_id INT REFERENCES job_categories (id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(250) UNIQUE NOT NULL,
    employment_type VARCHAR(50) NOT NULL, -- fulltime, parttime, etc
    work_mode VARCHAR(50) NOT NULL, -- remote, onsite, hybrid
    location VARCHAR(200),
    salary_min INT,
    salary_max INT,
    description TEXT,
    requirements TEXT,
    benefits TEXT,
    image TEXT,
    banner_url TEXT,
    contact_email VARCHAR(150),
    contact_phone VARCHAR(50),
    external_url TEXT,
    apply_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    is_approved BOOLEAN DEFAULT FALSE,
    approved_at TIMESTAMP,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- JOB APPLICATIONS
-- ============================================
CREATE TABLE job_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    job_id UUID NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    expected_salary INT,
    cv_url TEXT NOT NULL,
    portfolio_url TEXT,
    cover_letter TEXT,
    ip_address VARCHAR(200),
    user_agent VARCHAR(200),
    submitted_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- JOB VIEWS
-- ============================================
CREATE TABLE job_views (
    id SERIAL PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    ip_address VARCHAR(200),
    user_agent VARCHAR(200),
    viewed_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- JOB APPLY LOGS
-- ============================================
CREATE TABLE job_apply_logs (
    id SERIAL PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    application_id UUID REFERENCES job_applications (id) ON DELETE SET NULL,
    action_type VARCHAR(50) NOT NULL,
    ip_address VARCHAR(200),
    user_agent VARCHAR(200),
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- AUDIT LOGS
-- ============================================
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users (id) ON DELETE SET NULL,
    job_id UUID REFERENCES jobs (id) ON DELETE SET NULL,
    application_id UUID REFERENCES job_applications (id) ON DELETE SET NULL,
    action VARCHAR(150) NOT NULL,
    metadata JSONB,
    ip_address VARCHAR(200),
    user_agent VARCHAR(200),
    created_at TIMESTAMP DEFAULT NOW()
);