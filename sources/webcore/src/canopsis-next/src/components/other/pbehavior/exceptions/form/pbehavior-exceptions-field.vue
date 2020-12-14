<template lang="pug">
  v-layout(column)
    v-layout.mt-3(row)
      v-flex(xs12)
        v-layout
          slot(name="no-data", v-if="!exdates.length")
        pbehavior-exception-field.mb-3(
          v-for="(exdate, index) in exdates",
          v-field="exdates[index]",
          :key="exdate.key",
          :disabled="disabled",
          @delete="removeItemFromArray(index)"
        )
    v-layout(v-if="!disabled", row)
      v-flex
        v-btn.ml-0(color="secondary", @click="addExceptionDate") {{ $t('modals.createPbehaviorException.addDate') }}
    v-alert(:value="errors.has('exdates')", type="error") {{ errors.first('exdates') }}
</template>

<script>
import moment from 'moment';

import uid from '@/helpers/uid';

import formArrayMixin from '@/mixins/form/array';

import PbehaviorExceptionField from '@/components/other/pbehavior/calendar/partials/pbehavior-exception-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorExceptionField },
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  created() {
    this.$validator.attach({
      name: 'exdates',
      rules: 'required:true',
      getter: () => !!this.exdates.length,
      context: () => this,
      vm: this,
    });
  },
  methods: {
    addExceptionDate() {
      const startOfTodayMoment = moment().startOf('day');

      this.addItemIntoArray({
        key: uid(),
        begin: startOfTodayMoment.toDate(),
        end: startOfTodayMoment.toDate(),
        type: '',
      });
      this.$nextTick(() => this.$validator.validate('exdates'));
    },
  },
};
</script>
