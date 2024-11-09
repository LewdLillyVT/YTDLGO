# YTDLGO

YTDLGO is a YouTube downloader written in Go, designed to download both video and audio from YouTube. In addition to downloading, YTDLGO can display essential information about a YouTube video, including its length, title, and uploader.

## Features

- **Download Video**: Download videos directly from YouTube.
- **Download Audio**: Extract and download audio from YouTube videos.
- **Video Info Display**: View basic details about the video, including:
  - **Title**
  - **Uploader**
  - **Length**

## Preview

Here’s a preview of the YTDLGO GUI interface:

![YTDLGO Preview](https://cdn.hyrule.pics/a47d67564.png)


## Installation

When building YTDLGO yourself, please create a folder called `assets` in the same directory where you have the code. In the `assets` folder, add the following files:
- [yt-dlp.exe](https://github.com/yt-dlp/yt-dlp/releases/download/2024.10.22/yt-dlp.exe)
- `ffmpeg.exe` (also available from the [FFmpeg website](https://ffmpeg.org/download.html))

These files are essential dependencies for YTDLGO to function correctly.

## Usage

1. **Download Video**: Use YTDLGO to download a YouTube video in full quality.
2. **Download Audio**: Extract and download just the audio from the video.
3. **Display Video Information**: Use the tool to display key information about the YouTube link, including its title, length, and uploader.

## Example

Running YTDLGO will prompt you to input a YouTube link. The tool will then display the video's information and offer options to download either the video or audio.

## Dependencies

- [yt-dlp](https://github.com/yt-dlp/yt-dlp) (version 2024.10.22 or later)
- [FFmpeg](https://ffmpeg.org/)

