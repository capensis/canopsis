<template>
  <v-btn
    :loading="loading"
    icon
    small
    @click="$emit('expand', !expanded)"
  >
    <v-icon :color="color">
      {{ expanded ? 'keyboard_arrow_up' : 'keyboard_arrow_down' }}
    </v-icon>
  </v-btn>
</template>

<script>
export default {
  model: {
    prop: 'expanded',
    event: 'expand',
  },
  props: {
    expanded: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
