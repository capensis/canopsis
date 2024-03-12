<template>
  <v-layout class="mb-2 state-setting-information-row">
    <v-flex class="state-setting-information-row--label">
      <span class="font-weight-bold">{{ label }}</span>
    </v-flex>
    <v-flex>
      <slot>{{ value }}{{ unit }}</slot>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    value: {
      type: [String, Number],
      required: false,
    },
    unit: {
      type: String,
      required: false,
    },
  },
};
</script>

<style lang="scss">
.state-setting-information-row {
  &--label {
    max-width: 100px;
  }
}
</style>
