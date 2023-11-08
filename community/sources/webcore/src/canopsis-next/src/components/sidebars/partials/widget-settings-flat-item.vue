<template>
  <v-list-item class="widget-settings-flat-item">
    <v-list-item-content>
      <v-list-item-title>
        {{ title }}
      </v-list-item-title>
    </v-list-item-content>
    <v-list-item-action class="mr-3">
      <slot name="actions" />
    </v-list-item-action>
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
  .v-list-item__title {
    white-space: unset;
  }
}
</style>
