package gql

// func (r *Resolver) checkPermission(ctx context.Context) bool {
// 	auth := ctx.Value(domain.ContextKeyID).(*domain.Auth)

// userID := int64(userInfo["user_id"].(float64))
// roles := []string{}

// for _, role := range userInfo["roles"].([]interface{}) {
// 	roles = append(roles, role.(string))
// }

// auth := &domain.Auth{
// 	UserID: userID,
// 	Name:   userInfo["name"].(string),
// 	Email:  userInfo["email"].(string),
// 	Roles:  roles,
// }

// 	if len(auth.UserPermissions) == 0 {
// 		return false
// 	}

// 	userPermissions, err := r.permissionUseCase.GetPermissionsByUserID(ctx, auth.UserID)
// 	if err != nil {
// 		return false
// 	}

// 	for _, role := range auth.Roles {
// 		if role == "admin" {
// 			return true
// 		}
// 	}

// 	for index, userPermission := range userPermissions {
// 		if auth.UserPermissions[index] != userPermission.Name {
// 			return false
// 		}
// 	}

// 	return true
// }
