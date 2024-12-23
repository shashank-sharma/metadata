# Metadata

Fyne application - For metadata collection

### Mac

```
fyne install -icon icon.png
```

### Windows

```
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_ENABLED=1 fyne package -os windows -icon icon.png
```

### Android testing

```
adb devices
fyne package -os android -appID mytest.domain.metadata
adb install <filename>.apk
adb logs
```
