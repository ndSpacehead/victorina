package storage

import "fmt"

const migrationTableName = "migrations"

var initialMigration = fmt.Sprintf(`
CREATE TABLE %[1]s (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	version INTEGER NOT NULL
);

CREATE UNIQUE INDEX migrations_version_idx ON %[1]s (version DESC);

INSERT INTO %[1]s (version) VALUES (1);`, migrationTableName)

var migrations = []string{
	fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS questions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	external_id VARCHAR(36) NOT NULL,
	question TEXT NOT NULL,
	answer TEXT NOT NULL,
	score INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS questions_external_id_idx ON questions (external_id);

INSERT INTO %s (version) VALUES (2);`, migrationTableName),

	fmt.Sprintf(`
CREATE TABLE scenarios (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	external_id VARCHAR(36) NOT NULL,
	"name" VARCHAR(256) NOT NULL,
	description TEXT NOT NULL
);

CREATE UNIQUE INDEX scenarios_external_id_idx ON scenarios (external_id);

CREATE TABLE scenarios_questions (
	scenario_id INTEGER NOT NULL,
	question_id INTEGER NOT NULL,
	score INTEGER NOT NULL,
	FOREIGN KEY (scenario_id) REFERENCES scenarios (id) ON DELETE CASCADE,
	FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE,
	UNIQUE(scenario_id, question_id)
);

INSERT INTO scenarios (external_id, "name", description)
SELECT '00000000-0000-0000-0000-000000000001', 'v1.0.0', 'Обновление имеющейся базы вопросов на новую версию';

INSERT INTO scenarios_questions (scenario_id, question_id, score)
SELECT sc.id, q.id, q.score
FROM scenarios sc, questions q
WHERE
	EXISTS (SELECT TRUE FROM questions)
	AND sc.external_id = '00000000-0000-0000-0000-000000000001';

ALTER TABLE questions DROP COLUMN score;

INSERT INTO %s (version) VALUES (3);`, migrationTableName),
}
