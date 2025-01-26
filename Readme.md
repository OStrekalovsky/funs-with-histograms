# Step №1. Make directory grafana_data writable
```
chmod -R 777 grafana_data
```
# Step №2. Run docker compose
## On Linux
`docker compose -f docker-compose.yaml up -d`

## On MacOS
### Enable Docker Host Network
Host networking is supported on Docker Desktop version 4.34 and later. To enable this feature:

1. Sign in to your Docker account in Docker Desktop.
2. Navigate to Settings.
3. Under the Resources tab, select Network.
4. Check the Enable host networking option.
5. Select Apply and restart.

###  Run Docker Compose
` docker compose -f docker-compose.yaml up -d`

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
