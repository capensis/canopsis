<template lang="pug">
  span
    slot
    c-help-icon.priority-help-icon.ml-1(:text="$t('alarm.tooltips.priority')", size="16", top)
</template>

<style lang="scss" scoped>
.priority-help-icon /deep/ .v-icon {
  opacity: 1 !important;
  transform: none !important;
}
</style>
