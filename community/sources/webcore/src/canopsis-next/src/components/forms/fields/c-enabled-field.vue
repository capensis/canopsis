<template>
  <div :class="wrapperClass" class="c-enabled-field">
    <v-switch
      v-field="value"
      :class="{ 'ma-0': withBackground }"
      :label="label || $t('common.enabled')"
      :color="color"
      :disabled="disabled"
      :readonly="readonly"
      :hide-details="hideDetails || withBackground"
      v-on="listenersWithoutInput"
    >
      <template #label="">
        <slot name="label" />
      </template>
      <template #append="">
        <slot name="append" />
      </template>
    </v-switch>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { computed } from 'vue';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      default: true,
    },
    label: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    withBackground: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { listeners }) {
    const listenersWithoutInput = computed(() => (listeners ? omit(listeners, ['input']) : {}));
    const wrapperClass = computed(() => (
      props.withBackground
        ? {
          'c-enabled-field--with-background': true,
          'c-enabled-field--with-background--disabled': !props.value,
        }
        : {}
    ));

    return {
      listenersWithoutInput,
      wrapperClass,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-enabled-field {
  display: flex;
}

.c-enabled-field--with-background {
  padding: 8px;
  border-radius: 5px;
  border: 1px solid var(--v-success-base);
  background-color: var(--v-success-lighten3);

  ::v-deep .v-input {
    margin: 0;
    padding: 0;
  }

  &--disabled {
    border-color: var(--v-error-base);
    background-color: var(--v-error-lighten4);
  }
}
</style>
