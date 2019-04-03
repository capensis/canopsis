<template lang="pug">
  div
    v-layout(row)
      strong Exdates
    v-layout(v-for="(exdate, index) in exdates", :key="exdate.key", row)
      v-flex
        date-picker-field(
        :value="exdate.value",
        @input="updateFieldInArrayItem(index, 'value', $event)"
        )
        v-btn(color="error", icon, @click="removeItemFromArray(index)")
          v-icon delete
    v-btn.primary.ml-0(@click="addItem") Add exdate
</template>

<script>
import uid from '@/helpers/uid';

import formArrayMixin from '@/mixins/form/array';

import DatePickerField from '@/components/forms/fields/date-picker/date-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DatePickerField },
  mixins: [formArrayMixin],
  model: {
    prop: 'exdates',
    event: 'input',
  },
  props: {
    exdates: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray({
        key: uid(),
        value: null,
      });
    },
  },
};
</script>
