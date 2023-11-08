<template>
  <div>
    <h3 class="text--secondary">
      {{ $t('pbehavior.exceptions.title') }}
    </h3>
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
      <template #no-data="">
        <c-alert
          :value="!hasExceptionsOrExdates"
          type="info"
        >
          {{ $t('pbehavior.exceptions.emptyExceptions') }}
        </c-alert>
      </template>
      <template #actions="">
        <v-btn
          class="mr-2"
          color="primary"
          @click="addException"
        >
          {{ $t('pbehavior.exceptions.create') }}
        </v-btn>
        <pbehavior-recurrence-rule-exceptions-list-menu
          :value="exceptions"
          @input="updateExceptions"
        />
      </template>
    </pbehavior-exceptions-field>
  </div>
</template>

<script>
import { uid } from '@/helpers/uid';
import { convertDateToStartOfDayDateObject } from '@/helpers/date/date';

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
        end: convertDateToStartOfDayDateObject(),
        type: '',
      });
    },
  },
};
</script>
