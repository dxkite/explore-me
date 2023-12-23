FOR /F %%V IN ('git rev-list HEAD --count') DO SET COUNT=%%V
FOR /F %%V IN ('git rev-parse --short HEAD') DO SET COMMIT=%%V

SET TAG=dev
SET VERSION=%TAG%
SET PROJECT=explore-me
SET BUILD_PATH=./cmd/explore-me
SET FLAGS="-s -w"

@echo "build x64"
SET GOOS=windows
SET GOARCH=amd64
SET NAME=%PROJECT%-%VERSION%-%GOOS%-%GOARCH%
go build -o %NAME%.exe -ldflags=%FLAGS% %BUILD_PATH%
7z a %NAME%.exe.zip %NAME%.exe

@echo "build x86"
SET GOOS=windows
SET GOARCH=386
SET NAME=%PROJECT%-%VERSION%-%GOOS%-%GOARCH%
go build -o %NAME%.exe -ldflags=%FLAGS% %BUILD_PATH%
7z a %NAME%.exe.zip %NAME%.exe
