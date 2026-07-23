<template>
  <c-layout class="template-item" column>
    <span class="text-subtitle-2 mb-3">
      {{ title }}
    </span>

    <c-layout :class="statusProps.headerClass" class="template-item__header" align="stretch">
      <div class="template-item__header-section">
        <span>{{ t('common.input') }}</span>
      </div>

      <div class="template-item__header-section">
        <c-layout justify="space-between" align="center">
          <c-layout align="center">
            <span>{{ t('common.output') }}</span>

            <v-chip :color="statusProps.chipColor" class="ml-2" size="small">
              {{ statusProps.chipText }}
            </v-chip>
          </c-layout>

          <c-enabled-field
            v-model="isJson"
            :label="t('templateTesting.formatJson')"
            class="mt-0"
            hide-details
          />
        </c-layout>
      </div>
    </c-layout>

    <c-layout class="template-item__content">
      <div class="template-item__input-section">
        <c-payload-textarea-field
          v-if="textarea"
          :value="template"
          :variables="variables"
          :name="name"
          wrap="off"
          rows="2"
          class="template-item__input-editor"
          @update:value="updateTemplate"
        />

        <c-payload-text-field
          v-else
          :value="template"
          :variables="variables"
          :name="name"
          label=""
          class="template-item__input-editor"
          @update:value="updateTemplate"
        />
      </div>

      <div class="template-item__output-section">
        <c-payload-textarea-field
          :value="formattedOutput"
          :line-height="TEMPLATE_TESTING_TEST_VARIABLES_EDITOR_LINE_HEIGHT"
          :style="{ height: `${editorHeight}px` }"
          class="template-item__output-editor"
          rows="2"
          readonly
        />
      </div>
    </c-layout>
  </c-layout>
</template>

<script>
import { isUndefined } from 'lodash';
import { computed, ref, watch } from 'vue';

import {
  TEMPLATE_TESTING_TEST_VARIABLES_EDITOR_LINE_HEIGHT,
  TEMPLATE_TESTING_TEST_VARIABLES_MIN_EDITOR_LINES,
} from '@/constants';

import { getStringLinesCount } from '@/helpers/string';

import { useI18n } from '@/hooks/i18n';

export default {
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
    name: {
      type: String,
      default: '',
    },
    variables: {
      type: Array,
      default: () => [],
    },
    result: {
      type: Object,
      default: () => ({}),
    },
    lastRunValue: {
      type: String,
      default: '',
    },
  },

  emits: ['update:template'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const isJson = ref(false);

    const editorHeight = ref(
      TEMPLATE_TESTING_TEST_VARIABLES_MIN_EDITOR_LINES *
        TEMPLATE_TESTING_TEST_VARIABLES_EDITOR_LINE_HEIGHT,
    );

    const statusProps = computed(() => {
      const { err, is_valid: isValid } = props.result;

      if (props.lastRunValue !== props.template && !isUndefined(isValid)) {
        return {
          headerClass: 'template-item__header--warning',
          chipColor: 'warning',
          chipText: t('templateTesting.notRun'),
        };
      }

      if (err) {
        return {
          headerClass: 'template-item__header--error',
          chipColor: 'error',
          chipText: t('templateTesting.errorsInInput'),
        };
      }

      if (isValid) {
        return {
          headerClass: 'template-item__header--success',
          chipColor: 'success',
          chipText: t('templateTesting.success'),
        };
      }

      return {
        chipColor: 'grey',
        chipText: t('templateTesting.notRun'),
      };
    });

    const formattedOutput = computed(() => {
      if (isJson.value) {
        try {
          return JSON.stringify(JSON.parse(props.result.result), null, 2);
        } catch {
          return props.result.result;
        }
      }

      return props.result.result;
    });

    const updateTemplate = (value) => emit('update:template', value);

    watch(formattedOutput, (value = '') => {
      editorHeight.value = Math.max(
        getStringLinesCount(value) * TEMPLATE_TESTING_TEST_VARIABLES_EDITOR_LINE_HEIGHT,
        TEMPLATE_TESTING_TEST_VARIABLES_MIN_EDITOR_LINES *
          TEMPLATE_TESTING_TEST_VARIABLES_EDITOR_LINE_HEIGHT,
      );
    });

    return {
      t,
      statusProps,
      editorHeight,
      isJson,
      formattedOutput,
      TEMPLATE_TESTING_TEST_VARIABLES_EDITOR_LINE_HEIGHT,
      updateTemplate,
    };
  },
};
</script>

<style lang="scss" scoped>
.template-item {
  --paddings: 8px;
  --border-color: rgba(var(--v-border-color), var(--v-border-opacity));
  --header-background-color: rgb(var(--v-theme-secondary-lighten-2));

  &__header {
    background-color: var(--header-background-color);
    border: 1px solid var(--border-color);

    &--success {
      background-color: rgb(var(--v-theme-success-background)) !important;
    }

    &--warning {
      background-color: rgb(var(--v-theme-warning-background)) !important;
    }

    &--error {
      background-color: rgb(var(--v-theme-error-background)) !important;
    }

    &-section {
      flex: 0 0 50%;
      max-width: 50%;
      padding: var(--paddings);

      &:first-child {
        border-right: 1px solid var(--border-color);
      }
    }
  }

  &__content {
    .template-item__input-section,
    .template-item__output-section {
      flex: 0 0 50%;
      max-width: 50%;
      padding: var(--paddings);
      border-right: 1px solid var(--border-color);
      border-bottom: 1px solid var(--border-color);
    }

    .template-item__input-section {
      border-left: 1px solid var(--border-color);
    }
  }

  &__input-editor,
  &__output-editor {
    min-height: 40px;
  }
}
</style>
