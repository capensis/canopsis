<template>
  <c-action-btn
    v-bind="$attrs"
    custom-tooltip-activator="custom-tooltip-activator"
  >
    <template #button="">
      <v-btn
        :color="color"
        :small="small"
        fab="fab"
        dark="dark"
        @click.stop="$listeners.click"
      >
        <v-icon>{{ icon }}</v-icon>
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
    small: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
