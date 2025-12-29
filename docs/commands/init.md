# Command: init

The `init` command helps you get started quickly by generating a template configuration file.

## Usage

```bash
gopunch init [filename]
```

## Description

If no filename is provided, GoPunch creates `gopunch.json` in the current directory. This file contains all possible configuration options with sensible defaults.

## Generated Template

The generated file includes:
- Sample URLs to monitor.
- Default HTTP settings (timeout, method, headers).
- Concurrency and retry settings.
- A placeholder for alerting configuration (Discord/Slack webhooks).

## Example Workflow

1.  **Generate config**:
    ```bash
    gopunch init
    ```
2.  **Edit the file**: Open `gopunch.json` and add your URLs and Discord webhook URL.
3.  **Start monitoring**:
    ```bash
    gopunch watch
    ```

## Safety

The `init` command will **never** overwrite an existing file. If the target file already exists, the command will exit with an error to prevent accidental data loss.
