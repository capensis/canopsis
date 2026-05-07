<template>
  <v-layout
    :class="layoutClasses"
    class="c-form-block-row"
    align-start
  >
    <v-flex
      :style="labelFlexStyle"
      class="c-form-block-row__label text-break"
    >
      <h4 class="c-form-block-row__title text-subtitle-1 text-break">
        {{ label }}
      </h4>
    </v-flex>
    <v-flex class="c-form-block-row__field text-break">
      <slot />
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

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
  },
  setup(props) {
    const { hasChildrenError } = useValidationChildren();

    const layoutClasses = computed(() => ({
      [`c-form-block-row--depth-${props.depth}`]: !!props.depth,
      'c-form-block-row--error': hasChildrenError.value,
    }));

    const labelFlexStyle = computed(() => {
      if (typeof props.width === 'number') {
        const px = `${props.width}px`;

        return {
          flex: `0 0 ${px}`,
          maxWidth: px,
          minWidth: px,
        };
      }

      return {
        flex: `0 0 ${props.width}`,
        minWidth: props.width,
      };
    });

    return {
      layoutClasses,
      labelFlexStyle,
    };
  },
};
</script>

<style lang="scss">
.c-form-block-row {
  &:not(:last-child) {
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
    font-weight: 600;
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
  }

  &--error {
    .c-form-block-row__label {
      background-color: var(--v-error-background-base) !important;
    }
  }
}

.theme--dark .c-form-block-row {
  &:not(:last-child) {
    border-bottom-color: var(--v-application-background-lighten3);
  }

  .c-form-block-row__label {
    background-color: var(--v-application-background-lighten2);
    border-right-color: var(--v-application-background-lighten3);
  }

  .c-form-block-row__title {
    color: var(--v-active-color-base);
  }

  &--depth-1 .c-form-block-row__label {
    background-image: linear-gradient(rgba(255, 255, 255, 0.052), rgba(255, 255, 255, 0.052));
  }

  &--depth-2 .c-form-block-row__label {
    background-image: linear-gradient(rgba(255, 255, 255, 0.085), rgba(255, 255, 255, 0.085));
  }
}

.theme--light .c-form-block-row {
  .c-form-block-row__label {
    background-color: #EAEAEA;
  }

  &--depth-1 .c-form-block-row__label {
    background-color: #F0F0F0;
  }

  &--depth-2 .c-form-block-row__label {
    background-color: #F5F5F5;
  }
}
</style>
