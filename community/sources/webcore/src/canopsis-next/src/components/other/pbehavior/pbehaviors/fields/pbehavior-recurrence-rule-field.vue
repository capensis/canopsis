<template>
  <v-layout align-center>
    <v-btn
      v-if="!hasRecurrenceRule"
      color="primary"
      outlined
      @click="showCreateRecurrenceRuleModal"
    >
      {{ $t('pbehavior.buttons.addRRule') }}
    </v-btn>
    <template v-else>
      <v-layout>
        <recurrence-rule-information
          :rrule="rruleBodyForInformation"
          :exdates="form.exdates"
          :exceptions="form.exceptions"
          class="mr-2"
        />
        <c-action-btn
          :tooltip="$t('pbehavior.buttons.editRrule')"
          type="edit"
          @click="showCreateRecurrenceRuleModal"
        />
        <c-action-btn
          type="delete"
          @click="showConfirmRemoveRecurrenceRuleModal"
        />
      </v-layout>
    </template>
  </v-layout>
</template>

<script>
import { computed } from 'vue';
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import { formToRrule } from '@/helpers/entities/shared/recurrence-rule/form';

import { useModals } from '@/hooks/modals';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';

export default {
  components: {
    RecurrenceRuleInformation,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    withExdateType: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const modals = useModals();

    const hasRecurrenceRule = computed(() => !isEmpty(props.form.rrule));

    const rruleBodyForInformation = computed(() => formToRrule(props.form.rrule));

    const showConfirmRemoveRecurrenceRuleModal = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => emit('input', { ...props.form, rrule: '' }),
        },
      });
    };

    const showCreateRecurrenceRuleModal = () => {
      modals.show({
        name: MODALS.createRecurrenceRule,
        config: {
          rrule: props.form.rrule,
          exdates: props.form.exdates,
          exceptions: props.form.exceptions,
          start: props.form.tstart,
          withExdateType: props.withExdateType,
          action: ({ rrule, exdates, exceptions }) => emit('input', {
            ...props.form,
            rrule,
            exdates,
            exceptions,
          }),
        },
      });
    };

    return {
      hasRecurrenceRule,
      rruleBodyForInformation,
      showConfirmRemoveRecurrenceRuleModal,
      showCreateRecurrenceRuleModal,
    };
  },
};
</script>
