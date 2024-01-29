<template>
  <c-action-btn v-bind="$attrs">
    <template #button="{ on }">
      <v-btn
        v-on="on"
        :color="color"
        small
        fab
        dark
        @click.stop="$listeners.click"
      >
        <v-icon small>
          {{ icon }}
        </v-icon>
      </v-btn>
    </template>
  </c-action-btn>
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    icon: {
      type: String,
      required: true,
    },
    color: {
      type: String,
      default: '',
    },
  },
};
</script>
