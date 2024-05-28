import sqlite3

sql_script = open("db_init.sql", "r").read()

conn = sqlite3.connect('cosmic_db.sqlite')

cursor = conn.cursor()

cursor.executescript(sql_script)

conn.commit()
conn.close()

print("SQLite database cosmic_db.sqlite created and populated successfully.")
