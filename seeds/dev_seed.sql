-- RinnoGen MVP Development Seed Data
-- Apply after schema.sql
-- Admin user (password: "password", bcrypt hashed)
INSERT INTO users (email, password, username, full_name, role, is_active)
    VALUES ('admin@example.com', '$2a$10$rHSAVC0IFQOmprSeGM4cDO26pPVZMH/ObzUp0gp2DLlOiYqGnQI6a', 'admin', 'Admin User', 'admin', TRUE)
ON CONFLICT (email)
    DO NOTHING;

-- JavaScript language (matches Piston runtime)
INSERT INTO languages (name, piston_alias, piston_version, file_extension, default_file_name)
    VALUES ('JavaScript', 'node', '18.15.0', '.js', 'solution.js')
ON CONFLICT (piston_alias, piston_version)
    DO NOTHING;

-- Published subject
INSERT INTO subjects (title, slug, color, is_published, language_id)
    VALUES ('JavaScript', 'javascript', '#f7df1e', TRUE, 1)
ON CONFLICT (slug)
    DO NOTHING;

-- Session 01
INSERT INTO subject_sessions (subject_id, title, description, order_index)
    VALUES (1, 'Session 01 - Basics', 'Learn the basics of JavaScript programming.', 1)
ON CONFLICT (subject_id, order_index)
    DO NOTHING;

-- Lesson 01
INSERT INTO lessons (subject_session_id, title, content_md, order_index)
    VALUES (1, 'Lesson 01 - Input and Output', 'Learn how to read input and write output in JavaScript.', 1)
ON CONFLICT (subject_session_id, order_index)
    DO NOTHING;

-- Problem: Sum Two Numbers
INSERT INTO problems (slug, title, difficulty, problem_md, time_limit_ms, memory_limit_mb, is_published)
    VALUES ('sum-two-numbers', 'Sum Two Numbers', 'Easy', E'## Sum Two Numbers\n\nGiven two integers, return their sum.\n\n### Input\n\nTwo space-separated integers `a` and `b`.\n\n### Output\n\nThe sum `a + b`.', 1000, 128, TRUE)
ON CONFLICT (slug)
    DO NOTHING;

-- Link lesson to problem
INSERT INTO lesson_problems (lesson_id, problem_id, order_index)
    VALUES (1, 1, 1)
ON CONFLICT (lesson_id, problem_id)
    DO NOTHING;

-- Sample test case (visible to users)
INSERT INTO test_cases (problem_id, visibility, input_data, expected_output, order_index)
    VALUES (1, 'sample', '1 2', '3', 1);

-- Hidden test case (not visible to users)
INSERT INTO test_cases (problem_id, visibility, input_data, expected_output, order_index)
    VALUES (1, 'hidden', '5 7', '12', 2);

-- GitHub account — uncomment and fill in your real installation ID
-- Find it at: https://github.com/settings/apps/rinnogen/installations
-- INSERT INTO github_accounts (user_id, installation_id, github_owner, github_owner_type)
-- VALUES (1, '133822192', 'Test-GEM-Support', 'Organization');
