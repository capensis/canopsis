<template>
  <v-layout column>
    <v-layout
      class="mermaid-editor mb-2"
      :style="editorStyles"
    >
      <v-flex class="mermaid-editor__sidebar">
        <mermaid-code-editor
          class="fill-height"
          v-field="form.code"
        />
      </v-flex>
      <v-flex class="mermaid-editor__content">
        <v-layout
          class="mermaid-editor__toolbar px-2"
          align-center
          justify-end
        >
          <add-location-btn
            class="mr-2"
            v-model="addOnClick"
          />
          <mermaid-theme-field
            class="mermaid-editor__theme-picker"
            v-field="form.theme"
          />
        </v-layout>
        <div class="mermaid-editor__preview">
          <mermaid-code-preview
            :value="form.code"
            :theme="form.theme"
          />
          <mermaid-points-editor
            class="mermaid-editor__points"
            v-field="form.points"
            :add-on-click="addOnClick"
          />
        </div>
      </v-flex>
    </v-layout>
    <v-messages
      v-if="hasChildrenError"
      :value="errorMessages"
      color="error"
    />
  </v-layout>
</template>

<script>
import { COLORS } from '@/config';

import { formMixin, validationChildrenMixin } from '@/mixins/form';

import MermaidCodePreview from '@/components/other/map/partials/mermaid-code-preview.vue';

import AddLocationBtn from './add-location-btn.vue';
import MermaidCodeEditor from './mermaid-code-editor.vue';
import MermaidThemeField from './mermaid-theme-field.vue';
import MermaidPointsEditor from './mermaid-points-editor.vue';

export default {
  inject: ['$validator'],
  components: {
    AddLocationBtn,
    MermaidCodeEditor,
    MermaidThemeField,
    MermaidCodePreview,
    MermaidPointsEditor,
  },
  mixins: [formMixin, validationChildrenMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
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
  data() {
    return {
      addOnClick: false,
    };
  },
  computed: {
    errorMessages() {
      return [this.$t('mermaid.errors.emptyMermaid')];
    },

    editorStyles() {
      const minHeight = Math.max(...this.form.points.map(({ y }) => y), this.minHeight);

      return {
        minHeight: `${minHeight}px`,
        borderColor: this.hasChildrenError ? COLORS.error : undefined,
      };
    },
  },
  watch: {
    form: {
      deep: true,
      handler() {
        if (this.hasChildrenError) {
          this.$validator.validate(this.name);
        }
      },
    },
  },
  mounted() {
    this.attachRequiredRule();
  },
  beforeDestroy() {
    this.detachRules();
  },
  methods: {
    attachRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => !!this.form.code && !!this.form.points.length,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.name);
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
