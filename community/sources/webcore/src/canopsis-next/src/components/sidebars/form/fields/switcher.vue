<template>
  <v-container
    class="pa-3"
    fluid
  >
    <v-layout
      align-center
      justify-space-between
    >
      <v-flex xs11="xs11">
        <div class="subheading">
          {{ title }}
        </div>
      </v-flex>
      <v-layout>
        <v-switch
          class="ma-0 pa-0"
          :input-value="value"
          color="primary"
          hide-details
          @change="$emit('input', $event)"
        />
      </v-layout>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
    },
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
