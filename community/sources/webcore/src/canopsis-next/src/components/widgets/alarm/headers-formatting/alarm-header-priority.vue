<template lang="pug">
  span
    slot
    c-help-icon.priority-help-icon.ml-1(:text="$t('alarmList.tooltips.priority')", size="16", top)
</template>

<style lang="scss" scoped>
.priority-help-icon ::v-deep .v-icon {
  opacity: 1 !important;
  transform: none !important;
}
</style>
