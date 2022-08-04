<template lang="pug">
  v-layout.mermaid-editor
    v-flex.mermaid-editor__field(xs4)
      v-textarea.mt-0.pa-0(v-model="form.code", :error="!isValidMermaid", rows="27", hide-details)
    v-flex(xs8)
      v-layout.fill-height(column)
        v-flex.mermaid-editor__toolbar
          v-layout(row, align-center, justify-end)
            v-btn-toggle.mr-2
              v-btn.ma-0(v-model="addOnClick", flat, large)
                v-icon(left) add_location
                span {{ $t('mermaid.addPoint') }}
            v-flex(xs4)
              mermaid-theme-field.my-2(v-field="form.theme")
        v-flex
          div.position-relative
            mermaid-preview(:value="form.code", :config="config", :error="!isValidMermaid")
            mermaid-points.mermaid-editor__points(v-field="form.points", :add-on-click="addOnClick")
</template>

<script>
import { MERMAID_THEME_PROPERTIES_BY_NAME } from '@/constants';

import { isValidMermaidDiagram } from '@/helpers/mermaid';

import { formMixin } from '@/mixins/form';

import MermaidThemeField from './partials/mermaid-theme-field.vue';
import MermaidPreview from './mermaid-preview.vue';
import MermaidPoints from './mermaid-points.vue';

export default {
  components: { MermaidPoints, MermaidThemeField, MermaidPreview },
  mixins: [formMixin],
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
  data() {
    return {
      addOnClick: false,
    };
  },
  computed: {
    isValidMermaid() {
      return isValidMermaidDiagram(this.form.code);
    },

    config() {
      const themeProperties = MERMAID_THEME_PROPERTIES_BY_NAME[this.form.theme] ?? {};

      return {
        theme: this.form.theme,
        ...themeProperties,

        er: {
          useMaxWidth: false,
        },
        pie: {
          useMaxWidth: false,
        },
        sequence: {
          useMaxWidth: false,
        },
        flowchart: {
        },
        requirement: {
          useMaxWidth: false,
        },
      };
    },
  },
};
</script>

<style lang="scss">
$borderColor: #e5e5e5;

.mermaid-editor {
  border: 1px solid $borderColor;
  min-height: 500px;

  &__field {
    border-right: 1px solid $borderColor;
  }

  &__toolbar {
    height: 60px;
    max-height: 60px;
    border-bottom: 1px solid $borderColor;
  }

  &__points {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }
}
</style>
