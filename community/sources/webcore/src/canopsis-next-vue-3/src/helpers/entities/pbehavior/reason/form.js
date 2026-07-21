export const pbehaviorReasonToForm = (reason = {}) => ({
  name: reason.name ?? '',
  description: reason.description ?? '',
  visible: !(reason.hidden ?? false),
})

export const formToPbehaviorReason = (form = {}) => ({
  ...form,

  hidden: !form.visible,
})
