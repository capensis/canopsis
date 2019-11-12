<template lang="pug">
  div
    v-divider
    h3.my-3.grey--text {{ $t('modals.createPbehavior.steps.rrule.exdate') }}
    v-divider
    v-layout(v-for="(exdate, index) in exdates", :key="exdate.key", row)
      v-flex
        date-time-picker-field(
          :value="exdate.value",
          useSeconds,
          @input="updateFieldInArrayItem(index, 'value', $event)"
        )
        v-btn(color="error", icon, @click="removeItemFromArray(index)")
          v-icon delete
    v-btn.primary.ml-0(@click="addItem") {{ $t('modals.createPbehavior.steps.rrule.buttons.addExdate') }}
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
