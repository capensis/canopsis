<template>
  <v-row v-bind="$attrs" justify="center">
    <v-col :cols="cols" :offset="offset">
      <v-card class="pa-3">
        <slot name="prepend" />
        <slot />
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup>
defineOptions({
  inheritAttrs: false,
});

defineProps({
  cols: {
    type: [Number, String],
    default: 12,
  },
  offset: {
    type: [Number, String],
    default: undefined,
  },
});
</script>
