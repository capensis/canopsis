<template>
  <input
    v-if="active && isItemTypeValue"
    ref="inputElement"
    :value="item.value"
    type="text"
    autocomplete="off"
    @keydown="keydownInput"
    @input="input"
  >
  <v-chip
    v-else
    :outlined="isItemTypeValue"
    :input-value="active"
    :close="last"
    class="c-advanced-search-chip"
    @mousedown.stop=""
    @mouseup.stop=""
    @click="clickChip"
    @click:close="clickChipClose"
  >
    <span v-if="item.not" class="mr-2 font-italic">{{ $t('advancedSearch.not') }}</span>
    <span>{{ item.text }}</span>
  </v-chip>
</template>

<script>
import { ref, toRef, watch, nextTick } from 'vue';

import { useAdvancedSearchItemType } from '../hooks/advanced-search';

export default {
  model: {
    prop: 'item',
    event: 'input',
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
    active: {
      type: Boolean,
      default: false,
    },
    last: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const inputElement = ref();

    const { isItemTypeValue } = useAdvancedSearchItemType({ type: toRef(props.item, 'type') });

    watch(() => props.active, (active) => {
      if (active) {
        nextTick(() => inputElement.value?.focus());
      }
    });

    const input = ({ target: { value } = {} } = {}) => emit('input', { ...props.item, value, text: value });
    const clickChip = () => emit('click:item', props.item);
    const clickChipClose = () => emit('remove:item', props.item);
    const keydownInput = event => emit('keydown:input', event);

    return {
      inputElement,

      isItemTypeValue,

      input,
      clickChip,
      clickChipClose,
      keydownInput,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-advanced-search-chip {
  border-radius: 6px;
  height: 28px;
  margin: 2px;

  &.theme--light.v-chip--active::before {
    opacity: .25;
  }
}
</style>
