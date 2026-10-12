<template>
  <span>
    <slot />
    <c-help-icon
      :text="$t('alarm.tooltips.priority')"
      icon-class="priority-help-icon ml-1"
      size="16"
      top
    />
  </span>
</template>

<style lang="scss" scoped>
.priority-help-icon ::v-deep .v-icon {
  opacity: 1 !important;
  transform: none !important;
}
</style>
