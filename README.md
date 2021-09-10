# OIDC magic

This repository contains libraries to detect the presence of ambient OIDC
credentials (e.g. GKE workload identity, Github Actions OIDC) and furnish
them for use with OIDC-consuming systems.

This library draws inspiration from `k8s.io/kubernetes/pkg/credentialprovider`,
`k8schain`, and `docker-credential-magic`.

## Usage

To use this package, import the `providers` package, and link the "plugins"
you want registered for your application.

```go
import (
	"github.com/mattmoor/oidc-magic/pkg/providers"

	// These are the registered plugins
	_ "github.com/mattmoor/oidc-magic/pkg/providers/github"
	_ "github.com/mattmoor/oidc-magic/pkg/providers/google"
)
```

You can detect whether any ambient credentials are available by
checking:

```go
	if providers.Enabled(ctx) {
```

If there are providers available, then you can get yourself an OIDC token
with a particular audience via:

```go
	tok, err := providers.Provide(ctx, "this-is-my-audience")
```

## Examples

### GKE Workload identity

To see an example with GKE workload identity, look in
`gke-workload-identity-example.yaml`.

First, create a GCP service account and allow GKE workload identity to impersonate it:

```shell
PROJECT=<INSERT YOUR PROJECT ID>

gcloud iam service-accounts create example-identity
gcloud iam service-accounts add-iam-policy-binding --role roles/iam.workloadIdentityUser --member "serviceAccount:${PROJECT}.svc.id.goog[default/example]" example-identity@${PROJECT}.iam.gserviceaccount.com
gcloud projects add-iam-policy-binding ${PROJECT} --member=serviceAccount:example-identity@${PROJECT}.iam.gserviceaccount.com --role=roles/storage.admin
```

Next, create the service account that the workload will run with:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: example-identity
  annotations:
    iam.gke.io/gcp-service-account: example-identity@mattmoor-credit.iam.gserviceaccount.com
```

Now run the job with workload identity:

```shell
# Warning: this will print an identity token to the container logs!
ko apply -f gke-workload-identity-example.yaml
```

If you examine the logs, you should see that the workload ran and printed out an identity token.

> For extra credit, comment out `serviceAccountName: example-identity`, delete the previous job,
> and run the job again.  You should see that no providers are enabled!

### Github Actions

To see examples with Github Actions, look in
`.github/workflows/github-e2e-test.yaml` at the jobs named:
* `with-permission`: This will detect the github provider and furnish a token (censored)
* `without-permission`: This will not detect the github provider.


