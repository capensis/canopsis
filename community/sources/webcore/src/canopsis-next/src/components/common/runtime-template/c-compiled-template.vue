<template lang="pug">
  c-runtime-template(v-on="$listeners", :template="compiledTemplate", :parent="$parent")
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
    context: {
      type: Object,
      required: false,
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.template, this.context);

        return `<${this.parentElement}>${compiledTemplate}</${this.parentElement}>`;
      },
      default: '',
    },
  },
};
</script>
