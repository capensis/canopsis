<template>
  <c-textarea
    ref="field"
    :model-value="value"
    :label="label"
    :name="name"
    :rows="rows"
    :readonly="readonly"
    :disabled="disabled"
    :row-height="lineHeight"
    :style="textareaStyle"
    :error="!!linesErrors.length"
    :autofocus="autofocus"
    :wrap="wrap"
    :required="required"
    class="c-payload-textarea-field"
    auto-grow
    @blur="handleBlur"
    @update:model-value="updateModel"
    @input="debouncedOnSelectionChange"
  >
    <template #prepend-inner>
      <div :style="{ width: errorsOffsetPixel }" />
    </template>

    <template v-if="$slots.append || variables.length || lines.length" #append-inner>
      <div v-if="$slots.append" class="c-payload-textarea-field__append-slot">
        <slot name="append" />
      </div>

      <div class="c-payload-textarea-field__append">
        <variables-menu
          v-if="variables.length"
          :items="variables"
          :visible="variablesShown"
          :value="variablesMenuValue"
          :position-x="variablesMenuPosition.x"
          :position-y="variablesMenuPosition.y"
          ignore-click-outside
          show-value
          clickable-parent
          @update:value="pasteVariable"
          @close="hideVariablesMenu"
        />

        <span :style="linesStyle" class="c-payload-textarea-field__lines">
          <span
            v-for="(line, index) in lines"
            :key="index"
            :style="lineStyle"
            :class="{ 'c-payload-textarea-field__line--error': line.error }"
            class="c-payload-textarea-field__line"
          >
            <span
              v-if="selectedVariable && index === selectedVariable.index"
              class="c-payload-textarea-field__fake-line"
            >
              <span>{{ line.text.slice(0, selectedVariable.start) }}</span>
              <span ref="variable" class="c-payload-textarea-field__highlight">
                <span>{{ line.text.slice(selectedVariable.start, selectedVariable.end) }}</span>
              </span>
            </span>

            <v-tooltip v-if="line.error" location="top">
              <template #activator="{ props: activatorProps }">
                <c-icon
                  v-bind="activatorProps"
                  :size="lineHeight"
                  class="c-payload-textarea-field__warning-icon"
                  color="error"
                >
                  warning
                </c-icon>
              </template>

              <span>{{ line.error.message }}</span>
            </v-tooltip>

            <span class="c-payload-textarea-field__line-text" contenteditable="true">
              {{ line.text }}
            </span>
          </span>
        </span>
      </div>
    </template>
  </c-textarea>
</template>

<script>
import { keyBy } from 'lodash';

import { FORM_CONTEXT_KEY } from '@/constants';

import { getStringLinesCount } from '@/helpers/string';

import { payloadFieldMixin } from '@/mixins/payload/payload-field';

import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

export default {
  inject: {
    formContext: {
      from: FORM_CONTEXT_KEY,
      default: null,
    },
  },

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
    autofocus: {
      type: Boolean,
      default: false,
    },
    wrap: {
      type: String,
      default: 'on',
    },
    errorMessages: {
      type: [Array, String],
      default: () => [],
    },
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

    payloadErrors() {
      return this.normalizedErrorMessages.reduce(
        (acc, fullMessage) => {
          if (fullMessage.includes('|')) {
            const [line, message] = fullMessage.split('|');

            acc.lines.push({ line, message });
          } else {
            acc.inline.push(fullMessage);
          }

          return acc;
        },
        {
          lines: [],
          inline: [],
        },
      );
    },

    normalizedErrorMessages() {
      if (!this.errorMessages) {
        return [];
      }

      return Array.isArray(this.errorMessages) ? this.errorMessages : [this.errorMessages];
    },

    linesErrors() {
      return this.payloadErrors.lines;
    },

    inlineErrorMessage() {
      return this.payloadErrors.inline.join('\n') || undefined;
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

  watch: {
    inlineErrorMessage: {
      immediate: true,
      handler(message, oldMessage) {
        if (message) {
          this.formContext?.setFieldError(this.name, message);

          return;
        }

        const currentMessage = this.formContext?.errors?.value?.[this.name];

        if (!oldMessage || !currentMessage) {
          return;
        }

        const nextMessage =
          currentMessage
            .split('\n')
            .filter(part => part !== oldMessage)
            .join('\n') || undefined;

        this.formContext?.setFieldError(this.name, nextMessage);
      },
    },
  },

  mounted() {
    setTimeout(() => {
      const input = this.getInputElement();

      if (!input) {
        return;
      }

      const computedStyle = window.getComputedStyle(input);
      const lineHeight = parseFloat(computedStyle.lineHeight);
      const linesCount = getStringLinesCount(this.value);
      const height = Math.ceil(lineHeight * linesCount);
      const minHeight = parseInt(this.rows, 10) * parseFloat(this.lineHeight);

      input.style.height = `${Math.max(height, minHeight)}px`;
    }, 1000);
  },
};
</script>

<style lang="scss">
$iconBarWidth: 18px;

.c-payload-textarea-field {
  .v-field {
    position: relative;
  }

  .v-field__field {
    position: relative;
  }

  .v-field__append-inner {
    position: absolute;
    grid-area: 1 / 1 / -1 / -1;
    inset: 0;
    z-index: 1;
    width: 100%;
    height: 100%;
    padding: 0 !important;
    pointer-events: none;

    .v-icon {
      pointer-events: all;
    }
  }

  .v-field__prepend-inner {
    padding: 0 !important;
  }

  textarea {
    line-height: inherit;
  }

  &__lines {
    display: flex;
    flex-direction: column;

    max-height: 100%;
  }

  &__line,
  &__fake-line {
    line-height: 16px;
    white-space: pre-wrap;
    word-break: normal;
    text-align: start;
    overflow-wrap: break-word;
    color: transparent;
  }

  &__line {
    &-text {
      margin-left: -$iconBarWidth;
    }

    &--error {
      background: rgb(var(--v-theme-error));
      color: rgb(var(--v-theme-on-error));
    }
  }

  &__fake-line {
    position: absolute;
  }

  &__highlight {
    display: inline-block;
    line-height: inherit;
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
    inset: 0;

    .c-payload-textarea-field__lines {
      padding: var(--v-field-input-padding-top, var(--v-field-padding-top, 0))
        var(--v-field-padding-end, 16px) var(--v-field-input-padding-bottom, 0)
        var(--v-field-padding-start, 16px);
    }
  }

  &__append-slot {
    position: absolute;
    top: var(--v-field-padding-top, 0);
    right: var(--v-field-padding-end, 16px);
    z-index: 2;
    pointer-events: auto;
  }
}
</style>
