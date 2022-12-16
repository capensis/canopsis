<template lang="pug">
  v-textarea.c-payload-field(
    ref="textarea",
    v-validate="",
    v-field="value",
    :label="label",
    :name="name",
    :rows="rows",
    :readonly="readonly",
    :disabled="disabled",
    :error-messages="errors.collect(name)",
    :row-height="lineHeight",
    :style="textareaStyle",
    auto-grow,
    @change="debouncedOnSelectionChange"
  )
    template(#prepend-inner="")
      div(:style="{ width: `${lineHeight}px` }")
    template(#append="")
      variables-menu(
        v-if="variables",
        :variables="availableVariables",
        :visible="variablesShown",
        :value="variablesMenuValue",
        :position-x="variablesMenuPosition.x",
        :position-y="variablesMenuPosition.y",
        show-value,
        @input="pasteVariable"
      )
      span.c-payload-field__lines(:style="linesStyle")
        span.c-payload-field__line(v-for="(line, index) in lines", :key="index", :style="lineStyle")
          span.c-payload-field__fake-line(v-if="selectedVariable && index === selectedVariable.index")
            | {{ line.text.slice(0, selectedVariable.start) }}
            span.c-payload-field__highlight(ref="variable")
              | {{ line.text.slice(selectedVariable.start, selectedVariable.end) }}
          v-tooltip(v-if="line.error", top)
            template(#activator="{ on }")
              v-icon.c-payload-field__warning-icon(v-on="on", :size="lineHeight", color="error") warning
            span {{ line.error.message }}
          | {{ line.text }}
</template>

<script>
import { debounce, keyBy } from 'lodash';

import {
  findSelectedVariable,
  matchPayloadOperators,
  matchPayloadVariables,
} from '@/helpers/payload-json';

import { formBaseMixin } from '@/mixins/form';

import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

export default {
  inject: ['$validator'],
  components: { VariablesMenu },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'json',
    },
    rows: {
      type: [Number, String],
      default: 5,
    },
    variables: {
      type: Array,
      /** TODO: Should be removed after integrate */
      default: () => [{
        value: '.Alarms',
        enumerable: true,
        variables: [{
          value: '.Value.Component',
          text: 'Component',
        }],
      }, {
        value: '.Entity',
        enumerable: true,
        variables: [{
          value: '.Infos.%infos_name%.Value',
          text: 'Infos value',
        }],
      }],
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    lineHeight: {
      type: Number,
      default: 18,
    },
    linesErrors: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      selectionVariableStart: 0,
      selectionVariableEnd: 0,
      variableGroup: undefined,
      operatorGroup: undefined,
      newVariableGroup: undefined,
      variablesShown: false,
      variablesMenuPosition: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    lines() {
      return this.value.split(/\n/).map((text, index) => ({
        text,
        error: this.linesErrorsByLineNumber[index + 1],
      }));
    },

    selectedVariable() {
      if (!this.variablesShown) {
        return undefined;
      }

      let end = this.selectionVariableEnd;

      for (let index = 0; index < this.lines.length; index += 1) {
        const { text } = this.lines[index];

        const lineCharactersCount = text.length + 1;

        if (end < lineCharactersCount) {
          return {
            start: Math.max(0, this.selectionVariableStart - (this.selectionVariableEnd - end)),
            end,
            index,
          };
        }

        end -= lineCharactersCount;
      }

      return undefined;
    },

    availableVariables() {
      return this.variables.reduce((acc, variable) => {
        if (variable.enumerable) {
          acc.push(...variable.variables.map(({ value, text }) => ({
            text,
            value: (this.variableGroup || this.newVariableGroup) && this.operatorGroup
              ? `{{ ${value} }}`
              : `{{ range ${variable.value} }} {{ ${value} }} {{ end }}`,
          })));
        }

        return acc;
      }, []);
    },

    variablesMenuValue() {
      if (this.variableGroup) {
        return this.variableGroup[0];
      }

      return this.operatorGroup && this.operatorGroup[0];
    },

    valueVariables() {
      return matchPayloadVariables(this.value);
    },

    valueOperators() {
      return matchPayloadOperators(this.value);
    },

    linesErrorsByLineNumber() {
      return keyBy(this.linesErrors, 'lineNumber');
    },

    lineHeightPixel() {
      return `${this.lineHeight}px`;
    },

    lineStyle() {
      return {
        lineHeight: this.lineHeightPixel,
        minHeight: this.lineHeightPixel,
      };
    },

    textareaStyle() {
      return {
        lineHeight: this.lineHeightPixel,
        fontSize: this.lineHeightPixel,
      };
    },

    linesStyle() {
      return {
        marginLeft: this.lineHeightPixel,
        maxWidth: `calc(100% - ${this.lineHeightPixel})`,
      };
    },
  },
  created() {
    this.debouncedOnSelectionChange = debounce(this.onSelectionChange, 100);
  },
  mounted() {
    document.addEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  beforeDestroy() {
    document.removeEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  methods: {
    setVariableSelection(start, end) {
      this.selectionVariableStart = start;
      this.selectionVariableEnd = end;
    },

    setVariableSelectionByGroup(group) {
      const [value] = group;

      this.setVariableSelection(group.index, group.index + value.length);
    },

    resetVariableSelection() {
      this.selectionVariableStart = undefined;
      this.selectionVariableEnd = undefined;
    },

    pasteVariable(variable) {
      const newValue = `${this.value.slice(0, this.selectionVariableStart)}${variable}${this.value.slice(this.selectionVariableEnd)}`;

      this.updateModel(newValue);
      this.resetVariableSelection();
      this.hideVariablesMenu();
    },

    showVariablesMenu() {
      this.variablesShown = true;

      this.$nextTick(() => {
        const [variableElement] = this.$refs.variable;
        const { top, left, height } = variableElement.getBoundingClientRect();

        this.variablesMenuPosition.x = left;
        this.variablesMenuPosition.y = top + height;
      });
    },

    hideVariablesMenu() {
      this.variablesShown = false;
    },

    isVariableCreatingInsideOperatorContent() {
      const { index, groups } = this.operatorGroup;

      return (
        this.newVariableGroup.index > index + groups.open.length
        && this.newVariableGroup.index < index + groups.variable.length - groups.close.length
      );
    },

    onSelectionChange() {
      if (!this.$el.contains(document.activeElement)) {
        return;
      }

      const { selectionStart, selectionEnd } = this.$refs.textarea.$refs.input;

      this.variableGroup = findSelectedVariable(
        this.valueVariables,
        selectionStart,
        selectionEnd,
      );
      this.operatorGroup = findSelectedVariable(
        this.valueOperators,
        selectionStart,
        selectionEnd,
      );
      this.newVariableGroup = this.value.slice(0, selectionStart).match(/({{){1,2}$/);

      if (this.newVariableGroup && !this.variableGroup) {
        if (
          !this.operatorGroup
          || this.isVariableCreatingInsideOperatorContent()
        ) {
          this.setVariableSelection(this.newVariableGroup.index, selectionEnd);
          this.showVariablesMenu();
          return;
        }
      }

      if (this.variableGroup || this.operatorGroup) {
        this.setVariableSelectionByGroup(this.variableGroup || this.operatorGroup);
        this.showVariablesMenu();
        return;
      }

      this.hideVariablesMenu();
      this.resetVariableSelection();
    },
  },
};
</script>

<style lang="scss">
$iconBarWidth: 18px;

.c-payload-field {
  .v-input__append-inner {
    margin: 0;
    padding: 0;

    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;

    pointer-events: none;
  }

  .v-input__prepend-inner {
    padding: 0;
  }

  textarea {
    line-height: inherit;
  }

  &__lines  {
    display: flex;
    flex-direction: column;

    padding: 7px 0 8px 0;
    max-height: 100%;
  }

  &__line, &__fake-line {
    white-space: pre-wrap;
    word-break: normal;
    text-align: start;
    overflow-wrap: break-word;
    color: transparent;
  }

  &__fake-line {
    position: absolute;
  }

  &__highlight {
    outline: 1px solid grey;
    border-radius: 2px;
    background: rgba(grey, 0.2);
  }

  &__warning-icon {
    position: absolute;
    left: 0;
    z-index: 1;
    pointer-events: auto;
  }
}
</style>
