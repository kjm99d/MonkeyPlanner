# Enabling Homebrew formula auto-push

Once this secret is set, every `vX.Y.Z` tag pushed to `main` will:

1. Build release archives via goreleaser
2. Regenerate `Formula/monkey-planner.rb` in [`kjm99d/homebrew-tap`](https://github.com/kjm99d/homebrew-tap)
3. Commit + push the formula with message `chore(formula): update monkey-planner to vX.Y.Z`

Users can then install with:

```bash
brew tap kjm99d/tap
brew install monkey-planner
```

## One-time setup (maintainer only)

### 1. Generate a fine-grained PAT

1. Open <https://github.com/settings/personal-access-tokens/new>
2. Set:
   - **Token name**: `MonkeyPlanner homebrew-tap writer`
   - **Expiration**: 1 year (or no expiry)
   - **Repository access**: "Only select repositories" → `kjm99d/homebrew-tap`
   - **Repository permissions**:
     - `Contents`: **Read and write**
     - `Metadata`: Read-only (required)
3. Click **Generate token** and copy the value (starts with `github_pat_`).

### 2. Add the secret to MonkeyPlanner

1. Open <https://github.com/kjm99d/MonkeyPlanner/settings/secrets/actions>
2. Click **New repository secret**
3. **Name**: `HOMEBREW_TAP_TOKEN`
4. **Secret**: paste the PAT value
5. Click **Add secret**

### 3. Verify

Tag a patch release to test:

```bash
git tag v1.4.1
git push origin v1.4.1
```

Watch <https://github.com/kjm99d/MonkeyPlanner/actions/workflows/release.yml>.
After it succeeds, check <https://github.com/kjm99d/homebrew-tap/tree/master/Formula>
for `monkey-planner.rb`.

## Troubleshooting

**Release succeeds but the formula is not pushed**
Check the goreleaser log for `skip_upload=true` — it means `HOMEBREW_TAP_TOKEN`
was not in the env. Verify the secret name matches exactly.

**`403 Permission denied` on the tap repo**
The PAT's **Repository access** must explicitly include `kjm99d/homebrew-tap`.
"All repositories" also works but is broader than necessary.

**The formula is pushed to `main` but the tap default branch is `master`**
`.goreleaser.yml` explicitly sets `branch: master`. If the tap default branch
changes, update that field too.
