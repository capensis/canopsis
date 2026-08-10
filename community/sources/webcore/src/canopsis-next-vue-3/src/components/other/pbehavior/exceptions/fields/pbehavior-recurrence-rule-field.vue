<template>
  <c-layout gap="3" justify="center" column>
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
          :exdates="exdates"
          :exceptions="exceptions"
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

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import { MODALS } from '@/constants';
import { formToRrule } from '@/helpers/entities/shared/recurrence-rule/form';
import { useModals } from '@/hooks/modals';

export default {
  emits: ['update:rrule', 'update:exdates', 'update:exceptions'],
  components: {
    RecurrenceRuleInformation,
  },
  props: {
    rrule: {
      type: String,
      default: '',
    },
    exdates: {
      type: Array,
      default: () => [],
    },
    exceptions: {
      type: Array,
      default: () => [],
    },
    start: {
      type: [Number, Object],
      default: null,
    },
    withExdateType: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const modals = useModals();

    const hasRecurrenceRule = computed(() => !isEmpty(props.rrule));

    const rruleBodyForInformation = computed(() => formToRrule(props.rrule));

    const showConfirmRemoveRecurrenceRuleModal = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => emit('update:rrule', ''),
        },
      });
    };

    const showCreateRecurrenceRuleModal = () => {
      modals.show({
        name: MODALS.createRecurrenceRule,
        config: {
          rrule: props.rrule,
          exdates: props.exdates,
          exceptions: props.exceptions,
          start: props.start,
          withExdateType: props.withExdateType,
          action: ({ rrule, exdates, exceptions }) => {
            emit('update:rrule', rrule);
            emit('update:exdates', exdates);
            emit('update:exceptions', exceptions);
          },
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
