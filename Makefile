RunDashboardDev:
	cd dashboard && npm run dev

RunServiceDevAir:
	cd service && air .

RunServiceDev:
	cd service && go run main.go

ClearBadgerDB:
	rm -rf ./service/badger

