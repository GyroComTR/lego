---
title: "Gyro DNS"
date: 2019-03-03T16:39:46+01:00
draft: false
slug: gyro
dnsprovider:
  since:    "v5.7.0"
  code:     "gyro"
  url:      "https://www.gyrohosting.com/"
---

<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->
<!-- providers/dns/gyro/gyro.toml -->
<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->

The API token is an ACME key: it can only add and remove `_acme-challenge` TXT records.


<!--more-->

- Code: `gyro`
- Since: v5.7.0


Here is an example bash command using the Gyro DNS provider:

```bash
GYRO_API_TOKEN="xxx" \
lego run --dns gyro -d '*.example.com' -d example.com
```




## Credentials

| Environment Variable Name | Description |
|-----------------------|-------------|
| `GYRO_API_TOKEN` | ACME key (Gyro DNS service page, ACME / Let's Encrypt section) |

The environment variable names can be suffixed by `_FILE` to reference a file instead of a value.
More information [here]({{% ref "dns#configuration-and-credentials" %}}).


## Additional Configuration

| Environment Variable Name | Description |
|--------------------------------|-------------|
| `GYRO_BASE_URL` | API base URL (Default: https://www.gyrohosting.com/api/v1/client/) |
| `GYRO_HTTP_TIMEOUT` | API request timeout in seconds (Default: 30) |
| `GYRO_POLLING_INTERVAL` | Time between DNS propagation check in seconds (Default: 5) |
| `GYRO_PROPAGATION_TIMEOUT` | Maximum waiting time for DNS propagation in seconds (Default: 120) |

The environment variable names can be suffixed by `_FILE` to reference a file instead of a value.
More information [here]({{% ref "dns#configuration-and-credentials" %}}).




## More information

- [API documentation](https://www.gyrohosting.com/developers)

<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->
<!-- providers/dns/gyro/gyro.toml -->
<!-- THIS DOCUMENTATION IS AUTO-GENERATED. PLEASE DO NOT EDIT. -->
