<template lang="pug">
  div
    h3.my-3.grey--text {{ $t('pbehavior.exceptions.title') }}
    v-divider
    pbehavior-exceptions-list(v-if="exceptions.length", :exceptions="exceptions")
    pbehavior-exceptions-field(
      v-field="exdates",
      :disabled="disabled",
      :with-exdate-type="withExdateType"
    )
      template(#no-data="")
        c-alert(
          :value="!hasExceptionsOrExdates",
          type="info"
        ) {{ $t('pbehavior.exceptions.emptyExceptions') }}
      template(#actions="")
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

import PbehaviorExceptionsList from '../../pbehaviors/partials/pbehavior-exceptions-list.vue';

import PbehaviorExceptionsField from './pbehavior-exceptions-field.vue';

export default {
  components: {
    PbehaviorExceptionsList,
    PbehaviorExceptionsField,
  },
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
