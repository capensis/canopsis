<template>
  <c-layout class="pa-3">
    <v-col :class="cardFlexClass" class="pa-0">
      <v-card class="alarms-expand-panel-tab-item-wrapper__card">
        <slot name="prepend" />
        <v-card-text>
          <slot />
        </v-card-text>
      </v-card>
    </v-col>
  </c-layout>
</template>

<script>
export default {
  props: {
    cardFlexClass: {
      type: [String, Array, Object],
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.alarms-expand-panel-tab-item-wrapper__card {
  margin: auto;
}
</style>
