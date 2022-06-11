
run:
	docker exec -it project-db psql -U subd_project -d technopark

dump:
	docker exec project-db pg_dump -U subd_project -d technopark > db_dump/tp.sql

restore:
	docker exec project-db bash -c 'cd tmp && psql -U subd_project -d postgres -c "drop database technopark with (FORCE)" && psql -U subd_project -d postgres -c "create database technopark" && psql -U subd_project -d postgres -c "grant all privileges on database technopark to subd_project" && psql -U subd_project -d technopark -1 -f tp.sql'
