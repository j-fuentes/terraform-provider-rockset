---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rockset_role Resource - terraform-provider-rockset"
subcategory: ""
description: |-
  Manages a Rockset Role.
---

# rockset_role (Resource)

Manages a Rockset Role.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Role name.

### Optional

- `description` (String) Role description.
- `id` (String) The ID of this resource.
- `privilege` (Block Set) Privileges associated with the role. (see [below for nested schema](#nestedblock--privilege))

### Read-Only

- `created_at` (String) When the role was created.
- `created_by` (String) Who created the role.
- `owner_email` (String) The email of the user who currently owns the role.

<a id="nestedblock--privilege"></a>
### Nested Schema for `privilege`

Required:

- `action` (String) The action allowed by this privilege.

Optional:

- `cluster` (String) Rockset cluster ID for which this action is allowed. Only applies to Workspace actions. Defaults to '*ALL*' if not specified.
- `resource_name` (String) The resource on which this action is allowed. Defaults to 'All' if not specified.


