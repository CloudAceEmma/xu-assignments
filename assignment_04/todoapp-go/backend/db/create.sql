CREATE DATABASE todoapp;
CREATE USER 'monty'@'localhost' IDENTIFIED BY 'test001';
GRANT ALL PRIVILEGES ON todoapp.* TO 'monty'@'localhost';
FLUSH PRIVILEGES;
