dp-published-collection-api
================

This API is used to gather information about collections which have been published
using the new process.

### API

#### /publishedcollection
Returns an array of objects containing the most recent collection which have been 
published.

##### JSON structure
```
[{
	"id": Collection id,
	"name": Name of the collection,
	"publishDate": Published date in ISO 8601 format,
	"publishStartDate": Published start date in ISO 8601 format,
	"publishEndDate": Published end date in ISO 8601 format,
}, ...]
```
#### /publishedcollection/<collection id>
Returns an objects containing more detailed information about the published 
collection.

##### JSON structure
```
{
	"id": Collection id,
	"name": Name of the collection,
	"publishDate": Published date in ISO 8601 format,
	"publishStartDate": Published start date in ISO 8601 format,
	"publishEndDate": Published end date in ISO 8601 format,
	"publishResults":[
		{
			"startTime": The start time for the item when it was published
			"endTime": The time once this item was published to the website
			"uri": The uri of the item
			"size": The file size in bytes
		}, ...]
}
```
### Configuration

| Environment variable | Default | Description
| -------------------- | ---------- | -----------
| PORT                 | 9090       | The host and port to bind to
| DB_ACCESS            | localhost  | Parameters used to connect to Postgres

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
