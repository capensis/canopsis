<template>
  <input
    :value="value"
    :placeholder="placeholder"
    type="text"
    class="internal-search"
    autocomplete="off"
    @input="input"
    @blur="blur"
    @click="click"
    @keydown="keydown"
  >
</template>

<script>
import { KEY_CODES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    isValueType: {
      type: Boolean,
      default: false,
    },
    placeholder: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const input = event => updateModel(event.target.value);
    const click = () => emit('click');
    const apply = () => emit('apply');
    const blur = () => {
      if (props.isValueType && props.value) {
        apply();
      }
    };

    const keydown = (event) => {
      if (event.keyCode === KEY_CODES.backspace && !props.value) {
        emit('keydown:backspace', event);
        return;
      }

      if (
        [
          KEY_CODES.enter,
          KEY_CODES.esc,
          KEY_CODES.down,
          KEY_CODES.up,
          KEY_CODES.home,
          KEY_CODES.end,
        ].includes(event.keyCode)
      ) {
        emit('keydown:navigate', event);
      }

      if (event.keyCode === KEY_CODES.esc) {
        updateModel('');
        return;
      }

      if (props.isValueType && event.keyCode === KEY_CODES.enter) {
        apply();
      }
    };

    return {
      input,
      blur,
      click,
      keydown,
    };
  },
};
</script>

<style lang="scss" scoped>
.internal-search::placeholder {
  .theme--light & {
    color: var(--v-text-light-secondary, rgba(0, 0, 0, 0.6)) !important;
  }

  .theme--dark & {
    color: var(--v-text-dark-secondary, rgba(255, 255, 255, 0.7)) !important;
  }
}
</style>
