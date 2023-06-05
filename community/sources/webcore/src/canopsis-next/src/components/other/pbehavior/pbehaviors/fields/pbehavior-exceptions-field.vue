<template lang="pug">
  div
    h3.my-3.grey--text {{ $t('pbehavior.exceptions.title') }}
    v-divider
    pbehavior-exception-list(v-if="exceptions.length", :exceptions="exceptions")
    v-layout.mt-3(column)
      v-flex(xs12)
        v-alert(
          v-if="!hasExceptionsOrExdates",
          :value="true",
          type="info"
        ) {{ $t('pbehavior.exceptions.emptyExceptions') }}
      pbehavior-exception-field.mb-3(
        v-for="(exdate, index) in exdates",
        v-field="exdates[index]",
        :key="exdate.key",
        :with-type="withExdateType",
        :disabled="disabled",
        @delete="removeItemFromArray(index)"
      )
    v-layout(v-if="!disabled", row)
      v-flex
        v-btn.ml-0(outline, @click="addException") {{ $t('pbehavior.exceptions.create') }}
      v-flex
        v-btn.mr-0(outline, @click="showSelectExceptionModal") {{ $t('pbehavior.exceptions.choose') }}
</template>

<script>
import { MODALS } from '@/constants';

import uid from '@/helpers/uid';
import { convertDateToStartOfDayDateObject } from '@/helpers/date/date';

import { formArrayMixin } from '@/mixins/form';

import PbehaviorExceptionField from './pbehavior-exception-field.vue';

import PbehaviorExceptionList from '../partials/pbehavior-exception-list.vue';

export default {
  components: { PbehaviorExceptionList, PbehaviorExceptionField },
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
    exceptions: {
      type: Array,
      default: () => [],
    },
    withExdateType: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasExceptionsOrExdates() {
      return this.exdates.length || this.exceptions.length;
    },
  },
  methods: {
    showSelectExceptionModal() {
      this.$modals.show({
        name: MODALS.selectExceptionsLists,
        config: {
          exceptions: this.exceptions,
          action: exceptions => this.$emit('update:exceptions', exceptions),
        },
      });
    },

    addException() {
      this.addItemIntoArray({
        key: uid(),
        begin: convertDateToStartOfDayDateObject(),
        end: convertDateToStartOfDayDateObject(),
        type: '',
      });
    },
  },
};
</script>
