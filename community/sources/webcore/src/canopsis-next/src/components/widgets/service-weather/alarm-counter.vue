<template>
  <v-layout
    justify-space-around
    align-center
    v-on="$listeners"
  >
    <v-icon
      :style="{ color, caretColor: color }"
      size="18"
    >
      {{ icon }}
    </v-icon>
    <span class="text-caption">{{ count }}</span>
  </v-layout>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
    count: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'white',
    },
  },
};
</script>
