## Deployment Instructions

Create a api token with these permissions
Go to My Profile (top-right avatar) → API Tokens → Create Token. This is different from the account-level API Tokens page.
Add these permissions:
Zone | SSL and Certificates | Edit. It appears in this list; it wasn't in the account token UI.
The other permissions your config needs: Zone | DNS | Edit, Zone | Zone Settings | Edit, and Account | Access: Apps and Policies | Edit for the Zero Trust resources.
put it in .env or export

Create a Cloudflare API token scoped to the rnorwood.com zone and account with: DNS edit, zone settings edit, zone read, and Access apps and policies edit. Then run export CLOUDFLARE_API_TOKEN=... in your shell.
Enable a login method for Access in the Zero Trust dashboard (One-time PIN is the simplest).
Build and push the image to GHCR, then note its digest. I didn't check .github/workflows/build.yml for the exact tag scheme.
Create infra/envs/dev/terraform.tfvars with cloudflare_zone_id, cloudflare_account_id, container_path (ghcr.io/bobby060/site@sha256:...) and access_allowed_emails. Your email should go in there or in a TF_VAR_ variable, not in a committed file. dns_proxied stays at false.
Run terraform init and terraform plan from infra/envs/dev, then review the plan.
First apply with dns_proxied = false. Azure needs to validate the custom domain, and that fails while the record is proxied.
Check the origin lockdown. A request straight to the ACA FQDN should be rejected by the IP restriction, and the domain name should should reach the app once it's proxied.
Set dns_proxied = true and apply again. Visiting dev.rnorwood.com should now show the Access login page.