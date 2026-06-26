import { GROUPS_NAVIGATION_TYPES } from '@/constants'

import { mapIds } from '@/helpers/array'

export const userToForm = (user = {}) => ({
  _id: user._id ?? '',
  name: user.name ?? '',
  firstname: user.firstname ?? '',
  lastname: user.lastname ?? '',
  email: user.email ?? '',
  password: user.password ?? '',
  roles: user.roles || [],
  enabled: user.enabled ?? true,
  defaultview: user.defaultview ? user.defaultview._id : '',
  ui_language: user.ui_language ?? '',
  ui_tours: user.ui_tours ?? {},
  ui_groups_navigation_type: user.ui_groups_navigation_type ?? GROUPS_NAVIGATION_TYPES.sideBar,
  ui_theme: user.ui_theme ?? '',
})

export const formToUserRequest = form => ({
  ...form,

  roles: mapIds(form.roles),
})
