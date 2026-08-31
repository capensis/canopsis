<template>
  <div class="text-editor">
    <c-label
      v-if="label"
      :label="label"
      :error="hasError"
      :required="required"
      class="text-editor__label text-caption"
    />

    <div :class="{ 'error--text': hasError }" class="text-editor" @blur="$emit('blur', $event)">
      <div ref="textEditorElement" />

      <variables-menu
        v-if="hasVariables"
        :items="variables"
        :visible="variablesShown"
        :value="variablesMenuValue"
        :position-x="variablesMenuPosition.x"
        :position-y="variablesMenuPosition.y"
        dense
        clickable-parent
        @update:value="pasteVariable"
        @close="closeVariablesMenu"
      />

      <div class="text-editor__details">
        <v-messages :messages="errorMessages" color="error" />
      </div>
    </div>
  </div>
</template>

<script>
import { Jodit } from 'jodit';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import 'jodit/esm/plugins/all';
import 'jodit/es2021/jodit.min.css';

import { BASE_URL, DEFAULT_SANITIZE_OPTIONS } from '@/config';

import { sanitizeHtml } from '@/helpers/html';

import { useI18n } from '@/hooks/i18n';

import { useTextEditorUpload, useTextEditorVariables } from './hooks/text-editor';
import VariablesMenu from './variables-menu.vue';

const JODIT_TO_MATERIAL_ICONS_MAP = {
  source: 'code',
  bold: 'format_bold',
  italic: 'format_italic',
  strikethrough: 'strikethrough_s',
  underline: 'format_underlined',
  ul: 'format_list_bulleted',
  ol: 'format_list_numbered',
  font: 'font_download',
  fontsize: 'text_fields',
  brush: 'format_paint',
  paragraph: 'notes',
  image: 'image',
  table: 'table_chart',
  link: 'link',
  left: 'format_align_left',
  center: 'format_align_center',
  right: 'format_align_right',
  justify: 'format_align_justify',
  undo: 'undo',
  redo: 'redo',
};

export default {
  components: { VariablesMenu },

  props: {
    modelValue: {
      type: String,
      default: undefined,
    },
    value: {
      type: String,
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    buttons: {
      type: Array,
      default: () => [
        'source',
        '|',
        'bold',
        'italic',
        'strikethrough',
        'underline',
        '|',
        'ul',
        'ol',
        '|',
        'font',
        'fontsize',
        'brush',
        'paragraph',
        '|',
        'image',
        'table',
        'link',
        '|',
        'align',
        'undo',
        'redo',
        '|',
      ],
    },
    public: {
      type: Boolean,
      default: false,
    },
    extraButtons: {
      type: Array,
      default: () => [],
    },
    config: {
      type: Object,
      default: () => ({}),
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    maxFileSize: {
      type: Number,
      required: false,
    },
    variables: {
      type: Array,
      required: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    sanitizeOptions: {
      type: Object,
      required: false,
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'blur'],

  setup(props, { emit }) {
    const { locale } = useI18n();

    const editor = ref(null);
    const textEditorElement = ref(null);
    const sanitized = ref(false);

    const currentValue = computed(() => props.modelValue ?? props.value ?? '');
    const sanitizedValue = computed(
      () => sanitizeHtml(currentValue.value, props.sanitizeOptions || DEFAULT_SANITIZE_OPTIONS),
    );
    const hasError = computed(() => props.errorMessages.length > 0);

    const { controlsOptions, uploaderOptions } = useTextEditorUpload({
      editor,
      isPublic: computed(() => props.public),
      maxFileSize: computed(() => props.maxFileSize),
    });

    const {
      hasVariables,
      variablesShown,
      variablesMenuValue,
      variablesMenuPosition,
      variablesButton,
      variablesExtraIcon,
      pasteVariable,
      closeVariablesMenu,
    } = useTextEditorVariables({
      editor,
      variables: computed(() => props.variables),
    });

    const options = computed(() => {
      const config = {
        language: locale.value,
        enter: 'p',
        cleanHTML: {
          fillEmptyParagraph: false,
          replaceNBSP: false,
        },
        toolbarSticky: false,
        addNewLine: false,
        showCharsCounter: false,
        showWordsCounter: false,
        showXPathInStatusbar: false,
        hidePoweredByJodit: true,
        controls: controlsOptions.value,
        uploader: uploaderOptions.value,
        sourceEditor: 'ace',
        extraIcons: { ...variablesExtraIcon.value },
        getIcon(name) {
          if (JODIT_TO_MATERIAL_ICONS_MAP[name]) {
            return `<i class="material-icons v-icon v-icon--small" style="font-size: 18px;">${JODIT_TO_MATERIAL_ICONS_MAP[name]}</i>`;
          }

          return '';
        },
        sourceEditorCDNUrlsJS: [`${BASE_URL}scripts/libs/ace/1.43.3/ace.js`],
        ...props.config,
      };

      if (props.buttons.length) {
        config.buttons = props.buttons;
        config.buttonsMD = props.buttons;
        config.buttonsSM = props.buttons;
        config.buttonsXS = props.buttons;
      }

      if (props.dark && !props.config.theme) {
        config.theme = 'dark';
      }

      config.extraButtons = [];

      if (hasVariables.value) {
        config.extraButtons.push(variablesButton.value);
      }

      if (props.extraButtons.length) {
        config.extraButtons.push(...props.extraButtons);
      }

      if (props.autofocus) {
        config.autofocus = props.autofocus;
        config.cursorAfterAutofocus = 'start';
      }

      return config;
    });

    const changeValue = (value) => {
      emit('update:modelValue', value);
    };

    const makeEditor = () => {
      editor.value = Jodit.make(textEditorElement.value, options.value);
      editor.value.value = currentValue.value;
      editor.value.options.popupRoot = editor.value.container;
      editor.value.events.on('change', changeValue);
    };

    const destructEditor = () => {
      if (editor.value) {
        editor.value.events.off('change', changeValue);
        editor.value.destruct();
        editor.value = null;
      }
    };

    watch(options, () => {
      destructEditor();
      makeEditor();
    });

    watch(
      () => currentValue.value,
      (newValue) => {
        if (!editor.value || newValue === editor.value.value) {
          return;
        }

        if (newValue && !sanitized.value) {
          editor.value.value = sanitizedValue.value;
          sanitized.value = true;
        } else {
          editor.value.value = newValue;
        }
      },
      { immediate: true },
    );

    onMounted(makeEditor);
    onBeforeUnmount(destructEditor);

    return {
      editor,
      textEditorElement,
      variablesShown,
      variablesMenuValue,
      variablesMenuPosition,
      hasVariables,
      hasError,
      pasteVariable,
      closeVariablesMenu,
    };
  },
};
</script>

<style lang="scss">
.text-editor {
  --jd-font-default: 'Roboto', sans-serif;

  width: 100%;

  &__label {
    display: block;
    margin-bottom: 8px;
  }

  &__details {
    display: flex;
    flex: 1 0 auto;
    max-width: 100%;
    overflow: hidden;
  }

  &.error--text {
    .jodit-container {
      margin-bottom: 8px;
      border-color: rgb(var(--v-theme-error)) !important;
    }
  }

  .jodit-container {
    color: rgba(var(--v-theme-on-surface), var(--v-high-emphasis-opacity));
    border-color: rgba(var(--v-border-color), var(--v-border-opacity));

    &:focus-within {
      border-color: rgb(var(--v-theme-primary));
    }
  }

  .jodit-progress-bar div {
    background: rgb(var(--v-theme-primary));
  }

  .jodit-icon {
    color: var(--jd-color-icon);
  }

  .jodit-ui-button_variant_primary {
    background-color: rgb(var(--v-theme-primary));
    text-transform: uppercase;

    .jodit-ui-button__text {
      color: #ffffff !important;
    }

    &:hover:not([disabled]) {
      background-color: rgb(var(--v-theme-primary-darken-1));

      .jodit-ui-button__text {
        background-color: unset;
      }
    }
  }
}

.jodit-dialog__panel {
  pointer-events: all;

  .jodit-ui-button_variant_primary {
    background-color: rgb(var(--v-theme-primary));

    &:hover:not([disabled]) {
      background-color: rgb(var(--v-theme-primary-darken-1));
    }
  }

  .jodit-dialog__resizer {
    display: none;
  }
}
</style>
