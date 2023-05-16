<template lang="pug">
  c-runtime-template(:template="compiledTemplate")
</template>

<script>
import { compile } from '@/helpers/handlebars';

export default {
  props: {
    template: {
      type: String,
      required: false,
    },
    parentElement: {
      type: String,
      default: 'div',
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.template, this.variables);

        return `<${this.parentElement}>${compiledTemplate}</${this.parentElement}>`;
      },
      default: '',
    },
  },
};
</script>
