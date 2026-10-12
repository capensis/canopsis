<template>
  <span class="c-action-btn">
    <c-tooltip :location="location" :disabled="disabled">
      <template #activator="{ props: activatorProps }">
        <slot :props="activatorProps" name="button">
          <div class="c-action-btn__button-wrapper" v-bind="{ ...activatorProps, ...$attrs }">
            <v-btn
              :disabled="isButtonDisabled"
              :loading="loading"
              :size="small ? 'small' : undefined"
              :color="btnColor"
              :theme="dark ? 'dark' : undefined"
              :active="inputValue"
              class="c-action-btn__button"
              variant="text"
              icon
              @click.stop.prevent="$emit('click', $event)"
            >
              <c-icon :color="preparedProps.color" :size="iconSize" :small="iconSmall">
                {{ preparedProps.icon }}
              </c-icon>
            </v-btn>
          </div>
        </slot>
      </template>

      <span>{{ preparedProps.tooltip }}</span>
    </c-tooltip>
  </span>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  inheritAttrs: false,

  props: {
    type: {
      type: String,
      default: null,
      validator: value => ['edit', 'duplicate', 'delete'].includes(value),
    },
    icon: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: '',
    },
    tooltip: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledButton: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    top: {
      type: Boolean,
      required: false,
    },
    right: {
      type: Boolean,
      required: false,
    },
    left: {
      type: Boolean,
      required: false,
    },
    bottom: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      required: false,
    },
    btnColor: {
      type: String,
      required: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    inputValue: {
      type: Boolean,
      required: false,
    },
    iconSmall: {
      type: Boolean,
      default: false,
    },
    iconSize: {
      type: [Number, String],
      default: undefined,
    },
  },

  emits: ['click'],

  setup(props) {
    const { t } = useI18n();

    const location = computed(
      () => ({
          [props.top]: 'top',
          [props.right]: 'right',
          [props.left]: 'left',
          [props.bottom || (!props.top && !props.right && !props.left)]: 'bottom',
        }).true,
    );

    const isButtonDisabled = computed(() => props.disabled || props.disabledButton);

    const preparedProps = computed(() => {
      const typesPresetsMap = {
        edit: {
          tooltip: t('common.edit'),
          icon: 'edit',
        },
        duplicate: {
          tooltip: t('common.duplicate'),
          icon: 'file_copy',
        },
        delete: {
          tooltip: t('common.delete'),
          icon: 'delete',
          color: 'error',
        },
      };

      const preset = typesPresetsMap[props.type] || {};

      return {
        icon: props.icon || preset.icon,
        color: props.color || preset.color,
        tooltip: props.tooltip || preset.tooltip,
      };
    });

    return {
      isButtonDisabled,
      location,
      preparedProps,
    };
  },
};
</script>

<style lang="scss">
.c-action-btn {
  display: inline-flex;
  flex: 0 0 auto;
  justify-self: start;
  width: fit-content;
  height: fit-content;
  border-radius: 50%;

  &__button-wrapper {
    display: inline-flex;
    flex: 0 0 auto;
    align-items: center;
    justify-content: center;
    width: fit-content;
    height: fit-content;
    border-radius: 50%;
  }

  &__button {
    flex: 0 0 auto;
  }
}
</style>
