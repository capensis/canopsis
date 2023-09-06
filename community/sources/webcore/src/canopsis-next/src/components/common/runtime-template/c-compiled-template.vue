<template lang="pug">
  c-runtime-template(v-on="$listeners", v-bind="$attrs", :template="compiledTemplate", :parent="$parent")
</template>

<script>
import { compile } from '@/helpers/handlebars';

export default {
  inheritAttrs: false,
  props: {
    template: {
      type: String,
      default: '',
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
  data() {
    return {
      compiledTemplate: '',
    };
  },
  watch: {
    template: 'compileTemplate',
    context: 'compileTemplate',
    parentElement: 'compileTemplate',
  },
  created() {
    this.compileTemplate();
  },
  methods: {
    async compileTemplate() {
      const compiledTemplate = await compile(this.template, this.context);

      this.compiledTemplate = `<${this.parentElement}>${compiledTemplate}</${this.parentElement}>`;
    },
  },
};
</script>
