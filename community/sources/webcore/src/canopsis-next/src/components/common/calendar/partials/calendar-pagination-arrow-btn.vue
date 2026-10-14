<template>
  <c-tooltip location="bottom">
    <template #activator="{ props: tooltipProps }">
      <v-btn
        v-bind="tooltipProps"
        class="mx-2"
        variant="text"
        size="small"
        icon
        @click="$emit('click', $event)"
      >
        <v-icon>{{ icon }}</v-icon>
      </v-btn>
    </template>

    <span>{{ label }}</span>
  </c-tooltip>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
  },

  emits: ['click'],
};
</script>
