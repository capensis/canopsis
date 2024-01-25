<template>
  <v-layout
    class="mt-3"
    wrap
  >
    <v-flex
      v-show="label"
      xs12
    >
      <div class="text-body-2">
        {{ label }}
      </div>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      default: '',
    },
  },
};
</script>
