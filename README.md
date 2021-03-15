# Puppet Master

Authentication and Authorization app in Go.

## TODO Functions (Draft)

### Permissions

**getAllPermissions() []*Permission***
Lists all permissions.

**getPermission(id int) *Permission***
Gets info for a given permission.

**createPermission() *Permission***
Creates a new permission.

**updatePermission(user *Permission) *Permission***
Updates a given permission info.

**deletePermission(id int) bool***
Deletes a given permission.

### Roles

**getAllRoles() []*Role***
Lists all roles.

**getRole(id int) *Role***
Gets info for a given role.

**createRole() *Role***
Creates a new role.

**updateRole(user *Role) *Role***
Updates a given role info.

**deleteRole(id int) bool***
Deletes a given role.

### Users

**getAllUsers() []*User***
Lists all users.

**getUser(id int) *User***
Gets info for a given user.

**createUser() *User***
Creates a new User.

**updateUser(user *User) *User***
Updates a given user info.

**deleteUser(id int) bool***
Deletes a given user.

**isSuperAdmin(userID int) bool**
Check if the user is Super Admin.

**checkPermission(userID int) bool**
Check user permissions.

**getPermissions(userID int) []int**
Get user permissions.

