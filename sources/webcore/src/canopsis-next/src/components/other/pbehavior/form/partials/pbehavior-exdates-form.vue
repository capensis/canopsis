<template lang="pug">
  div
    v-layout(row)
      strong {{ $t('modals.createPbehavior.fields.exdate') }}
    v-layout(v-for="(exdate, index) in exdates", :key="exdate.key", row)
      v-flex
        date-time-picker-field(
          :value="exdate.value",
          useSeconds,
          @input="updateFieldInArrayItem(index, 'value', $event)"
        )
        v-btn(color="error", icon, @click="removeItemFromArray(index)")
          v-icon delete
    v-btn.primary.ml-0(@click="addItem") {{ $t('modals.createPbehavior.buttons.addExdate') }}
</template>

<script>
import uid from '@/helpers/uid';

import formArrayMixin from '@/mixins/form/array';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerField },
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
