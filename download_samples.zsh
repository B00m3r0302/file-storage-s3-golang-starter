#!/usr/bin/env zsh
set -euo pipefail

# Create output directory
mkdir -p samples

# URLs to download
image_urls=(
  "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/boots-image-horizontal.png"
  "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/boots-image-vertical.png"
  "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/boots-video-horizontal.mp4"
  "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/boots-video-vertical.mp4"
  "https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/is-bootdev-for-you.pdf"
)

# Download files
for url in "${image_urls[@]}"; do
  file_name="${url:t}"   # zsh-native basename
  echo "Downloading $file_name..."
  curl -fL --progress-bar -o "samples/$file_name" "$url"
done

echo "All files downloaded to ./samples"
