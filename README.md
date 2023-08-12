# herpstat_spyderweb_exporter
[![Current Version](https://img.shields.io/badge/Version-0.0.1a-brightgreen)](https://github.com/jjack/herpstat_spyderweb_exporter/releases/latest)

A [Prometheus Exporter](https://prometheus.io/) for the Herpstat Sypderweb thermostat/humidity controller by [Spyder Robotics](https://spyderrobotics.com/).

It runs by default on port `10010` and the only flag required is for your `herpstat.address`, which can be its IP address or domain.

## Installation and Usage

### Herpstat SpyderWeb Configuration

Enable the `/RAWSTATUS` page on your device:
1. Log in to your Herpstat SpyderWeb and access the admin panel at `http://herpstat.address/handleAdminControl`
2. Check the box for "Show Advanced Status Integration Options"
3. Check the box for "Enable /RAWSTATUS Page"
4. Click "SUBMIT"

> **Note**
> The Herpstat SpyderWeb docs suggest only hitting `http://herpstat.address/RAWSTATUS` once every 10 seconds. If you go over this, `herpstat_spyderweb_exporter` will return cached data to prevent blocking any necessary actions.

### Prometheus Config (prometheus.yml)

```
scrape_configs:
  - job_name: herpstat_spyderweb_exporter
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:10010']
```

### Docker
```
docker run -d \
    -p 10010:10010 
    ghcr.io/jjack/herpstat_spyderweb_exporter:latest
    --herpstat.address 1.2.3.4
```

### Docker Compose
```
---
version: '3.8'

services:
  herpstat_spyderweb_exporter:
    image: ghcr.io/jjack/herpstat_spyderweb_exporter:latest
    container_name: herpstat_spyderweb_exporter
    command:
      - '--herpstat.address=1.2.3.4'
    port:
      - 10010:10010
    restart: unless-stopped
```

### Available Options:
|  CLI Flag | Docker Env Var | Description  |  Default |  Required |
|---|---|---|---|---|
| --herpstat.address | HERPSTAT_SPYDERWEB_EXPORTER_ADDRESS | Address of your Herpstat Spyderweb |  | YES |
| --web.port | HERPSTAT_EXPORTER_WEB_PORT | The port on which herpstat_spyderweb_exporter listens | 10010 |  |
| --web.telemetry-path | HERPSTAT_EXPORTER_TELEMETRY_PATH | The path on whcih herpstat_spyderweb_exporter exposes metrics. | /metrics |  |
| --web.disable-exporter-metrics | HERPSTAT_EXPORTER_DISABLE_EXPORTER_METRICS |Exclude metrics about the exporter itself (promhttp_*, process_*, go_*). | no |  |
| --help | n/a | Show context-sensitive help | no | |
| --debug | HERPSTAT_EXPORTER_DEBUG | Enable debugging log output. (It's noisy!) | no | |


## Metrics Collected

| Name | Description | Labels | Misc Info |
|---|---|---|---|
| herpstat_system_info | Metadata information about the Herpstat Spyderweb system itself | name, firmware, ip, mac, # of outputs | |
| herpstat_system_safetyrelay  | Safety Relays enabled | name, relay | Has a value of 0 until a relay is triggered. Then it becomes 1 and "relay" becomes the relay message.|
| herpstat_system_temp  | Current internal temperature | name | |
| herpstat_system_reset_total  | Number of times the Herpstat Spyderweb has lost power and/or been reset | name |  Value comes from the Herpstat, not `herpstat_spyderweb_exporter` |
| herpstat_output_info | Metadata information about Herpstat Spyderweb Outputs itself | id, name, system, mode | |
| herpstat_output_power | This output's current power output % | id, system | |
| herpstat_output_power | This output's current power output limit % | id, system | |
| herpstat_output_probe_temperature | This output's probe's current temperature reading | id, system | |
| herpstat_output_probe_humidity | This output's probe's current humidity reading | id, system | |
| herpstat_output_ramping  | Is this output ramping? | id, system | |
| herpstat_output_alarm_high  | This output's high alarm value | id, system | |
| herpstat_output_alarm_low  | This output's low alarm value | id, system | |
| herpstat_output_alarm_enabled  | Does this output have a high/low alarm? | id, system | |
| herpstat_output_error | This output's error state/number | id, system, error | |
