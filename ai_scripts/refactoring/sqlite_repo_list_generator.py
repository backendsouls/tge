import os
import re

file_path = "internal/adapter/sqlite/rpg_repository.go"
with open(file_path, "r") as f:
    content = f.read()

helper = """
func queryList[T any](ctx context.Context, db *sql.DB, query string, scan func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var list []T
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}
"""

def replace_basic(receiver, return_type, table):
    global content
    pattern = r"func \(r \*" + receiver + r"\) List\(ctx context\.Context\) \(\[\]rpg\." + return_type + r", error\) \{.*?\n\}\n"
    replacement = f"""func (r *{receiver}) List(ctx context.Context) ([]rpg.{return_type}, error) {{
	return queryList(ctx, r.db, `SELECT name, description FROM {table} ORDER BY name`, func(rows *sql.Rows) (rpg.{return_type}, error) {{
		var item rpg.{return_type}
		err := rows.Scan(&item.Name, &item.Description)
		return item, err
	}})
}}
"""
    content = re.sub(pattern, replacement, content, flags=re.DOTALL)

def replace_effect():
    global content
    pattern = r"func \(r \*EffectRepository\) List\(ctx context\.Context\) \(\[\]rpg\.Effect, error\) \{.*?\n\}\n"
    replacement = """func (r *EffectRepository) List(ctx context.Context) ([]rpg.Effect, error) {
	return queryList(ctx, r.db, `SELECT name, kind, description FROM effects ORDER BY name`, func(rows *sql.Rows) (rpg.Effect, error) {
		var e rpg.Effect
		var kind string
		err := rows.Scan(&e.Name, &kind, &e.Description)
		e.Kind = rpg.EffectKind(kind)
		return e, err
	})
}
"""
    content = re.sub(pattern, replacement, content, flags=re.DOTALL)

def replace_equipment():
    global content
    pattern = r"func \(r \*EquipmentRepository\) List\(ctx context\.Context\) \(\[\]rpg\.Equipment, error\) \{.*?\n\}\n"
    replacement = """func (r *EquipmentRepository) List(ctx context.Context) ([]rpg.Equipment, error) {
	return queryList(ctx, r.db, `SELECT name, slot, str, agi, intel, vit, dex, wis, cha, luk FROM equipment ORDER BY name`, func(rows *sql.Rows) (rpg.Equipment, error) {
		return scanEquipment(rows)
	})
}
"""
    content = re.sub(pattern, replacement, content, flags=re.DOTALL)

replace_basic("AbilityRepository", "Ability", "abilities")
replace_basic("SkillRepository", "Skill", "skills")
replace_basic("ItemRepository", "Item", "items")
replace_basic("ProfessionRepository", "Profession", "professions")
replace_basic("ClassRepository", "Class", "classes")
replace_effect()
replace_equipment()

content += helper

with open(file_path, "w") as f:
    f.write(content)
print("done")
