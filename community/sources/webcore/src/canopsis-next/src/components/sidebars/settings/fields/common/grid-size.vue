<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        span {{ title }}
        v-icon.ml-2 {{ sizeIcon }}
    v-container
      c-column-size-field(v-field="value", :mobile="mobile", :tablet="tablet")
</template>

<script>
export default {
  props: {
    value: {
      type: Number,
      required: false,
    },
    title: {
      type: String,
      required: true,
    },
    mobile: {
      type: Boolean,
      default: false,
    },
    tablet: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    sizeIcon() {
      if (this.mobile) {
        return 'phone_android';
      }

      if (this.tablet) {
        return 'tablet';
      }

      return 'laptop';
    },
  },
};
</script>
