<template>
  <c-layout gap="1" column>
    <v-btn
      v-bind="$attrs"
      color="primary"
    >
      <slot />
    </v-btn>
    <div v-if="error" class="text-error text-caption">
      {{ error }}
    </div>
  </c-layout>
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    error: {
      type: String,
      default: '',
    },
  },
};
</script>
