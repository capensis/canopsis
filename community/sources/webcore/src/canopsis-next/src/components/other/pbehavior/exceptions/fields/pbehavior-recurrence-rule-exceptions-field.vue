<template>
  <v-layout class="gap-3" column>
    <pbehavior-exceptions-list
      v-if="exceptions.length"
      :exceptions="exceptions"
      @input="updateExceptions"
    />
    <pbehavior-exceptions-field
      v-field="exdates"
      :disabled="disabled"
      :with-exdate-type="withExdateType"
    >
      <template #actions="">
        <v-layout class="gap-2 fill-height" align-center>
          <v-btn
            color="primary"
            outlined
            @click="addException"
          >
            {{ $t('pbehavior.exceptions.create') }}
          </v-btn>
          <pbehavior-recurrence-rule-exceptions-list-menu
            :value="exceptions"
            @input="updateExceptions"
          />
        </v-layout>
      </template>
    </pbehavior-exceptions-field>
  </v-layout>
</template>

<script>
import { uid } from '@/helpers/uid';
import { convertDateToStartOfDayDateObject, convertDateToEndOfDayDateObject } from '@/helpers/date/date';

import { formArrayMixin } from '@/mixins/form';

import PbehaviorExceptionsList from '../../pbehaviors/partials/pbehavior-exceptions-list.vue';

import PbehaviorExceptionsField from './pbehavior-exceptions-field.vue';
import PbehaviorRecurrenceRuleExceptionsListMenu from './pbehavior-recurrence-rule-exceptions-list-menu.vue';

export default {
  components: {
    PbehaviorRecurrenceRuleExceptionsListMenu,
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
    updateExceptions(exceptions) {
      this.$emit('update:exceptions', exceptions);
    },

    addException() {
      this.addItemIntoArray({
        key: uid(),
        begin: convertDateToStartOfDayDateObject(),
        end: convertDateToEndOfDayDateObject(),
        type: '',
      });
    },
  },
};
</script>
