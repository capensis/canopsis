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
          :with-type="withExdateType",
          @delete="removeItemFromArray(index)"
        )
    v-layout(v-if="!disabled", row)
      slot(name="actions")
        v-flex
          v-btn.ml-0(color="secondary", @click="addExceptionDate") {{ $t('modals.createPbehaviorException.addDate') }}
    v-alert(:value="errors.has('exdates')", type="error") {{ errors.first('exdates') }}
</template>

<script>
import uid from '@/helpers/uid';

import { convertDateToStartOfDayDateObject } from '@/helpers/date/date';

import { formArrayMixin } from '@/mixins/form';
import { entitiesFieldPbehaviorFieldTypeMixin } from '@/mixins/entities/pbehavior/types-field';

import PbehaviorExceptionField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exception-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorExceptionField },
  mixins: [
    formArrayMixin,
    entitiesFieldPbehaviorFieldTypeMixin,
  ],
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
    withExdateType: {
      type: Boolean,
      default: false,
    },
  },
  created() {
    this.$validator.attach({
      name: 'exdates',
      rules: 'required:true',
      getter: () => !!this.exdates.length,
      vm: this,
    });
  },
  mounted() {
    this.fetchFieldPbehaviorTypesList();
  },
  methods: {
    addExceptionDate() {
      this.addItemIntoArray({
        key: uid(),
        begin: convertDateToStartOfDayDateObject(),
        end: convertDateToStartOfDayDateObject(),
        type: '',
      });
      this.$nextTick(() => this.$validator.validate('exdates'));
    },
  },
};
</script>
