# Output Formats

GoPunch supports multiple output formats to suit different needs, from human-readable tables to machine-parseable JSON or CSV.

## Table (Default)
Designed for interactive use. It uses colors to highlight success (green), failure (red), and warnings (yellow).

**Example:**
```text
STATUS  TARGET                   CODE/INFO  TIME   NOTE
✓       https://google.com       200        145ms  1.2KB
✗       https://github.com       500        89ms   Internal Error
!       tcp://localhost:5432     Open       5ms    
```

## JSON
Ideal for integration with other tools or scripts (e.g., using `jq`).

**Example:**
```json
[
  {
    "url": "https://google.com",
    "info": "200",
    "duration_ms": 145,
    "size": 1234,
    "success": true,
    "error": null
  }
]
```

## CSV
Useful for data analysis in Excel or other spreadsheet software.

**Header:** `url,info,duration_ms,size,success,error`

**Example:**
```csv
https://google.com,200,145,1234,true,
https://github.com,500,89,0,false,Internal Error
```

## Minimal
The most concise format, printing only the status and the URL. Perfect for piping into other CLI tools.

**Example:**
```text
200 https://google.com
ERR https://github.com
Open tcp://localhost:5432
```

## Quiet Mode (`-q` / `--quiet`)
In this mode, GoPunch produces **no output**. It only returns an exit code:
- `0`: All checks passed.
- `1`: One or more checks failed.

This is the preferred way to use GoPunch in CI/CD pipelines or shell scripts where you only care about the result.
