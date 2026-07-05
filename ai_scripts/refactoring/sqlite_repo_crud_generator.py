import os
import re

file_path = "internal/adapter/sqlite/rpg_repository.go"
with open(file_path, "r") as f:
    content = f.read()

def update_repo(receiver, struct_type, table, err_exists, err_not_found):
    global content
    
    # Save
    pattern_save = r"func \(r \*" + receiver + r"\) Save\(ctx context\.Context, [a-z] rpg\." + struct_type + r"\) error \{.*?return nil\n\}"
    repl_save = f"""func (r *{receiver}) Save(ctx context.Context, item rpg.{struct_type}) error {{
	_, err := r.db.ExecContext(ctx, `INSERT INTO {table} (name, description, grade) VALUES (?, ?, ?)`, item.Name, item.Description, item.Grade)
	if err != nil {{
		if isUniqueConstraint(err) {{
			return port.{err_exists}
		}}
		return fmt.Errorf("save {struct_type.lower()}: %w", err)
	}}
	return nil
}}"""
    content = re.sub(pattern_save, repl_save, content, flags=re.DOTALL)
    
    # FindByName
    pattern_find = r"func \(r \*" + receiver + r"\) FindByName\(ctx context\.Context, name string\) \(rpg\." + struct_type + r", error\) \{.*?return [a-z], nil\n\}"
    repl_find = f"""func (r *{receiver}) FindByName(ctx context.Context, name string) (rpg.{struct_type}, error) {{
	var item rpg.{struct_type}
	err := r.db.QueryRowContext(ctx, `SELECT name, description, grade FROM {table} WHERE name = ?`, name).Scan(&item.Name, &item.Description, &item.Grade)
	if errors.Is(err, sql.ErrNoRows) {{
		return rpg.{struct_type}{{}}, fmt.Errorf("%w: %q", port.{err_not_found}, name)
	}}
	if err != nil {{
		return rpg.{struct_type}{{}}, fmt.Errorf("find {struct_type.lower()}: %w", err)
	}}
	return item, nil
}}"""
    content = re.sub(pattern_find, repl_find, content, flags=re.DOTALL)

    # List
    pattern_list = r"func \(r \*" + receiver + r"\) List\(ctx context\.Context\) \(\[\]rpg\." + struct_type + r", error\) \{.*?\n\}\n"
    repl_list = f"""func (r *{receiver}) List(ctx context.Context) ([]rpg.{struct_type}, error) {{
	return queryList(ctx, r.db, `SELECT name, description, grade FROM {table} ORDER BY name`, func(rows *sql.Rows) (rpg.{struct_type}, error) {{
		var item rpg.{struct_type}
		err := rows.Scan(&item.Name, &item.Description, &item.Grade)
		return item, err
	}})
}}
"""
    content = re.sub(pattern_list, repl_list, content, flags=re.DOTALL)

update_repo("AbilityRepository", "Ability", "abilities", "ErrAbilityExists", "ErrAbilityNotFound")
update_repo("SkillRepository", "Skill", "skills", "ErrSkillExists", "ErrSkillNotFound")
update_repo("ItemRepository", "Item", "items", "ErrItemExists", "ErrItemNotFound")
update_repo("ProfessionRepository", "Profession", "professions", "ErrProfessionExists", "ErrProfessionNotFound")
update_repo("ClassRepository", "Class", "classes", "ErrClassExists", "ErrClassNotFound")

with open(file_path, "w") as f:
    f.write(content)
print("done")
