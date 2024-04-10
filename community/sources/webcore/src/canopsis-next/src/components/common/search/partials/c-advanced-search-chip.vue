<template>
  <div>
    <input
      v-if="active && isTypeValue"
      :value="item.value"
      type="text"
      autocomplete="off"
    >
    <v-chip
      v-else
      :outlined="isTypeValue"
      :input-value="active"
      @mousedown.stop=""
      @mouseup.stop=""
      @click="clickChip"
    >
      <span v-if="item.not" class="mr-2 font-italic">{{ $t('advancedSearch.not') }}</span>
      <span>{{ item.text }}</span>
    </v-chip>
  </div>
</template>

<script>
import { toRef } from 'vue';

import { useType } from '../hooks/advanced-search';

export default {
  props: {
    item: {
      type: Object,
      required: true,
    },
    active: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { isTypeValue } = useType(toRef(props.item, 'type'));

    const clickChip = () => emit('select:item', props.item);

    return {
      isTypeValue,

      clickChip,
    };
  },
};
</script>
