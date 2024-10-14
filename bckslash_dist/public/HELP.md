# Traefik Setup with Docker Compose

### 2. `traefik.yml`

This file contains the detailed configuration for Traefik, including entry points, Let's Encrypt certificate management, and middleware for authentication.

#### Key Sections:

- **Entry Points**: Defines two entry points—`web` (for HTTP) and `websecure` (for HTTPS).
- **Certificates**: Traefik uses Let's Encrypt to automatically manage SSL certificates with the `myresolver` resolver.
- **Middlewares**: Includes basic authentication for the Traefik dashboard. You can add more middlewares here (e.g., rate limiting).
- **Dashboard Router**: This section configures the Traefik dashboard. It’s accessible at a domain you define and is protected by basic authentication.

#### Modifications:

- To **change the domain for the Traefik dashboard**:

  In `traefik.yml`, update the following line:

  ```
  rule: "Host(`traefik.example.com`)"
  ```

  Replace `traefik.example.com` with your desired domain (e.g., `traefik.mydomain.com`).

- To **add or change dashboard authentication**:

  In `traefik.yml`, update the basic authentication middleware:

  ```
  users:
    - "user:$$apr1$$hashed_password"
  ```

  Replace `user` and `hashed_password` with your desired username and password. You can generate a hashed password using the `htpasswd` tool:

  ```
  htpasswd -nb yourusername yourpassword
  ```

### 3. Usage

1. **Access the Dashboard**:

   Visit `https://traefik.example.com` (or your chosen domain) to access the Traefik dashboard. Log in with the username and password you defined.

2. **Adding Docker Services**:

   Traefik will automatically discover services defined in Docker Compose files that are connected to the `bckslash` network and include the necessary Traefik labels.
   The labels will be generated for you when you add the service int the dashboard.

4. **Modifying Configuration**:

   If you modify `traefik.yml` or `docker-compose.yml`, restart Traefik to apply changes. Traefik can be rastarted from the menu found on the Bckslash homepage "Restart Traefik".

### 4. Additional Resources

- [Traefik Documentation](https://doc.traefik.io/traefik/)

