# Safe Processing
- Avoid writing files to FS during extraction
- Use inbuilt db

# Global Store
- Single store for all evidence
- Shared space
- Lossless compression
- Map data bits to preserve space
- Homomorphic encryption?

# OSINT
- Collect APIs to download data from OSNs
- Collect other relevant OSINT APIs
- Agent for safe data download
- Agent for safe data transfer

# Fast processing
- RAM Disk for data acquisition and analysis speed up
- In-memory analysis and reporting
- In-memory file system vs RAM Disk?

# Criminal behaviour detection
- Weapon detection / classification

# Audio Analytics (Low Priority)
- Audio match
- Audio to text

# Video Analytics
- Include features of image analytics
    - PoI Identification
    - Weapon Detection / Classification

# Text Analytics (Low Priority)
- Semantic Analysis
- Sentiment analysis

# Forgery detection
- Watermark
- Logo
- Pattern detection
- Video Foregery

# Map Artifacts
- CKC
- MITRE ATT&CK TTPs

----------------------------------
# Proposed DB to use
- dgraph
- influxdb
----------------------------------

# Memory Forensics
- Use other memory forensics tools and see what capabilities can be enhanced

# Crypto Forensics
- Wallet search
- Key search
- Blocckahin/DB search

# Improved Disk Imaging
- Live / Dead Imaging
- CRC Checks for header and data (or something better)
- Multiple Files
- Appending md5 or sha256 hashes of chunk files at the end

# Improved Disk Image Reading
- Mount in Virtual File System
- Read without mount

----------------------------------

# User Action Traceback
- History of fired commands in terminal

----------------------------------
# Misc Features
- Maintain history of fired commands in CLI mode