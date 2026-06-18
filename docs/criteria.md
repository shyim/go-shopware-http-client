# Criteria builder

[← Docs index](./README.md)

`Criteria` is a fluent builder that renders to the DAL search payload. Methods
return the receiver, so calls chain:

```go
c := shopware.NewCriteria().
	SetLimit(25).
	SetPage(2).
	SetTerm("t-shirt").
	SetTotalCountMode(shopware.ExactTotalCount).
	AddFilter(shopware.Equals("active", true)).
	AddFilter(shopware.EqualsAny("manufacturerId", []any{idA, idB})).
	AddFilter(shopware.Range("stock", map[string]any{"gte": 1})).
	AddSorting(shopware.Sort("name", "ASC")).
	AddFields("id", "name").
	AddIncludes(map[string][]string{"product": {"id", "name"}})
```

Filter constructors: `Equals`, `EqualsAny`, `Contains`, `Prefix`, `Suffix`,
`Range`, `Not`, `Multi`.

Sorting constructors: `Sort`, `NaturalSort`, `CountSort`.

Aggregation constructors: `TermsAggregation`, `FilterAggregation`,
`EntityAggregation`, `HistogramAggregation`, `AvgAggregation`, `SumAggregation`,
`MinAggregation`, `MaxAggregation`, `CountAggregation`, `StatsAggregation` — see
[Aggregations](./aggregations.md) for decoding the results.

## Nested associations

`AddAssociation` ensures the whole path exists; `GetAssociation` returns the
sub-criteria so you can refine it:

```go
c := shopware.NewCriteria()
c.AddAssociation("categories")
c.GetAssociation("manufacturer.media").
	AddFilter(shopware.Equals("private", false))
```

## Using a Criteria as a raw body

A `Criteria` implements `json.Marshaler`, so you can pass it as a raw request
body if you ever need to bypass the repository:

```go
resp, err := client.Post(ctx, "/search/product", c)
```

## See also

- [Entities & repositories](./entities.md) — running a criteria via `Search`.
- [Aggregations](./aggregations.md) — typed aggregation results.
