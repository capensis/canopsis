<template>
  <v-list-item class="widget-settings-flat-item">
    <template #title>
      <slot name="title">{{ title }}</slot>
    </template>

    <template #append>
      <div class="mr-3">
        <slot name="actions" />
      </div>
    </template>
  </v-list-item>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
    },
  },
};
</script>

<style lang="scss">
.widget-settings-flat-item {
  .v-list-item-title {
    white-space: unset;
  }
}
</style>
