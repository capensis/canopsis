<template>
  <c-runtime-template
    v-bind="$attrs"
    :template="compiledTemplate"
    :parent="$parent"
    v-on="$listeners"
  />
</template>

<script>
import { sanitizeHtml, linkifyHtml, normalizeHtml } from '@/helpers/html';
import { compile, runTemplate, hasTemplate } from '@/helpers/handlebars';

export default {
  inject: ['$system'],
  inheritAttrs: false,
  props: {
    template: {
      type: String,
      default: '',
    },
    templateId: {
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
    sanitizeOptions: {
      type: Object,
      required: false,
    },
    linkifyOptions: {
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
    sanitizeOptions: 'compileTemplate',
    linkifyOptions: 'compileTemplate',
    '$system.theme': 'compileTemplate',
  },
  created() {
    this.compileTemplate();
  },
  methods: {
    async compileTemplate() {
      try {
        const context = {
          theme: {
            ...this.$system.theme,

            dark: this.$system.dark,
          },

          ...this.context,
        };

        const compiledTemplate = hasTemplate(this.templateId)
          ? await runTemplate(this.templateId, context)
          : await compile(this.template, context);

        this.compiledTemplate = `<${this.parentElement}>${
          normalizeHtml(
            sanitizeHtml(
              linkifyHtml(compiledTemplate, this.linkifyOptions),
              this.sanitizeOptions,
            ),
          )
        }</${this.parentElement}>`;
      } catch (err) {
        console.error(err);

        this.compiledTemplate = '';
      }
    },
  },
};
</script>
