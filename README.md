# tabula-rasa

(still in development)

TabulaRasa is a demo web app written in Golang and React.js for managing relational databases.

The user can read, create, update and delete tables and their data through a front-end interface.

The account registration is handled by an RPC server, "account_service".

Each new registered user will have a fresh new database created for them, on which they can create custom tables 
and manage them and their data.
It supports Integer, Float, String, Boolean data types, as well as table references. 

The app uses reflection for the creation of custom structs according to the database table structure.
