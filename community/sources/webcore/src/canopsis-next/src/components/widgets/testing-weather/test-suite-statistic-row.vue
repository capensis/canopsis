<template>
  <v-layout>
    <v-flex xs6>
      {{ count }}
    </v-flex>
    <v-flex
      class="text-lowercase"
      xs6
    >
      {{ label }}
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
    count: {
      type: [Number, String],
      default: 0,
    },
  },
};
</script>
