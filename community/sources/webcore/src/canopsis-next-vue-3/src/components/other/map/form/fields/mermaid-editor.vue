<template>
  <c-layout column>
    <c-layout :style="editorStyles" class="mermaid-editor mb-2">
      <v-col class="mermaid-editor__sidebar pa-0">
        <mermaid-code-editor
          :value="form.code"
          class="fill-height"
          @input="updateField('code', $event)"
        />
      </v-col>
      <v-col class="mermaid-editor__content pa-0">
        <c-layout class="mermaid-editor__toolbar px-2" align="center" justify="end">
          <add-location-btn :value="addOnClick" class="mr-2" @input="addOnClick = $event" />
          <mermaid-theme-field
            :value="form.theme"
            class="mermaid-editor__theme-picker"
            @input="updateField('theme', $event)"
          />
        </c-layout>
        <div class="mermaid-editor__preview">
          <mermaid-code-preview :value="form.code" :theme="form.theme" />
          <mermaid-points-editor
            :value="form.points"
            :add-on-click="addOnClick"
            class="mermaid-editor__points"
            @input="updateField('points', $event)"
          />
        </div>
      </v-col>
    </c-layout>
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import MermaidCodePreview from '@/components/other/map/partials/mermaid-code-preview.vue';

import AddLocationBtn from './add-location-btn.vue';
import MermaidCodeEditor from './mermaid-code-editor.vue';
import MermaidThemeField from './mermaid-theme-field.vue';
import MermaidPointsEditor from './mermaid-points-editor.vue';

export default {
  components: {
    AddLocationBtn,
    MermaidCodeEditor,
    MermaidThemeField,
    MermaidCodePreview,
    MermaidPointsEditor,
  },
  mixins: [formMixin],
  props: {
    form: {
      type: Object,
      required: true,
    },
    minHeight: {
      type: Number,
      default: 500,
    },
    name: {
      type: String,
      default: 'parameters',
    },
  },
  emits: ['update:form'],
  data() {
    return {
      addOnClick: false,
    };
  },
  computed: {
    editorStyles() {
      const minHeight = Math.max(...this.form.points.map(({ y }) => y), this.minHeight);

      return {
        minHeight: `${minHeight}px`,
      };
    },
  },
};
</script>

<style lang="scss">
$borderColor: #e5e5e5;
$contentBackgroundColor: #f9f9f9;
$contentWidth: 800px;
$toolbarHeight: 60px;

.mermaid-editor {
  border: 1px solid $borderColor;

  &__sidebar {
    flex: 1;
    min-width: 200px;
  }

  &__content {
    background: $contentBackgroundColor;

    border-left: 1px solid $borderColor;

    position: relative;

    width: $contentWidth;
    max-width: $contentWidth;
  }

  &__toolbar {
    background: white;

    position: absolute;
    right: 0;

    height: $toolbarHeight;
    max-height: $toolbarHeight;

    border-bottom: 1px solid $borderColor;
    border-left: 1px solid $borderColor;
  }

  &__preview {
    position: relative;

    margin-top: $toolbarHeight;
  }

  &__theme-picker {
    width: 200px;
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
