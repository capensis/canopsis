<template>
  <v-list-group
    :disabled="disabled"
    class="widget-settings-group"
    color=""
    eager
  >
    <template #activator="{ props }">
      <v-list-item
        v-bind="props"
        :title="title"
      />
    </template>

    <v-list class="widget-settings-group__list px-2 py-0">
      <slot />
    </v-list>
  </v-list-group>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.widget-settings-group {
  &__list {
    background-color: rgba(var(--v-theme-on-surface), 0.04);
  }
}
</style>
