pdf-service/
├── cmd/
│   └── server/
│       └── main.go            # entrypoint do serviço
├── internal/
│   ├── handlers/
│   │   ├── convert.go          # convertHandler, downloadHandler
│   │   └── health.go           # healthHandler
│   ├── middleware/
│   │   └── logging.go          # loggingMiddleware e statusRecorder
│   ├── models/
│   │   └── image_metadata.go   # DownloadResponse, ImageMetadata
│   ├── storage/
│   │   ├── files.go            # saveUploadedFile, ensureDir
│   │   ├── metadata.go         # appendMetadataToFile, cleanupExpiredImages
│   │   └── cleanup.go          # startImageCleanupJob
│   └── services/
│       └── pdf_converter.go    # convertPDFToImages
├── downloads/                  # pasta para PDFs e PNGs gerados
├── go.mod
└── go.sum
