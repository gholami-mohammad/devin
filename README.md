<h1 dir="rtl">
 نصب و راه اندازی در محیط Docker
</h1>
<p dir="rtl">
این کار را به سه شیوه به شرح زیر میتوانید انجام دهید که راه اول ساده ترین و کم حجم ترین container را ایجاد میکند.
</p>
<h3 dir="rtl">
1. کامپایل در سیستم و راه اندازی در container
</h3>
<p dir="rtl">
در این شیوه برنامه را در محیط توسعه خود برای linux کامپایل میکنید و در یک کانتینر ساده اجرا میکنید.
در این روش لازم است که قبلا نسخه ی دلخواه
Golang
 در سیستم تان نصب باشد.
 <br>
 مراحل این کار به شرح زیر است:
</p>

**Build the Docker image**
```
docker build -t devin:light -f Dockerfile.light .
```

<p dir="rtl">
این مرحله تنها یک بار لازم است اجرا شود.
پس از ساخته شدن image
دستور زیر را اجرا کنید:
</p>

```
make docker-light
```

<h3 dir="rtl">
2. استفاده از Golang Docker image برای کامپایل و راه اندازی
</h3>

```
make docker-go
```

<h3 dir="rtl">
3. ساخت image مخصوص برنامه و اجرای آن
</h3>

**Build the Docker image**

```
docker build -t devin:go -f Dockerfile .
```
<p dir="rtl">
پس از ساخت image
برای هر بار راه اندازی از دستور زیر استفاده کنید.
</p>

```
make docker
```
