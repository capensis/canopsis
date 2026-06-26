<template>
  <v-tooltip :location="top ? 'top' : 'bottom'">
    <template #activator="{ props }">
      <slot
        :on="props"
        :props="props"
        name="activator"
      />
    </template>
    <span>{{ content }}</span>
  </v-tooltip>
</template>

<script>
export default {
  props: {
    content: {
      type: String,
      default: '',
    },
    top: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
