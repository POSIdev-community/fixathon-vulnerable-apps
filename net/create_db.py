import sqlite3
import codecs

sql_script = codecs.open("db_init.sql", "r", "utf-8").read()

conn = sqlite3.connect('cosmic_db.sqlite')

cursor = conn.cursor()

cursor.executescript(sql_script)

conn.commit()
conn.close()

print("SQLite database cosmic_db.sqlite created and populated successfully.")
