# twocli

**twocli** is a command-line two-factor authentication (2FA) application written in Go. It allows you to manage and generate Time-based One-Time Passwords (TOTPs) for your accounts securely from the terminal. All secrets are encrypted using AES-256-GCM encryption with a user-provided master password.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
    - [Add an Account](#add-an-account)
    - [List Accounts](#list-accounts)
    - [Generate TOTP Code](#generate-totp-code)
    - [Update an Account](#update-an-account)
    - [Delete an Account](#delete-an-account)
- [Security Considerations](#security-considerations)
- [Examples](#examples)
- [Testing](#testing)
- [License](#license)
- [Contributing](#contributing)
- [Contact](#contact)

---

## Features

- **Add Accounts**: Securely store multiple 2FA accounts with names and secrets.
- **List Accounts**: View all saved account names.
- **Generate TOTP Codes**: Generate TOTP codes for your accounts with:
  - Real-time countdown timer
  - Color-coded progress bar
  - Automatic code refresh
  - Visual remaining time indicator
- **Update Accounts**: Update the secret key of an existing account.
- **Delete Accounts**: Remove accounts you no longer need.
- **Secure Encryption**: All secrets are encrypted using AES-256-GCM with a master password.
- **Cross-Platform**: Works on Unix-like systems and Windows.

---

## Prerequisites

- **Go Programming Language**: You need Go installed (version 1.17 or later).
    - Download and install Go from [golang.org](https://golang.org/dl/).

---

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/bykclk/twocli.git
   cd twocli
   ```

2. **Initialize Go Modules**

   ```bash
   go mod tidy
   ```

3. **Build the Application**

   ```bash
   go build -o twocli ./cmd/twocli
   ```

---

## Usage

**General Syntax:**

```bash
./twocli [command] [options]
```

### Available Commands

- `add`     - Add a new account
- `list`    - List all saved accounts
- `code`    - Generate TOTP code for an account
- `update`  - Update the secret key of an existing account
- `delete`  - Delete an existing account

### Global Options

- `-h`, `--help`  - Show help information

---

### Add an Account

Add a new account with a name and secret key.

**Syntax:**

```bash
./twocli add -name ACCOUNT_NAME -secret SECRET_KEY
```

**Options:**

- `-name`   - The name of the account
- `-secret` - The base32-encoded secret key for the account

**Example:**

```bash
./twocli add -name GitHub -secret JBSWY3DPEHPK3PXP
```

---

### List Accounts

List all saved accounts.

**Syntax:**

```bash
./twocli list
```

**Example:**

```bash
./twocli list
```

---

### Generate TOTP Code

Generate a TOTP code for a specified account.

**Syntax:**

```bash
./twocli code -name ACCOUNT_NAME [-auto]
```

**Options:**

- `-name` - The name of the account
- `-auto` - Automatically generate new codes when the current one expires

**Example:**

```bash
# Generate a single code
./twocli code -name GitHub

# Generate codes automatically
./twocli code -name GitHub -auto
```

**Features:**
- Color-coded progress bar that changes based on remaining time:
  - Green: > 15 seconds
  - Yellow: 6-15 seconds
  - Red: ≤ 5 seconds
- Visual countdown timer
- Automatic code refresh (with -auto flag)
- Clean and modern UI

---

### Update an Account

Update the secret key of an existing account.

**Syntax:**

```bash
./twocli update -name ACCOUNT_NAME -secret NEW_SECRET_KEY
```

**Options:**

- `-name`   - The name of the account
- `-secret` - The new base32-encoded secret key for the account

**Example:**

```bash
./twocli update -name GitHub -secret NEWSECRETKEY
```

---

### Delete an Account

Delete an existing account.

**Syntax:**

```bash
./twocli delete -name ACCOUNT_NAME
```

**Options:**

- `-name` - The name of the account

**Example:**

```bash
./twocli delete -name GitHub
```

---

## Security Considerations

- **Master Password**: A master password is required to encrypt and decrypt your account secrets. Choose a strong, memorable password.
- **Password Input**: When prompted for your master password, input is hidden for security.
- **Encryption**: Secrets are encrypted using AES-256-GCM with a key derived from your master password using PBKDF2 with SHA-256 and 100,000 iterations.
- **Data Storage**: Account data is stored in the `data/accounts.db` file with restrictive permissions (`0600`).
- **Failed Attempts**: After 3 incorrect master password attempts, the application will exit to prevent brute-force attacks.

---

## Examples

### Adding and Using an Account

1. **Add an Account**

   ```bash
   ./twocli add -name Gmail -secret JBSWY3DPEHPK3PXP
   ```

    - Enter your master password when prompted.

2. **List Accounts**

   ```bash
   ./twocli list
   ```

    - Enter your master password.

3. **Generate TOTP Code**

   ```bash
   ./twocli code -name Gmail
   ```

    - Enter your master password.
    - The TOTP code will be displayed.

### Updating an Account’s Secret

```bash
./twocli update -name Gmail -secret NEWSECRETKEY
```

- Enter your master password.

### Deleting an Account

```bash
./twocli delete -name Gmail
```

- Enter your master password.
- Confirm the deletion by typing `yes` when prompted.

---

## Testing

The project includes unit tests for critical components.

### Running Tests

```bash
go test ./...
```

---

## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

---

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the Repository**

   Click the “Fork” button at the top right of this page.

2. **Clone Your Fork**

   ```bash
   git clone https://github.com/bykclk/twocli.git
   ```

3. **Create a Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Commit Your Changes**

   ```bash
   git commit -am 'Add new feature'
   ```

5. **Push to the Branch**

   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**

   Navigate to the original repository and click “New Pull Request”.

---

## Contact

For questions or support, please open an issue on the [GitHub repository](https://github.com/bykclk/twocli/issues).

---

**Disclaimer**: Use this tool responsibly. The author is not responsible for any loss of data or security breaches resulting from the use of this application.
