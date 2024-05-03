echo "Warning, this will drop all tables from the database before migrating, you WILL lose all your data."
read -p "Continue to migrate database? [Y/n] " choice
[[ $choice == "y" || $choice == "Y" ]] && docker exec ecss-locker-db './home/push-schema.sh'
