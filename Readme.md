# Step №1. Make directory grafana_data and grafana.db writable
```
chmod 777 -R grafana_data
```
# Step №2. Run docker compose
- On linux: `docker compose -f docker-compose-linux.yaml up -d`
- On MacOS: //todo
# Step №3. Run Go program
```
cd src
go run main.go
```
# Step №4. View dashboard
1. Open Grafana at http://localhost:3000/
2. User/password: admin/admin
3. Go to "Dashboards"
4. Click on "funs-with-histograms"
