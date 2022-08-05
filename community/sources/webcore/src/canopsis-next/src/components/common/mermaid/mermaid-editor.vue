<template lang="pug">
  v-layout.mermaid-editor
    v-flex.mermaid-editor__sidebar
      mermaid-code-editor.fill-height(v-field="form.code")
    v-flex.mermaid-editor__content
      v-layout.fill-height(column)
        v-flex.mermaid-editor__toolbar
          v-layout(row, align-center, justify-end)
            v-flex(xs4)
              mermaid-theme-field.my-2(v-field="form.theme")
        v-flex
          mermaid-preview(:value="form.code", :config="config")
</template>

<script>
import { MERMAID_THEME_PROPERTIES_BY_NAME } from '@/constants';

import MermaidCodeEditor from './partials/mermaid-code-editor.vue';
import MermaidThemeField from './partials/mermaid-theme-field.vue';
import MermaidPreview from './mermaid-preview.vue';

export default {
  components: { MermaidCodeEditor, MermaidThemeField, MermaidPreview },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    config() {
      const themeProperties = MERMAID_THEME_PROPERTIES_BY_NAME[this.form.theme] ?? {};

      return {
        theme: this.form.theme,
        ...themeProperties,

        er: {
          useMaxWidth: false,
        },
        sequence: {
          useMaxWidth: false,
        },
        flowchart: {
          useMaxWidth: false,
        },
      };
    },
  },
};
</script>

<style lang="scss">
$borderColor: #e5e5e5;
$sideBarWidth: 500px;
$contentWidth: 800px;

.mermaid-editor {
  min-height: 500px;

  &__sidebar {
    width: $sideBarWidth;
    max-width: $sideBarWidth;
    border: 1px solid $borderColor;
    border-right: none;
  }

  &__content {
    width: $contentWidth;
    max-width: $contentWidth;
    border: 1px solid $borderColor;
  }

  &__toolbar {
    height: 60px;
    max-height: 60px;
    border-bottom: 1px solid $borderColor;
  }
}
</style>
