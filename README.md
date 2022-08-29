# YT

## Dependencies:
- [youtube-dl](https://github.com/ytdl.org/youtube-dl)
- [ffmpeg](https://ffmpeg.org/)

****

## Build
```go
go build yt.go
```

## Usage
```sh
yt -u 'url'

yt -u 'url' -f mp3

yt -u 'url' -o mixtape.mp3
```

## Examples
```sh
yt -u 'https://www.youtube.com/watch?v=b9fUdJdlExU'

yt -u 'https://www.youtube.com/watch?v=b9fUdJdlExU' -f mp4

yt -u 'https://www.youtube.com/watch?v=b9fUdJdlExU' -o dark_techno.mp3

yt -u 'https://www.youtube.com/watch?v=b9fUdJdlExU' -f mp4 -o dark_techno.mp3
```
