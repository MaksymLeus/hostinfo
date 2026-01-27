# GH Actions Template Documentations


## Notify
### 🧩 Notification Template Enhancements
Platforms supported now:
| Platform | Method          | Description                                 |
| -------- | --------------- | ------------------------------------------- |
| Slack    | OAuth Bot Token | Send a formatted message to a Slack channel |
| Telegram | Bot Token API   | Simple bot text message                     |
| Email    | SMTP            | Send an email via SMTP (e.g., Gmail)        |
| Discord  | Webhook         | Send text embeds to a Discord channel       |
| MS Teams | Webhook         | Send card-style message to MS Teams channel |

### 🔐 Required Secrets Summary
| Platform | Secrets                                        | Description          |
| -------- | ---------------------------------------------- | -------------------- |
| Slack    | `SLACK_TOKEN`, `SLACK_CHANNEL`                 | Bot token + channel  |
| Telegram | `TELEGRAM_TOKEN`, `TELEGRAM_CHAT_ID`           | Bot token + chat     |
| Email    | `EMAIL_USERNAME`, `EMAIL_PASSWORD`, `EMAIL_TO` | SMTP credentials     |
| Discord  | `DISCORD_WEBHOOK`                              | Webhook URL          |
| MS Teams | `TEAMS_WEBHOOK`                                | Incoming webhook URL |

### 📦 Usage Example (All Platforms)
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Building..."

  notify:
    needs: [build]
    uses: ./.github/workflows/template-notifications.yml
    with:
      slack: true
      telegram: true
      email: true
      discord: true
      teams: true
      status: ${{ needs.build.result }}
      message: "Build pipeline finished"
    secrets:
      SLACK_TOKEN: ${{ secrets.SLACK_TOKEN }}
      SLACK_CHANNEL: ${{ secrets.SLACK_CHANNEL }}
      TELEGRAM_TOKEN: ${{ secrets.TELEGRAM_TOKEN }}
      TELEGRAM_CHAT_ID: ${{ secrets.TELEGRAM_CHAT_ID }}
      EMAIL_USERNAME: ${{ secrets.EMAIL_USERNAME }}
      EMAIL_PASSWORD: ${{ secrets.EMAIL_PASSWORD }}
      EMAIL_TO: ${{ secrets.EMAIL_TO }}
      DISCORD_WEBHOOK: ${{ secrets.DISCORD_WEBHOOK }}
      TEAMS_WEBHOOK: ${{ secrets.TEAMS_WEBHOOK }}
```

### 🔍 Added GitHub metadata collected at runtime:
Captured fields:
| Field              | Source                                                                                |
| ------------------ | ------------------------------------------------------------------------------------- |
| Repository         | `${{ github.repository }}`                                                            |
| Actor              | `${{ github.actor }}`                                                                 |
| Branch             | `${{ github.ref_name }}`                                                              |
| Commit SHA (short) | `git rev-parse --short HEAD`                                                          |
| Commit message     | via `git log -1 --pretty=%B`                                                          |
| Workflow run URL   | `${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}` |
| Commit URL         | `${{ github.server_url }}/${{ github.repository }}/commit/${{ github.sha }}`          |


### 🔎 Example Slack Output:
```sh
📢 CI Notification
Status: `failure`
Repository: `org/app-service`
Branch: `feature/xyz`
Commit: `abc1234` (link)
Message: "Fix broken integration tests"
Actor: `MaksymLeus`
Run: View Workflow (link)
```
This works across Slack / Discord / Teams / Telegram / Email because all use `steps.compose.outputs.msg`.

