<template lang="pug">
  c-runtime-template(v-on="$listeners", v-bind="$attrs", :template="compiledTemplate", :parent="$parent")
</template>

<script>
import { sanitizeHtml, linkifyHtml, normalizeHtml } from '@/helpers/html';
import { compile } from '@/helpers/handlebars';

export default {
  inject: ['$system'],
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
    '$system.theme': 'compileTemplate',
  },
  created() {
    this.compileTemplate();
  },
  methods: {
    async compileTemplate() {
      try {
        const compiledTemplate = await compile(this.template, {
          theme: {
            ...this.$system.theme,

            dark: this.$system.dark,
          },

          ...this.context,
        });

        this.compiledTemplate = `<${this.parentElement}>${
          normalizeHtml(sanitizeHtml(linkifyHtml(compiledTemplate)))
        }</${this.parentElement}>`;
      } catch (err) {
        console.error(err);

        this.compiledTemplate = '';
      }
    },
  },
};
</script>
