<template lang="pug">
  div
    h3.my-3.grey--text {{ $t('pbehaviorExceptions.title') }}
    v-divider
    pbehavior-exception-list(v-if="exceptions.length", :exceptions="exceptions")
    v-layout.mt-3(row)
      v-flex(xs12)
        v-alert(:value="!exdates.length", type="info") {{ $t('pbehaviorExceptions.emptyExceptions') }}
        pbehavior-exception-field.mb-3(
          v-for="(exdate, index) in exdates",
          v-field="exdates[index]",
          :key="exdate.key",
          @delete="removeItemFromArray(index)"
        )
    v-layout(row)
      v-flex
        v-btn.ml-0(outline, @click="addException") {{ $t('pbehaviorExceptions.create') }}
      v-flex
        v-btn.mr-0(outline, @click="showSelectExceptionModal") {{ $t('pbehaviorExceptions.choose') }}
</template>

<script>
import moment from 'moment';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

import formArrayMixin from '@/mixins/form/array';

import PbehaviorExceptionField from '@/components/other/pbehavior/calendar/partials/pbehavior-exception-field.vue';
import PbehaviorExceptionList from '@/components/other/pbehavior/calendar/partials/pbehavior-exception-list.vue';

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
      const startOfTodayMoment = moment().startOf('day');

      this.addItemIntoArray({
        key: uid(),
        begin: startOfTodayMoment.toDate(),
        end: startOfTodayMoment.toDate(),
        type: '',
      });
    },
  },
};
</script>
