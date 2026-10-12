<template>
  <v-window-item v-bind="$attrs" class="pa-3">
    <c-expand-panel-card :cols="cols" :offset="offset">
      <template v-if="$slots.prepend" #prepend>
        <slot name="prepend" />
      </template>
      <slot />
    </c-expand-panel-card>
  </v-window-item>
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
