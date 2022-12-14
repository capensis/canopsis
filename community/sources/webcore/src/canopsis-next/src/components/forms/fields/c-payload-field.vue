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
    @update:searchInput="debouncedOnSelectionChange"
  )
    template(#prepend-inner="")
      div(:style="{ width: `${lineHeight}px` }")
    template(#append="")
      span.c-payload-field__lines(:style="linesStyle")
        span.c-payload-field__line(v-for="(line, index) in lines", :key="index", :style="lineStyle")
          span.c-payload-field__fake-line(v-if="selectedVariable && index === selectedVariable.index")
            | {{ line.text.slice(0, selectedVariable.start) }}
            span.c-payload-field__highlight {{ line.text.slice(selectedVariable.start, selectedVariable.end) }}
          v-tooltip(v-if="line.error", top)
            template(#activator="{ on }")
              v-icon.c-payload-field__warning-icon(v-on="on", :size="lineHeight", color="error") warning
            span {{ line.error.message }}
          | {{ line.text }}
</template>

<script>
import { debounce, keyBy } from 'lodash';

export default {
  inject: ['$validator'],
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
      default: () => [],
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
      return this.variables;
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
    this.debouncedOnSelectionChange = debounce(this.onSelectionChange, 50);
  },
  mounted() {
    document.addEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  beforeDestroy() {
    document.removeEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  methods: {
    onSelectionChange() {
      if (!this.$el.contains(document.activeElement)) {
        return;
      }

      const { selectionStart, selectionEnd } = this.$refs.textarea.$refs.input;
      /** TODO: Need to check variable and replace on variable coordinate */
      this.selectionVariableStart = selectionStart;
      this.selectionVariableEnd = selectionEnd;
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
    outline: 1px solid red;
  }

  &__warning-icon {
    position: absolute;
    left: 0;
    z-index: 1;
    pointer-events: auto;
  }
}
</style>
