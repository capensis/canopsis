<template>
  <v-textarea
    class="c-payload-textarea-field"
    ref="field"
    v-validate="rules"
    v-field="value"
    :label="label"
    :name="name"
    :rows="rows"
    :readonly="readonly"
    :disabled="disabled"
    :error-messages="payloadErrors.inline"
    :row-height="lineHeight"
    :style="textareaStyle"
    :error="!!linesErrors.length"
    auto-grow
    @blur="handleBlur"
    @input="debouncedOnSelectionChange"
  >
    <template #prepend-inner="">
      <div :style="{ width: errorsOffsetPixel }" />
    </template>
    <template #append="">
      <div class="c-payload-textarea-field__append">
        <variables-menu
          v-if="variables"
          :variables="availableVariables"
          :visible="variablesShown"
          :value="variablesMenuValue"
          :position-x="variablesMenuPosition.x"
          :position-y="variablesMenuPosition.y"
          ignore-click-outside
          show-value
          @input="pasteVariable"
        />
        <span
          class="c-payload-textarea-field__lines"
          :style="linesStyle"
        >
          <span
            class="c-payload-textarea-field__line"
            v-for="(line, index) in lines"
            :key="index"
            :style="lineStyle"
          >
            <span
              class="c-payload-textarea-field__fake-line"
              v-if="selectedVariable && index === selectedVariable.index"
            >
              <span>{{ line.text.slice(0, selectedVariable.start) }}</span>
              <span
                class="c-payload-textarea-field__highlight"
                ref="variable"
              >
                <span>{{ line.text.slice(selectedVariable.start, selectedVariable.end) }}</span>
              </span>
            </span>
            <v-tooltip
              v-if="line.error"
              top
            >
              <template #activator="{ on }">
                <v-icon
                  class="c-payload-textarea-field__warning-icon"
                  v-on="on"
                  :size="lineHeight"
                  color="error"
                >
                  warning
                </v-icon>
              </template>
              <span>{{ line.error.message }}</span>
            </v-tooltip>
            <span>{{ line.text }}</span>
          </span>
        </span>
      </div>
    </template>
  </v-textarea>
</template>

<script>
import { keyBy } from 'lodash';

import { payloadFieldMixin } from '@/mixins/payload/payload-field';

import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

export default {
  inject: ['$validator'],
  components: { VariablesMenu },
  mixins: [payloadFieldMixin],
  props: {
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
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

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

    payloadErrors() {
      return this.errors.collect(this.name, null, false).reduce((acc, item) => {
        if (item.msg.includes('|')) {
          const [line, message] = item.msg.split('|');

          acc.lines.push({ line, message });
        } else {
          acc.inline.push(item.msg);
        }

        return acc;
      }, {
        lines: [],
        inline: [],
      });
    },

    linesErrors() {
      return this.payloadErrors.lines;
    },

    linesErrorsByLineNumber() {
      return keyBy(this.linesErrors, 'line');
    },

    lineHeightPixel() {
      return `${this.lineHeight}px`;
    },

    errorsOffsetPixel() {
      return this.linesErrors.length ? this.lineHeightPixel : 0;
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
        marginLeft: this.errorsOffsetPixel,
        maxWidth: `calc(100% - ${this.lineHeightPixel})`,
      };
    },
  },
};
</script>

<style lang="scss">
$iconBarWidth: 18px;

.c-payload-textarea-field {
  .v-input__append-inner {
    pointer-events: none;
  }

  .v-input__prepend-inner {
    padding: 0 !important;
  }

  textarea {
    line-height: inherit;
  }

  &__lines  {
    display: flex;
    flex-direction: column;

    max-height: 100%;
  }

  &__line, &__fake-line {
    line-height: 16px;
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
    transform: translateX(-100%);
  }

  &__append {
    margin: 0;
    padding: 0;

    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
  }
}
</style>
