<template>
  <span
    v-html="preparedContent"
    :class="themeClass"
    class="v-icon"
  />
</template>
<script>
export default {
  inject: ['$system'],
  props: {
    content: {
      type: String,
      default: '',
    },
  },
  computed: {
    themeClass() {
      return `theme--${this.$system.dark ? 'dark' : 'light'}`;
    },

    preparedContent() {
      const element = document.createElement('div');

      element.innerHTML = this.content;

      const svg = element.getElementsByTagName('svg')[0];

      if (svg) {
        svg.classList.add('v-icon__component', this.themeClass);
      }

      return element.innerHTML;
    },
  },
};
</script>
