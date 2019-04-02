<template lang="pug">
  div
    h2 Exdate
    v-layout(v-for="(item, index) in value", :key="`exdate-item-${index}`", row)
      v-flex(xs11)
        date-time-picker-field(
        v-validate="'required'",
        :value="item",
        :name="`exdate[${index}]`",
        @input="updateItemInArray(index, $event)"
        )
      v-flex(xs1)
        v-btn(color="error", icon, @click="removeItemFromArray(index)")
          v-icon delete
    v-btn.primary(@click="addItem") Add
</template>

<script>
import formArrayMixin from '@/mixins/form/array';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerField },
  mixins: [formArrayMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray(null);
    },
  },
};
</script>
