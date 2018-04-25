<template lang="pug">
  div(ref="stickyWrapperElement")
    div(ref="stickyElement", :style="stickyStyle")
      slot
    div(v-if="keepSpace && hasSticky", :style="keepSpaceStyle")
</template>

<script>
export default {
  props: {
    keepSpace: Boolean,
  },
  data() {
    return {
      height: 0,
      hasSticky: false,
    };
  },

  created() {
    document.addEventListener('scroll', this.sticky);
  },
  mounted() {
    this.setHeight();
  },
  beforeDestroy() {
    document.removeEventListener('scroll', this.sticky);
  },
  computed: {
    keepSpaceStyle() {
      return {
        height: `${this.height}px`,
      };
    },
    stickyStyle() {
      if (this.hasSticky) {
        return {
          top: 0,
          position: 'fixed',
          'z-index': 100,
        };
      }

      return {};
    },
  },
  methods: {
    setHeight() {
      if (!this.keepSpace) {
        return;
      }

      this.height = this.$refs.stickyElement.getBoundingClientRect().height;
    },
    sticky() {
      this.hasSticky = this.$refs.stickyWrapperElement.getBoundingClientRect().top <= 0;
    },
  },
};
</script>
