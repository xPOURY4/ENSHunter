# ENSHunter 🔍

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Twitter Follow](https://img.shields.io/twitter/follow/TheRealPourya?style=social)](https://twitter.com/TheRealPourya)

**ENSHunter** is a high-performance Go tool that helps you find available Ethereum Name Service (ENS) domains from a list. It scans domains in parallel, checks their availability directly on the Ethereum blockchain, and saves the available domains to a file.


## ✨ Features

- 🚀 **High Performance**: Check multiple domains concurrently using goroutines
- 📊 **Progress Visualization**: Real-time progress bar showing scan completion
- 🔄 **Auto-retry**: Automatic retry for failed requests with configurable attempts
- ⏱️ **Rate Limiting**: Configurable request rate to avoid API throttling
- 🎛️ **Configurable**: Extensive command-line options for customization
- 🧠 **Smart Processing**: Handles domains with or without `.eth` suffix
- 🔐 **Secure Configuration**: Store API keys securely in config files

## 📋 Requirements

- Go 1.18 or higher
- Infura API key ([Get one for free here](https://infura.io/register))
- Internet connection to access Ethereum network

## 🔧 Installation

### Method 1: Using `go install` (Recommended)

The easiest way to install ENSHunter is via Go's module system:

```bash
go install github.com/xPOURY4/ENSHunter@latest
```

After installation, you should be able to run `enshunter` from your terminal.

### Method 2: From Source

For the latest features or to modify the code:

```bash
# Clone the repository
git clone https://github.com/xPOURY4/ENSHunter.git

# Navigate to the project directory
cd ENSHunter

# Build the executable
go build -o enshunter

# Optional: Move to your PATH
sudo mv enshunter /usr/local/bin/  # Linux/macOS
# OR
# Add to your PATH manually on Windows
```

### Method 3: Download Binary (Coming Soon)

Pre-built binaries for various platforms will be available on the [Releases](https://github.com/xPOURY4/ENSHunter/releases) page.

## 📦 Dependencies

The application requires the following Go packages:

```bash
go get github.com/ethereum/go-ethereum
go get github.com/schollz/progressbar/v3
go get github.com/fatih/color
go get github.com/joho/godotenv
```

## 🔑 Configuration

ENSHunter supports several ways to configure your API keys and settings:

### 1. Environment File (.env)

Create a `.env` file in the directory where you run the application:

```
INFURA_KEY=your_infura_key_here
WORKERS=10
RATE_LIMIT=50
RETRIES=3
TIMEOUT=30
```

### 2. Configuration File

ENSHunter will check for a config file at `~/.enshunter/config.json`. You can save your current settings with:

```bash
enshunter -infura YOUR_INFURA_KEY -workers 10 -save-config
```

This creates a configuration file like:

```json
{
  "infura_key": "your_infura_key",
  "workers": 10,
  "rate_limit": 50,
  "retries": 3,
  "timeout": 30
}
```

### 3. Command Line Arguments (highest priority)

These override any settings from config files.

## 🚀 Usage

### Basic Usage

```bash
enshunter -input domains.txt -output available.txt
```

### Full Options

```bash
enshunter \
  -infura YOUR_INFURA_KEY \
  -input domains.txt \
  -output available.txt \
  -workers 10 \
  -rate 50 \
  -retries 3 \
  -timeout 30 \
  -verbose \
  -save-config
```

### Command Line Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `-infura` | Your Infura project ID | From config |
| `-input` | Input file containing domain names | `esn.txt` |
| `-output` | Output file for available domains | `ens_available.txt` |
| `-workers` | Number of concurrent workers | From config or `5` |
| `-rate` | Rate limit in milliseconds between requests | From config or `100` |
| `-verbose` | Enable verbose logging | `false` |
| `-retries` | Number of retries for failed requests | From config or `3` |
| `-timeout` | Request timeout in seconds | From config or `30` |
| `-save-config` | Save current settings to config file | `false` |

## 📝 Input File Format

Create a text file with one domain name per line. The domain names can be with or without the `.eth` suffix:

```
example
test123
cool-domain.eth
myfavoritename
```

You can find a sample input file at `esn.txt.sample`.

## 📊 Example Output

When the program completes, you'll see a summary like this:

```
Summary:
Total domains checked: 150
Available domains: 23
Errors: 2
Available domains saved to: available.txt
```

The output file will contain a list of available domains with the `.eth` suffix:

```
example.eth
test123.eth
myfavoritename.eth
```

## 🔍 How It Works

ENSHunter connects directly to the Ethereum blockchain using the Ethereum Name Service (ENS) contracts through Infura's API. For each domain in your input file:

1. The program checks if the domain is already registered
2. If available, it's added to the output file
3. The process continues until all domains are checked

The application uses goroutines to perform multiple checks simultaneously, significantly speeding up the process compared to sequential checking.

## 🌟 Use Cases

- Finding short or memorable ENS domains still available for registration
- Bulk checking domain availability for clients or personal use
- Monitoring domain availability for future registration
- Finding valuable domains before others discover them

## ❓ Troubleshooting

### Common Issues

- **Rate limiting**: If you encounter rate limiting errors, increase the `-rate` value
- **Timeout errors**: Try increasing the `-timeout` value or reducing the number of workers
- **Missing API key**: Ensure your Infura API key is provided via one of the configuration methods

### Debug Mode

Use the `-verbose` flag to enable detailed logging, which can help identify issues:

```bash
enshunter -input domains.txt -verbose
```

## 🛣️ Roadmap

Future features planned for ENSHunter:

- Price estimation for domain registration
- Domain expiration monitoring
- Web interface for easier use
- Support for additional name services beyond ENS
- Batch domain registration

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/xPOURY4/ENSHunter/issues) for open issues or create a new one.

To contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is [MIT](LICENSE) licensed.

## 📞 

 [@TheRealPourya](https://twitter.com/TheRealPourya) - Twitter

Project Link: [https://github.com/xPOURY4/ENSHunter](https://github.com/xPOURY4/ENSHunter)

---

<div dir="rtl">



[![نسخه Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![مجوز](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![دنبال کردن در توییتر](https://img.shields.io/twitter/follow/TheRealPourya?style=social)](https://twitter.com/TheRealPourya)

** ENS** یک ابزار قدرتمند نوشته شده با زبان Go است که به شما کمک می‌کند دامنه‌های قابل ثبت سرویس نام اتریوم (ENS) را از یک لیست پیدا کنید. این ابزار دامنه‌ها را به صورت موازی اسکن می‌کند، قابلیت ثبت آنها را مستقیماً روی بلاکچین اتریوم بررسی می‌کند و دامنه‌های قابل ثبت را در یک فایل ذخیره می‌کند.


## ✨ ویژگی‌ها

- 🚀 **کارایی بالا**: بررسی همزمان چندین دامنه با استفاده از goroutines
- 📊 **نمایش پیشرفت**: نوار پیشرفت در زمان واقعی برای نمایش وضعیت اسکن
- 🔄 **تلاش مجدد خودکار**: تلاش مجدد خودکار برای درخواست‌های ناموفق با تعداد قابل تنظیم
- ⏱️ **کنترل نرخ درخواست**: نرخ درخواست قابل تنظیم برای جلوگیری از محدودیت API
- 🎛️ **قابل تنظیم**: گزینه‌های گسترده خط فرمان برای سفارشی‌سازی
- 🧠 **پردازش هوشمند**: کار با دامنه‌های با یا بدون پسوند `.eth`
- 🔐 **پیکربندی امن**: ذخیره کلیدهای API به صورت امن در فایل‌های پیکربندی

## 📋 پیش‌نیازها

- Go 1.18 یا بالاتر
- کلید API اینفورا ([دریافت رایگان از اینجا](https://infura.io/register))
- اتصال به اینترنت برای دسترسی به شبکه اتریوم

## 🔧 نصب

### روش اول: استفاده از `go install` (توصیه شده)

ساده‌ترین روش نصب شکارچی ENS از طریق سیستم ماژول Go است:

```bash
go install github.com/xPOURY4/ENSHunter@latest
```

پس از نصب، باید بتوانید `enshunter` را از ترمینال خود اجرا کنید.

### روش دوم: از کد منبع

برای آخرین ویژگی‌ها یا تغییر کد:

```bash
# کلون کردن مخزن
git clone https://github.com/xPOURY4/ENSHunter.git

# رفتن به دایرکتوری پروژه
cd ENSHunter

# ساخت فایل اجرایی
go build -o enshunter

# اختیاری: انتقال به PATH
sudo mv enshunter /usr/local/bin/  # لینوکس/مک
# یا
# اضافه کردن به PATH به صورت دستی در ویندوز
```

### روش سوم: دانلود باینری (به زودی)

نسخه‌های از پیش ساخته شده برای پلتفرم‌های مختلف در صفحه [Releases](https://github.com/xPOURY4/ENSHunter/releases) در دسترس خواهند بود.

## 📦 وابستگی‌ها

این برنامه به پکیج‌های Go زیر نیاز دارد:

```bash
go get github.com/ethereum/go-ethereum
go get github.com/schollz/progressbar/v3
go get github.com/fatih/color
go get github.com/joho/godotenv
```

## 🔑 پیکربندی

شکارچی ENS از چندین روش برای پیکربندی کلیدهای API و تنظیمات پشتیبانی می‌کند:

### ۱. فایل محیطی (.env)

یک فایل `.env` در دایرکتوری اجرای برنامه ایجاد کنید:

```
INFURA_KEY=کلید_اینفورای_شما
WORKERS=10
RATE_LIMIT=50
RETRIES=3
TIMEOUT=30
```

### ۲. فایل پیکربندی

شکارچی ENS یک فایل پیکربندی در مسیر `~/.enshunter/config.json` را بررسی می‌کند. می‌توانید تنظیمات فعلی خود را با دستور زیر ذخیره کنید:

```bash
enshunter -infura کلید_اینفورای_شما -workers 10 -save-config
```

این دستور یک فایل پیکربندی مانند زیر ایجاد می‌کند:

```json
{
  "infura_key": "کلید_اینفورای_شما",
  "workers": 10,
  "rate_limit": 50,
  "retries": 3,
  "timeout": 30
}
```

### ۳. آرگومان‌های خط فرمان (بالاترین اولویت)

این آرگومان‌ها هر تنظیماتی از فایل‌های پیکربندی را لغو می‌کنند.

## 🚀 نحوه استفاده

### استفاده ساده

```bash
enshunter -input دامنه‌ها.txt -output دامنه‌های_موجود.txt
```

### تمام گزینه‌ها

```bash
enshunter \
  -infura کلید_اینفورای_شما \
  -input دامنه‌ها.txt \
  -output دامنه‌های_موجود.txt \
  -workers 10 \
  -rate 50 \
  -retries 3 \
  -timeout 30 \
  -verbose \
  -save-config
```

### آرگومان‌های خط فرمان

| آرگومان | توضیح | پیش‌فرض |
|----------|-------------|---------|
| `-infura` | شناسه پروژه اینفورای شما | از پیکربندی |
| `-input` | فایل ورودی حاوی نام‌های دامنه | `esn.txt` |
| `-output` | فایل خروجی برای دامنه‌های قابل ثبت | `ens_available.txt` |
| `-workers` | تعداد کارگران همزمان | از پیکربندی یا `5` |
| `-rate` | محدودیت نرخ به میلی‌ثانیه بین درخواست‌ها | از پیکربندی یا `100` |
| `-verbose` | فعال‌سازی گزارش‌دهی دقیق | `false` |
| `-retries` | تعداد تلاش‌های مجدد برای درخواست‌های ناموفق | از پیکربندی یا `3` |
| `-timeout` | مهلت درخواست به ثانیه | از پیکربندی یا `30` |
| `-save-config` | ذخیره تنظیمات فعلی در فایل پیکربندی | `false` |

## 📝 فرمت فایل ورودی

یک فایل متنی با یک نام دامنه در هر خط ایجاد کنید. نام‌های دامنه می‌توانند با یا بدون پسوند `.eth` باشند:

```
example
test123
cool-domain.eth
myfavoritename
```

شما می‌توانید یک نمونه فایل ورودی را در `esn.txt.sample` پیدا کنید.

## 📊 نمونه خروجی

وقتی برنامه تمام می‌شود، یک خلاصه مانند این خواهید دید:

```
خلاصه:
کل دامنه‌های بررسی شده: 150
دامنه‌های قابل ثبت: 23
خطاها: 2
دامنه‌های قابل ثبت در فایل زیر ذخیره شدند: available.txt
```

فایل خروجی شامل لیستی از دامنه‌های قابل ثبت با پسوند `.eth` خواهد بود:

```
example.eth
test123.eth
myfavoritename.eth
```

## 🔍 نحوه کار

شکارچی ENS با استفاده از API اینفورا مستقیماً به بلاکچین اتریوم و قراردادهای هوشمند سرویس نام اتریوم (ENS) متصل می‌شود. برای هر دامنه در فایل ورودی شما:

1. برنامه بررسی می‌کند که آیا دامنه قبلاً ثبت شده است
2. اگر قابل ثبت باشد، به فایل خروجی اضافه می‌شود
3. این فرآیند تا بررسی تمام دامنه‌ها ادامه می‌یابد

برنامه از goroutines برای انجام چندین بررسی به طور همزمان استفاده می‌کند، که سرعت فرآیند را در مقایسه با بررسی متوالی به طور قابل توجهی افزایش می‌دهد.

## 🌟 موارد کاربرد

- یافتن دامنه‌های ENS کوتاه یا جذاب که هنوز برای ثبت موجود هستند
- بررسی دسته‌ای قابلیت ثبت دامنه برای مشتریان یا استفاده شخصی
- نظارت بر قابلیت ثبت دامنه‌ها برای ثبت در آینده
- یافتن دامنه‌های ارزشمند قبل از اینکه دیگران آنها را کشف کنند

## ❓ عیب‌یابی

### مشکلات رایج

- **محدودیت نرخ**: اگر با خطاهای محدودیت نرخ مواجه شدید، مقدار `-rate` را افزایش دهید
- **خطاهای زمان انتظار**: مقدار `-timeout` را افزایش دهید یا تعداد کارگران را کاهش دهید
- **کلید API گم شده**: اطمینان حاصل کنید که کلید API اینفورای شما از طریق یکی از روش‌های پیکربندی ارائه شده است

### حالت اشکال‌زدایی

از پرچم `-verbose` برای فعال‌سازی گزارش‌دهی دقیق استفاده کنید، که می‌تواند به شناسایی مشکلات کمک کند:

```bash
enshunter -input دامنه‌ها.txt -verbose
```

## 🛣️ نقشه راه

ویژگی‌های آینده برنامه‌ریزی شده برای شکارچی ENS:

- تخمین قیمت برای ثبت دامنه
- نظارت بر انقضای دامنه
- رابط وب برای استفاده آسان‌تر
- پشتیبانی از سرویس‌های نام دیگر علاوه بر ENS
- ثبت دامنه به صورت دسته‌ای

## 🤝 مشارکت

مشارکت‌ها، گزارش مشکلات و درخواست‌های ویژگی جدید مورد استقبال قرار می‌گیرند! لطفاً [صفحه مشکلات](https://github.com/xPOURY4/ENSHunter/issues) را برای مشکلات باز بررسی کنید یا یک مورد جدید ایجاد کنید.

برای مشارکت:

1. مخزن را فورک کنید
2. یک شاخه ویژگی ایجاد کنید (`git checkout -b feature/amazing-feature`)
3. تغییرات خود را کامیت کنید (`git commit -m 'Add some amazing feature'`)
4. به شاخه پوش کنید (`git push origin feature/amazing-feature`)
5. یک درخواست Pull باز کنید

## 📄 مجوز

این پروژه تحت مجوز [MIT](LICENSE) منتشر شده است.

## 📞 تماس

[@TheRealPourya](https://twitter.com/TheRealPourya) - توییتر

لینک پروژه: [https://github.com/xPOURY4/ENSHunter](https://github.com/xPOURY4/ENSHunter)

</div>

## 📝 فرمت فایل ورودی

یک فایل متنی با یک نام دامنه در هر خط ایجاد کنید. نام‌های دامنه می‌توانند با یا بدون پسوند `.eth` باشند:

```
example
test123
cool-domain.eth
myfavoritename
```

## 📊 نمونه خروجی

```
خلاصه:
کل دامنه‌های بررسی شده: 150
دامنه‌های قابل ثبت: 23
خطاها: 2
دامنه‌های قابل ثبت در فایل زیر ذخیره شدند: available.txt
```

## 🌟 موارد کاربرد

- یافتن دامنه‌های ENS کوتاه یا جذاب که هنوز برای ثبت موجود هستند
- بررسی دسته‌ای قابلیت ثبت دامنه برای مشتریان یا استفاده شخصی
- نظارت بر قابلیت ثبت دامنه‌ها برای ثبت در آینده

## 🤝 مشارکت

مشارکت‌ها، گزارش مشکلات و درخواست‌های ویژگی جدید مورد استقبال قرار می‌گیرند! لطفاً [صفحه مشکلات](https://github.com/xPOURY4/ENSHunter/issues) را بررسی کنید.

## 📄 مجوز

این پروژه تحت مجوز [MIT](LICENSE) منتشر شده است.



[@TheRealPourya](https://twitter.com/TheRealPourya) - توییتر

</div>
