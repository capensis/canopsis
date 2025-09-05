<template>
  <v-layout class="template-item" column>
    <span class="text-subtitle-2 mb-3">
      {{ title }}
    </span>
    <v-layout :class="statusProps.headerClass" class="template-item__header" align-stretch>
      <v-flex class="template-item__header-section" xs6>
        <span>{{ $t('common.input') }}</span>
      </v-flex>
      <v-flex class="template-item__header-section" xs6>
        <v-layout justify-space-between align-center>
          <v-layout align-center>
            <span>{{ $t('common.output') }}</span>
            <v-chip
              :color="statusProps.chipColor"
              class="ml-2"
              text-color="white"
              small
            >
              {{ statusProps.chipText }}
            </v-chip>
          </v-layout>
          <c-enabled-field
            v-model="isJson"
            :label="$t('templateTesting.formatJson')"
            class="mt-0"
            hide-details
          />
        </v-layout>
      </v-flex>
    </v-layout>
    <v-layout class="template-item__content">
      <v-flex class="template-item__input-section" xs6>
        <c-payload-textarea-field
          v-if="textarea"
          v-field="template"
          :placeholder="$t('templateTesting.inputPlaceholder')"
          wrap="off"
          rows="2"
          class="template-item__input-editor"
        />
        <c-payload-text-field
          v-else
          v-field="template"
          :placeholder="$t('templateTesting.inputPlaceholder')"
          class="template-item__input-editor"
        />
      </v-flex>
      <v-flex class="template-item__output-section" xs6>
        <simple-code-editor
          :value="formattedOutput"
          :language="isJson ? 'json' : 'plaintext'"
          :style="{ height: `${editorHeight}px` }"
          class="template-item__output-editor"
          readonly
        />
      </v-flex>
    </v-layout>
  </v-layout>
</template>

<script>
import { ref, computed, watch } from 'vue';

import { getStringLinesCount } from '@/helpers/string';

import { useI18n } from '@/hooks/i18n';

import SimpleCodeEditor from '@/components/common/code-editor/simple-code-editor.vue';

const EDITOR_LINE_HEIGHT = 19;
const MIN_EDITOR_LINES = 2;

export default {
  components: { SimpleCodeEditor },
  model: {
    prop: 'template',
    event: 'input',
  },
  props: {
    template: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
    textarea: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t, tc } = useI18n();

    const isRunning = ref(false);
    const isJson = ref(false);

    const editorHeight = ref(MIN_EDITOR_LINES * EDITOR_LINE_HEIGHT);

    const errorsCount = 4;

    const statusProps = computed(() => {
      if (true) { // SUCCESS
        return {
          headerClass: 'template-item__header--success',
          chipColor: 'success',
          chipText: t('templateTesting.success'),
        };
      }

      if (false) { // WARNING
        return {
          headerClass: 'template-item__header--warning',
          chipColor: 'warning',
          chipText: t('templateTesting.notRun'),
        };
      }

      if (false) { // ERROR
        return {
          headerClass: 'template-item__header--error',
          chipColor: 'error',
          chipText: tc('templateTesting.errorsInInput', errorsCount, { count: errorsCount }),
        };
      }

      return {
        chipColor: 'grey',
        chipText: t('templateTesting.notRun'),
      };
    });

    const output = '{"test": 321, "asd": "asdasd", "arr": [1, 2, 3]}';

    const formattedOutput = computed(() => {
      if (!output) {
        return 'asdasdasdasdasdasd {} {{ .asdasd.asd }} []\ndslakdlsadk';
      }

      if (isJson.value) {
        try {
          return JSON.stringify(JSON.parse(output), null, 2);
        } catch {
          return output;
        }
      }

      return output;
    });

    watch(formattedOutput, (value) => {
      editorHeight.value = getStringLinesCount(value) * EDITOR_LINE_HEIGHT;
    });

    return {
      statusProps,
      editorHeight,
      isJson,
      formattedOutput,
      isRunning,
    };
  },
};
</script>

<style lang="scss" scoped>
.template-item {
  --paddings: 8px;

  --border-color: var(--v-application-background-darken2);
  --border-color-dark: var(--v-application-background-lighten3);

  --header-background-color: var(--v-application-background-darken1);
  --header-background-color-success: var(--v-success-background-base);
  --header-background-color-warning: var(--v-warning-background-base);
  --header-background-color-error: var(--v-error-background-base);

  .template-item__header {
    background-color: var(--header-background-color);
    border: 1px solid var(--border-color);

    &--success {
      background-color: var(--header-background-color-success) !important;
    }

    &--warning {
      background-color: var(--header-background-color-warning) !important;
    }

    &--error {
      background-color: var(--header-background-color-error) !important;
    }

    &-section {
      padding: var(--paddings);

      &:first-child {
        border-right: 1px solid var(--border-color);
      }
    }
  }

  .template-item__content {
    .template-item__input-section,
    .template-item__output-section {
      padding: var(--paddings);
      border-right: 1px solid var(--border-color);
      border-bottom: 1px solid var(--border-color);
    }

    .template-item__input-section {
      border-left: 1px solid var(--border-color);

      & ::v-deep .v-text-field__slot {
        max-width: 100%;

        textarea {
          overflow-x: auto;
        }
      }
    }

    .template-item__input-editor {
      &, & ::v-deep textarea {
        min-height: 40px;
      }
    }

    .template-item__output-editor {
      min-height: 100%;
    }
  }
}

.theme--dark {
  .template-item {
    border-color: var(--border-color-dark);

    .template-item__header {
      background-color: var(--header-background-color-dark);
      border-bottom-color: var(--border-color-dark);

      &-section:first-child {
        border-right-color: var(--border-color-dark);
      }
    }

    .template-item__input-section {
      border-right-color: var(--border-color-dark);
    }
  }
}
</style>
