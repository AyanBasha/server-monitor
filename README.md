Server Monitor

A lightweight, self-hosted dashboard that shows whether your homelab services are up and links to each one from a single page. It's written in Go and runs in a Proxmox LXC container.

Features
Add a service with a name and address; http:// is added automatically if you leave it out
Each service's status is shown as up, unusual, down, or check config, and you can open the HTTP detail behind it
Edit or delete services from the dashboard
Optional skip TLS verification for services with self-signed certificates (e.g. Proxmox)

Requirements
Go 1.22 or newer
Network access from the host to the services you want to monitor

Usage
Click Add a service and enter a name and address, e.g. 192.168.1.10:8096.
For HTTPS services, type https:// yourself. Tick skip TLS verification if the service uses a self-signed certificate.
On the dashboard, click a service name to open it, info to edit it, or the status to see details.
