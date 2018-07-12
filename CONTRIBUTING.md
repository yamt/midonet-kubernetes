## Contribution

### Submitting patches

We use [GerritHub][gerrithub] to submit patches.

[gerrithub]: https://review.gerrithub.io/q/project:midonet%252Fmidonet-kubernetes

We don't use GitHub pull requests.

### Reviewing patches

Everyone is enouraged to review [patches for this repository][patches to review].

[patches to review]: https://review.gerrithub.io/q/project:midonet%252Fmidonet-kubernetes+status:open

If you want to be notified of patches, you can add this repository to
["Watched Projects"][watched projects] in your GerritHub settings.

[watched projects]: https://review.gerrithub.io/#/settings/projects

We have a voting CI named "Midokura Bot".
Unfortunately, its test logs are not publicly available.
If it voted -1 on your patch, please ask one of Midokura employees
to investigate the log.

### Merging patches

Unless it's urgent, a patch should be reviewed by at least one person
other than the submitter of the patch before being merged.

Right now, members of [GerritHub midonet group][midonet group] have the permission to merge patches.
If you are interested in being a member, please reach out the existing members.

[midonet group]: https://review.gerrithub.io/#/admin/groups/80,members

### Issue tracker

Bugs and Tasks are tracked in [MidoNet jira][jira].
We might consider alternatives if the traffic goes up.

[jira]: https://midonet.atlassian.net/

We don't use GitHub issues.

## Release process

Right now, our releases are tags on master branch.

1. Create and push a git tag for the release.

2. Build and push the docker images. (See the above sections about docker images)

3. Build and push the manifest list.

4. Submit a patch to update docker image tags in our kubernetes manifests.

6. Review and merge the patch.
