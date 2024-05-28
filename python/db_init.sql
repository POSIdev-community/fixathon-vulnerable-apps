-- Create or connect to the SQLite database file
-- SQLite doesn't have a concept of CREATE DATABASE or USE statement. Just connect to the file.

-- Create the Users table
CREATE TABLE IF NOT EXISTS Users (
    userId INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

-- Create the Articles table
CREATE TABLE IF NOT EXISTS Articles (
    articleId INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    userId INTEGER,
    FOREIGN KEY (userId) REFERENCES Users(userId)
);

-- Insert users
INSERT INTO Users (username, password) VALUES
('GalacticExplorer', 'explorer123'),
('StellarTraveler', 'traveler456'),
('CosmicAdventurer', 'adventurer789');

-- Insert articles
INSERT INTO Articles (title, content, userId) VALUES
('Исследование Галактики Андромеды', 'Наше путешествие в Галактику Андромеды раскрыло захватывающие виды и таинственные явления.', 1),
('Открытие пришельцев', 'Астрономы обнаружили сигналы из далекой звездной системы, возможно, указывающие на присутствие интеллектуальных форм жизни.', 2),
('Раскрытие тайн черных дыр', 'Новые наблюдения проливают свет на загадочную природу черных дыр, вызывая сомнения в нашем понимании Вселенной.', 3),
('Путешествие к краю Вселенной', 'Отправляйтесь в эпическое путешествие к самым дальним пределам космоса, где время и пространство искривляются в невообразимых масштабах.', 1),
('Поиск экзопланет', 'Ученые выявили перспективного кандидата на обитаемую экзопланету, обращающуюся вокруг близкой звезды, разжигая надежды на обнаружение внеземной жизни.', 2);
