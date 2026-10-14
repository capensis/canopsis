<template>
  <c-layout gap="3" justify-center column>
    <c-label :label="$t('common.recurrence')" />

    <div v-if="!hasRecurrenceRule">
      <v-btn color="primary" variant="outlined" @click="showCreateRecurrenceRuleModal">
        {{ $t('pbehavior.buttons.addRRule') }}
      </v-btn>
    </div>

    <template v-else>
      <c-layout>
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

        <c-action-btn type="delete" @click="showConfirmRemoveRecurrenceRuleModal" />
      </c-layout>
    </template>
  </c-layout>
</template>

<script>
import { isEmpty } from 'lodash';
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { formToRrule } from '@/helpers/entities/shared/recurrence-rule/form';

import { useModals } from '@/hooks/modals';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';

export default {
  components: {
    RecurrenceRuleInformation,
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const modals = useModals();

    const hasRecurrenceRule = computed(() => !isEmpty(props.form.rrule));

    const rruleBodyForInformation = computed(() => formToRrule(props.form.rrule));

    const showConfirmRemoveRecurrenceRuleModal = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => emit('update:form', { ...props.form, rrule: '' }),
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
          action: ({ rrule, exdates, exceptions }) => emit('update:form', {
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
