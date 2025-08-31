# networksecurity-approachtrace-project

A project to help red teams trace each other's approaches in offensive and defensive scenarios, utilizing WebSocket and Vue's Cytoscape for timely and convenient tracking.

## Overview

This tool enables red team members to track their teammates' methodologies during security operations. It addresses the challenge of coordinating distributed team members working on different attack paths, helping to avoid redundant efforts and streamline collaboration when team members need to pivot between tasks.

## Features

- **Modern Frontend**: Built with Vue 3 for a clean, intuitive user interface
- **Lightweight Backend**: Uses Go-Gin framework, deployable as a standalone executable
- **Real-time Updates**: WebSocket integration ensures instant message synchronization
- **Focus Mode**: Reduces noise from unrelated attack paths to maintain focus

## Getting Started

### Prerequisites

- Node.js (for frontend build)
- Go (for backend execution)
- Nginx (for frontend deployment)

### Installation & Deployment

1. **Frontend Deployment**
   ```bash
   # Build the frontend
   npm run build
   
   # Deploy the generated 'dist' directory on your Nginx server
   ```

2. **Backend Execution**
   ```bash
   #run db_init.sql first in database(mysql) and change the options of database connection
   # Run the backend directly
   go build main.go(windows) / go build -o main-linux.exe main.go(windows for linux)
   
   # Or compile to executable for your platform
   # (compilation commands depend on your target OS)
   ```

## Motivation

Based on real-world red team experience: when team members operate on separate attack paths, mid-mission roadblocks often require reallocation to other tasks. This platform simplifies tracking of each member's approach, reducing leadership overhead and preventing team energy depletion from redundant work.
