import { PBEHAVIOR_TYPE_TYPES } from '@/constants'

export const pbehaviorTypeToForm = (type = {}) => ({
  name: type.name ?? '',
  description: type.description ?? '',
  type: type.type ?? PBEHAVIOR_TYPE_TYPES.active,
  priority: type.priority ?? 5,
  icon_name: type.icon_name ?? '',
  color: type.color ?? '',
  visible: !(type.hidden ?? false),
})

export const formToPbehaviorType = (form = {}) => ({
  ...form,

  hidden: !form.visible,
})
