<template lang="pug">
  v-list-tile.subheading
    v-list-tile-avatar(v-if="!hideIcon")
      v-icon {{ icon }}
    v-list-tile-content
      v-layout(align-center)
        div.subheading
          span.font-weight-bold {{ label }}
          span(v-if="value") : {{ value }}
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    value: {
      type: String,
      default: '',
    },
    hideIcon: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
