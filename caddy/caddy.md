### Step 1: Update and Install Dependencies

1. **Update the package list**:

   ```bash
   sudo apt update
   sudo apt upgrade -y
   ```

2. **Install necessary dependencies**:

   ```bash
   sudo apt install -y apt-transport-https curl
   ```

### Step 2: Install Caddy

1. **Install the Caddy signing key**:

   ```bash
   curl -fsSL https://deb.caddyserver.com/api/gpg.key | sudo gpg --dearmor -o /usr/share/keyrings/caddy.gpg
   ```

2. **Add the Caddy repository**:

   ```bash
   echo "deb [signed-by=/usr/share/keyrings/caddy.gpg] https://deb.caddyserver.com/ubuntu/ focal main" | sudo tee /etc/apt/sources.list.d/caddy.list
   ```

3. **Install Caddy**:

   ```bash
   sudo apt update
   sudo apt install -y caddy
   ```

### Step 3: Create Your Caddy Configuration

1. **Create a directory for Caddy configuration under the `bckslash` user**:

   ```bash
   mkdir -p /home/bckslash/caddy
   ```

2. **Create the main Caddyfile**:

   ```bash
   nano /home/bckslash/caddy/Caddyfile
   ```

   Add the following content, replacing `your-email@example.com` with your email for Let's Encrypt notifications:

   ```caddyfile
   {
       email your-email@example.com  # Your email for Let's Encrypt notifications
       auto_https on                  # Enable automatic HTTPS
   }

   # Include project-specific Caddyfiles
   include /home/bckslash/bckslash/projects/*/bckslash.caddy
   ```

### Step 4: Configure the Caddy Service

1. **Create a systemd service file**:

   ```bash
   sudo nano /etc/systemd/system/caddy.service
   ```

   Add the following content:

   ```ini
   [Unit]
   Description=Caddy Web Server
   Documentation=https://caddyserver.com/docs/
   After=network.target

   [Service]
   User=bckslash
   Group=bckslash
   ExecStart=/usr/bin/caddy run --config /home/bckslash/caddy/Caddyfile --adapter caddyfile
   ExecReload=/usr/bin/caddy reload --config /home/bckslash/caddy/Caddyfile --adapter caddyfile
   Restart=always

   [Install]
   WantedBy=multi-user.target
   ```

### Step 5: Start and Enable the Caddy Service

1. **Start the Caddy service**:

   ```bash
   sudo systemctl start caddy
   ```

2. **Enable the Caddy service to start on boot**:

   ```bash
   sudo systemctl enable caddy
   ```

### Step 6: Check Caddy Status

To verify that Caddy is running without issues, use:

```bash
sudo systemctl status caddy
```

### Step 7: Reload Configuration

Whenever you make changes to the Caddyfile or any project-specific Caddyfiles, reload the Caddy configuration with:

```bash
sudo systemctl reload caddy
```

### Additional Notes

- Ensure that your domain’s DNS is pointed correctly to your server's IP address for Let's Encrypt validation.
- Make sure that ports **80** and **443** are open on your firewall to allow HTTP and HTTPS traffic.

### Summary

This guide provides a straightforward way to set up Caddy as a system-wide service while keeping all configuration files under the `bckslash` user. You can easily reload configurations as needed without stopping the service. If you have any further questions or need assistance, feel free to ask!
