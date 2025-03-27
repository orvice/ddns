# ddns

Support dns:

* cloudflare
* aliyun


## Usage

This application can be configured using environment variables:

### Required Environment Variables
- `DOMAIN`: The domain name to update
- `DNS_PROVIDER`: DNS provider to use (either "cloudflare" or "aliyun")

### Provider-specific Variables

#### Cloudflare
- `CF_TOKEN`: Cloudflare API token

#### Aliyun
- `ALIYUN_ACCESS_KEY_ID`: Aliyun Access Key ID
- `ALIYUN_ACCESS_KEY_SECRET`: Aliyun Access Key Secret

### Optional Telegram Notification
- `TELEGRAM_TOKEN`: Telegram bot token
- `TELEGRAM_CHATID`: Telegram chat ID for notifications

### Example Usage

Using Cloudflare:
```env
DNS_PROVIDER=cloudflare
DOMAIN=example.com
CF_TOKEN=your_cloudflare_token
```

Using Aliyun:
```env
DNS_PROVIDER=aliyun
DOMAIN=example.com
ALIYUN_ACCESS_KEY_ID=your_access_key_id
ALIYUN_ACCESS_KEY_SECRET=your_access_key_secret
```

With Telegram notifications:
```env
TELEGRAM_TOKEN=your_bot_token
TELEGRAM_CHATID=your_chat_id
```

