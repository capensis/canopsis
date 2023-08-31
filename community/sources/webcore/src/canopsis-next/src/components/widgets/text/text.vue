<template lang="pug">
  div.position-relative
    c-runtime-template(:template="compiledTemplate")
</template>

<script>
import { compile } from '@/helpers/handlebars';

export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.widget.parameters?.template ?? '');

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
};
</script>
