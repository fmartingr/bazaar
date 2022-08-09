# bazaar

A service/library to extract product information from URLs.

## Configuration

| Variable name  | Default | Description                    |
| -------------- | ------- | ------------------------------ |
| `HTTP_ENABLED` | `true`  | Enable/Disable the HTTP server |
| `HTTP_PORT`    | `8080`  | Port to serve the HTTP in      |

## Servers

### HTTP

- `POST /item`

  Parameters:
  - **url**: The URL to extract information from

  Responses:
  - `200`: Information extracted
  - `400`: Shop not supported, missing parameters
  - `500`: Internal error, check logs


- `GET /liveness`

  Responses:
  - `200`: Server running


## Data model

Currently, this information is extracted from the site (if possbile):

``` js
{
    "description": "...",
    "image_url": "https://...",
    "in_stock": false,
    "name": "...",
    "price": 0.0,
    "price_text": "0,0 €",
    "release_date": "2019-03-08T00:00:00Z",
    "url": "https://..."
}
```

[pkg/models/product.go](./pkg/models/product.go)

## Supported sites

Support is handled in a _best effort_ basis. Some sites do not provided all exposed fields.

- [Amazon.es](https://amazon.es)
- [Amazon.com](https://amazon.com)
- [AkiraComics](https://www.akiracomics.com)
- [Casa del libro](https://www.casadellibro.com)
- [Heroes De Papel](https://heroesdepapel.es)
- [Steam](https://store.steampowered.com)
