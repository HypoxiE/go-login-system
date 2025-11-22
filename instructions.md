**1. `make install clean`**

**2. Edit files in `./cmd/go-login/text_templates`**

**3. To `systemctl edit getty@tty1`** put
```
[Service]
ExecStart=
ExecStart=-/sbin/agetty --skip-login --noissue --noclear --login-program /bin/gologin %I $TERM
Type=idle
NoNewPrivileges=no
```