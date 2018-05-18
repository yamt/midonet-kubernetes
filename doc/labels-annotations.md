# Well-Known Labels, Annotations, etc

## Labels

| Label                 | Resource    | Description                         |
|:----------------------|:------------|:------------------------------------|
| midonet.org/owner-uid | Translation | UID of the k8s resource to which this Translation belongs |

## Annotations

| Annotation            | Resource    | Description                         |
|:----------------------|:------------|:------------------------------------|
| midonet.org/host-id   | Node        | The corresponding MidoNet Host ID   |

## Finalizers

| Finalizer             | Resource    | Description                         |
|:----------------------|:------------|:------------------------------------|
| midonet.org/deleter   | Translation | Postpone deletion for MidoNet API sync |