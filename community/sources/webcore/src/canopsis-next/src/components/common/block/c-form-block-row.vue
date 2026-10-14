<template>
  <c-layout :class="layoutClasses" class="c-form-block-row" align="start">
    <div :style="labelStyle" class="c-form-block-row__label text-break">
      <h4 class="c-form-block-row__title text-subtitle-1 text-break">
        {{ label }}
      </h4>
    </div>

    <div :class="fieldClasses" class="text-break">
      <slot />
    </div>
  </c-layout>
</template>

<script>
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import { useValidationChildren } from '@/hooks/validator/validation-children';

export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    width: {
      type: [Number, String],
      default: '28%',
    },
    depth: {
      type: [Number, String],
      default: 0,
    },
    indented: {
      type: Boolean,
      default: false,
    },
    topBorder: {
      type: Boolean,
      default: false,
    },
    bottomBorder: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
    alignCenter: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const theme = useTheme();
    const { hasChildrenError } = useValidationChildren();
    const depth = computed(() => Number(props.depth) || 0);

    const layoutClasses = computed(() => ({
      [`c-form-block-row--depth-${props.depth}`]: !!depth.value,
      'c-form-block-row--depth-more': depth.value > 2,
      'c-form-block-row--top-border': !!props.topBorder,
      'c-form-block-row--bottom-border': !!props.bottomBorder,
      'c-form-block-row--error': hasChildrenError.value || props.error,
      'c-form-block-row--dark': theme.current.value.dark,
      'c-form-block-row--light': !theme.current.value.dark,
    }));

    const fieldClasses = computed(() => ({
      'c-form-block-row__field': true,
      'py-3': props.indented,
      'c-form-block-row__field--align-center': props.alignCenter,
    }));

    const labelWidth = computed(() => Number.isFinite(props.width) ? `${props.width}px` : props.width,
    );

    const labelStyle = computed(() => ({
      flex: `0 0 ${labelWidth.value}`,
      minWidth: labelWidth.value,
      maxWidth: labelWidth.value,
    }));

    return {
      layoutClasses,
      fieldClasses,
      labelStyle,
    };
  },
};
</script>

<style lang="scss">
:root {
  --c-form-block-row-label-background-color-light: #eaeaea;
  --c-form-block-row-label-background-color-light-depth-1: #f4f4f4;
  --c-form-block-row-label-background-color-light-depth-2: #fafafa;
  --c-form-block-row-label-background-color-light-depth-more: #ffffff;

  --c-form-block-row-label-background-color-dark: var(--v-application-background-lighten2);
  --c-form-block-row-label-background-color-dark-depth-1: #404040;
  --c-form-block-row-label-background-color-dark-depth-2: #505050;
  --c-form-block-row-label-background-color-dark-depth-more: #606060;
}

.c-form-block-row {
  &:not(:last-child) {
    border-bottom: 1px solid var(--v-application-background-darken2);
  }

  &--top-border {
    border-top: 1px solid var(--v-application-background-darken2);
  }

  &--bottom-border {
    border-bottom: 1px solid var(--v-application-background-darken2);
  }

  .c-form-block-row__label {
    padding-left: 24px;
  }

  &--depth-1 .c-form-block-row__label {
    padding-left: 40px;
  }

  &--depth-2 .c-form-block-row__label {
    padding-left: 56px;
  }

  &--depth-more .c-form-block-row__label {
    padding-left: 72px;
  }

  &__label {
    align-self: stretch;
    display: flex;
    align-items: flex-start;
    border-right: 1px solid var(--v-application-background-darken2);
    padding-top: 28px;
    padding-right: 24px;
    padding-bottom: 28px;
  }

  &__title {
    margin: 0;
    font-weight: 400;
    word-break: break-word;
    overflow-wrap: anywhere;
    max-width: 100%;
  }

  &__field {
    flex: 1 1 auto !important;
    min-width: 0;
    align-self: stretch;
    word-break: break-word;
    overflow-wrap: anywhere;
    padding: 4px 16px;

    &--align-center {
      display: flex;
      align-items: center;
    }
  }

  &--error {
    .c-form-block-row__label {
      background-color: var(--v-error-background-base) !important;
    }
  }
}

.theme--light .c-form-block-row,
.c-form-block-row--light {
  .c-form-block-row__label {
    background-color: var(--c-form-block-row-label-background-color-light);
  }
}

.theme--light .c-form-block-row--depth-1 .c-form-block-row__label,
.c-form-block-row--light.c-form-block-row--depth-1 .c-form-block-row__label {
  background-color: var(--c-form-block-row-label-background-color-light-depth-1);
}

.theme--light .c-form-block-row--depth-2 .c-form-block-row__label,
.c-form-block-row--light.c-form-block-row--depth-2 .c-form-block-row__label {
  background-color: var(--c-form-block-row-label-background-color-light-depth-2);
}

.theme--light .c-form-block-row--depth-more .c-form-block-row__label,
.c-form-block-row--light.c-form-block-row--depth-more .c-form-block-row__label {
  background-color: var(--c-form-block-row-label-background-color-light-depth-more);
}

.theme--dark .c-form-block-row,
.c-form-block-row--dark {
  &:not(:last-child) {
    border-color: var(--v-application-background-lighten3);
  }

  .c-form-block-row__label {
    background-color: var(--c-form-block-row-label-background-color-dark);
    border-color: var(--v-application-background-lighten3);
  }

  .c-form-block-row__title {
    color: var(--v-active-color-base);
  }
}

.theme--dark .c-form-block-row--depth-1 .c-form-block-row__label,
.c-form-block-row--dark.c-form-block-row--depth-1 .c-form-block-row__label {
  background-color: var(--c-form-block-row-label-background-color-dark-depth-1);
}

.theme--dark .c-form-block-row--depth-2 .c-form-block-row__label,
.c-form-block-row--dark.c-form-block-row--depth-2 .c-form-block-row__label {
  background-color: var(--c-form-block-row-label-background-color-dark-depth-2);
}

.theme--dark .c-form-block-row--depth-more .c-form-block-row__label,
.c-form-block-row--dark.c-form-block-row--depth-more .c-form-block-row__label {
  background-color: var(--c-form-block-row-label-background-color-dark-depth-more);
}
</style>
