<template>
  <v-row>
    <v-col cols="6">
      {{ count }}
    </v-col>
    <v-col class="text-lowercase" cols="6">
      {{ label }}
    </v-col>
  </v-row>
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
