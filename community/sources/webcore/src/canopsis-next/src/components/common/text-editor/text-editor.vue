<template>
  <div class="text-editor">
    <v-label v-show="label">
      {{ label }}
    </v-label>
    <div
      :class="{ 'error--text': hasError }"
      class="text-editor"
      @blur="$emit('blur', $event)"
    >
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
        @input="pasteVariable"
        @close="closeVariablesMenu"
      />
      <div class="text-editor__details">
        <v-messages
          :value="errorMessages"
          color="error"
        />
      </div>
    </div>
  </div>
</template>

<script>
import {
  ref,
  computed,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';
import { Jodit } from 'jodit';

import 'jodit/esm/plugins/all';

import 'jodit/es2021/jodit.min.css';

import { BASE_URL, DEFAULT_SANITIZE_OPTIONS } from '@/config';

import { sanitizeHtml } from '@/helpers/html';

import { useI18n } from '@/hooks/i18n';

import { useTextEditorUpload, useTextEditorVariables } from './hooks/text-editor';
import VariablesMenu from './variables-menu.vue';

export default {
  components: { VariablesMenu },
  props: {
    value: {
      type: String,
    },
    label: {
      type: String,
      default: '',
    },
    buttons: {
      type: Array,
      default: () => [
        'source', '|',
        'bold', 'italic', 'strikethrough', 'underline', '|',
        'ul', 'ol', '|',
        'font', 'fontsize', 'brush', 'paragraph', '|',
        'image', 'table', 'link', '|',
        'align', 'undo', 'redo', '|',
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
    sanitizeOptions: {
      type: Object,
      required: false,
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { locale } = useI18n();

    const editor = ref(null);
    const textEditorElement = ref(null);
    const sanitized = ref(false);

    const sanitizedValue = computed(() => sanitizeHtml(props.value, props.sanitizeOptions || DEFAULT_SANITIZE_OPTIONS));
    const hasError = computed(() => props.errorMessages.length);

    const { controlsOptions, uploaderOptions } = useTextEditorUpload({
      editor,
      isPublic: props.public,
      maxFileSize: props.maxFileSize,
    });

    const {
      hasVariables,
      variablesShown,
      variablesMenuValue,
      variablesMenuPosition,
      variablesButton,
      pasteVariable,
      closeVariablesMenu,
    } = useTextEditorVariables({
      editor,
      variables: props.variables,
    });

    const options = computed(() => {
      const config = {
        language: locale.value,
        enter: 'p',
        cleanHTML: {
          fillEmptyParagraph: false,
        },
        toolbarSticky: false,
        addNewLine: false,
        hidePoweredByJodit: true,
        controls: controlsOptions.value,
        uploader: uploaderOptions.value,
        sourceEditor: 'ace',
        sourceEditorCDNUrlsJS: [
          `${BASE_URL}scripts/libs/ace/1.43.3/ace.js`,
        ],
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

    const changeValue = value => emit('input', value);

    const makeEditor = () => {
      editor.value = Jodit.make(textEditorElement.value, options.value);
      editor.value.value = props.value;
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

    watch(() => props.value, (newValue) => {
      if (!editor.value || newValue === editor.value.value) {
        return;
      }

      if (newValue && !sanitized.value) {
        editor.value.value = sanitizedValue.value;
        sanitized.value = true;
      } else {
        editor.value.value = newValue;
      }
    }, { immediate: true });

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
  &__details {
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-flex: 1;
    -ms-flex: 1 0 auto;
    flex: 1 0 auto;
    max-width: 100%;
    overflow: hidden;
  }

  &.error--text {
    .jodit-container {
      margin-bottom: 8px;
      border-color: var(--v-error-base) !important;
    }
  }

  .jodit-progress-bar div {
    background: var(--v-primary-base);
  }
}

.jodit-dialog__panel {
  pointer-events: all;

  .jodit-ui-button_variant_primary {
    background-color: var(--v-primary-base);

    &:hover:not([disabled]) {
      background-color: var(--v-primary-darken1);
    }
  }

  .jodit-dialog__resizer {
    display: none;
  }
}
</style>
