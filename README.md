# bazaar

A service/library to extract product information from URLs.

## Data model

Currently, this information is extracted from the site (if possbile):

``` js
{
    "image_url": "<url>", // (string) URL to an image file
    "in_stock": false, // (bool) If the item is currently available for purchase
    "name": "<name>", // (string) The name of the product as it appears on the site
    "price": 14.21, // (optional, float) The price of the product [parsed by the library]
    "price_text": "14,21 €", // (optional, string) The price of the product as it appears on the site (with currency)
    "release_date": "2021-03-22T00:00:00Z", // (optional, string RFC3339) the release date of the item
    "url": "<url>" // (string) The URL of the item
}
```

## Supported sites

Support is handled in a _best effort_ basis. Some sites do not provided all exposed fields.

- [Amazon.es](https://amazon.es)
- [Amazon.com](https://amazon.com)
- [Akira Comics](https://www.akiracomics.com)
- [Heroes De Papel](https://heroesdepapel.es)
- [Steam](https://store.steampowered.com)

## Running

```
go run cmd/server/main.go
```
