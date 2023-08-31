<template lang="pug">
  div.position-relative
    c-compiled-template(:template="template")
</template>

<script>
export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    template() {
      return this.widget.parameters?.template;
    },
  },
};
</script>
