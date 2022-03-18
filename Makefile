all:
	GOOS=linux GOARCH=amd64 go build -o builds/cloudprober_external_playwright-linux-amd64/cloudprober_external_playwright
	GOOS=linux GOARCH=arm64 go build -o builds/cloudprober_external_playwright-linux-arm64/cloudprober_external_playwright
	cd builds && tar -zcf cloudprober_external_playwright-0.0.3-linux-amd64.tgz  cloudprober_external_playwright-linux-amd64
	cd builds && tar -zcf cloudprober_external_playwright-0.0.3-linux-arm64.tgz  cloudprober_external_playwright-linux-arm64
