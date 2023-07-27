# Jabra busy light for Linux

This service connects to the Jabra and when it became active (audio is emitted) the busy light is enable.

## Dependencies

- `systemd` to start the service
- `jabra-sdk-linux` The linux Jabra SDK [here](https://developer.jabra.com/site/global/sdk/linux/index.gsp)

```shell
sudo cp udev/99-jabra.rules /etc/udev/rules.d/99-jabra.rules
sudo udevadm control --reload
```

## Build

```shell
CGO_LDFLAGS="-Ljabra-sdk-linux_1.12.2.0/JabraLibLinux/library/ubuntu/64-bit" go build
```

## Test

```shell
LD_LIBRARY_PATH=jabra-sdk-linux_1.12.2.0/JabraLibLinux/library/ubuntu/64-bit ./jabra-busylight
```

## Deploy

```
mkdir -p ~/bin/jabra-busylight
cp jabra-busylight ~/bin/jabra-busylight/
cp jabra-sdk-linux_1.12.2.0/JabraLibLinux/library/ubuntu/64-bit/* ~/bin/jabra-busylight/
mkdir -p ~/.config/systemd/user/
cp systemd/jabra-busylight.service ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable jabra-busylight
systemctl --user start jabra-busylight
```

## Logs

`journalctl --user -u jabra-busylight -f`

```
2022/10/01 13:30:11 testing C binding: this line must be print
2022/10/01 13:30:11 attach  Jabra Link 380
2022/10/01 13:30:11 first scan
2022/10/01 13:30:13 attach  Jabra Evolve2 85
2022/10/01 13:30:26 Set busy light on Jabra Evolve2 85 to true
2022/10/01 13:30:36 Set busy light on Jabra Evolve2 85 to false
```
