<template>
  <v-layout column>
    <v-layout
      class="secondary lighten-2 tab-panel-item"
      align-center
    >
      <slot name="title">
        <span class="white--text">{{ tab.title }}</span>
      </slot>
    </v-layout>
    <v-divider dark />
  </v-layout>
</template>

<script>
export default {
  props: {
    tab: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
  .tab-panel-item {
    display: flex;
    align-items: center;
    height: 48px;
  }
</style>
