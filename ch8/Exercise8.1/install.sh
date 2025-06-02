TZ=Asia/Tokyo ./clock -port 8020 &
TZ=Europe/London ./clock -port 8030 &
TZ=US/Eastern ./clock -port 8010 &
./clockwall/clockwall NewYork=localhost:8010 London=localhost:8020 Tokyo=localhost:8030