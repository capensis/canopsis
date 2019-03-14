<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") Sort Order
    v-container
      v-select(
      :value="value",
      @input="$emit('input', $event)",
      :items="orders"
      )
</template>

<script>
import formMixin from '@/mixins/form';

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
      default: 'ASC',
    },
  },
  data() {
    return {
      orders: ['ASC', 'DESC'],
    };
  },
};
</script>
