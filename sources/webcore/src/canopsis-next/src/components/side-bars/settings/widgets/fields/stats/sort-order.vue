<template lang="pug">
  v-list-group(data-test="sortOrder")
    v-list-tile(slot="activator") {{ $t('settings.statsNumbers.sortOrder') }}
    v-container(data-test="sortOrderSelectLayout")
      v-select(
        :value="value",
        @input="$emit('input', $event)",
        :items="orders"
      )
</template>

<script>
import formMixin from '@/mixins/form';

import { SORT_ORDERS } from '@/constants';

/**
 * Component to select the default sort order on settings
 *
 * @prop {String} [value] - Sort order value
 *
 * @event value#input
 */
export default {
  mixins: [formMixin],
  props: {
    value: {
      type: String,
      default: SORT_ORDERS.asc,
    },
  },
  computed: {
    orders() {
      return Object.values(SORT_ORDERS);
    },
  },
};
</script>
