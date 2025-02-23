CREATE TABLE gachi_teams (
    id SERIAL PRIMARY KEY,
    team_name VARCHAR(50) NOT NULL,
    motto VARCHAR(100) NOT NULL
);

INSERT INTO gachi_teams (team_name, motto) VALUES
    ('Dungeon Lords', 'Pain is pleasure'),
    ('Muscle Brothers', 'Strength and honor'),
    ('Oil Masters', 'Slip into submission'),
    ('Thicc Squad', 'The heavier, the better'),
    ('Holy Disciples', 'Divine domination'),
    ('The Alpha Pack', 'Only the strongest survive');

CREATE TABLE gachi_personas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    strength_level INT NOT NULL,
    special_move VARCHAR(100) NOT NULL,
    favorite_drink VARCHAR(50) NOT NULL,
    battle_cry VARCHAR(100) NOT NULL,
    team_id INT REFERENCES gachi_teams(id)
);

INSERT INTO gachi_personas (name, strength_level, special_move, favorite_drink, battle_cry, team_id) VALUES
    ('Billy Herrington', 100, 'Anvil Drop', 'Protein Shake', 'Are you ready?', 2),
    ('Van Darkholme', 95, 'Whip of Submission', 'Red Wine', 'I like this kind of stuff', 1),
    ('Ricardo Milos', 90, 'Twerk of Power', 'Pina Colada', 'Let’s dance, boys!', 4),
    ('Mark Wolff', 85, 'Wolf Howl Slam', 'Whiskey on the Rocks', 'Awooo!', 2),
    ('Kazuhiko', 80, 'Smiling Slam', 'Green Tea', 'Good night, boy', 5),
    ('Dungeon Master', 99, 'Chains of Domination', 'Dark Ale', 'You have no choice', 1),
    ('Chad Thundercock', 98, 'Alpha Stomp', 'Pre-Workout Drink', 'Kneel before me!', 6),
    ('Big Boss', 97, 'Phantom Grip', 'Espresso', 'You’re pretty good', 6),
    ('Julius Belmont', 92, 'Holy Chains', 'Holy Water', 'Feel the power of discipline!', 5),
    ('Leather Baron', 91, 'Spanking of Justice', 'Black Coffee', 'You’ve been naughty!', 1),
    ('Hard Rod', 88, 'Steel Pipe Crush', 'Energy Drink', 'Let’s get HARD!', 2),
    ('Sweaty Steve', 84, 'Greased Lightning', 'Coconut Water', 'Dripping wet!', 3),
    ('Oil Overlord', 83, 'Slippery Escape', 'Olive Oil Shot', 'Too slick for you!', 3),
    ('Thicc Thunder', 81, 'Clap of Doom', 'Banana Smoothie', 'Feel the THICCNESS!', 4),
    ('Muscle Daddy', 79, 'Bear Hug Crush', 'Chocolate Milkshake', 'Come to daddy!', 2);
