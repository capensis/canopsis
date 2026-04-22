<template>
  <c-chip
    v-if="statusColor"
    :color="statusColor"
    :text-color="textColor"
    class="declare-ticket-rule-ticket-status-chip px-2"
    small
  >
    <strong>{{ statusText }}</strong>
  </c-chip>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';
import { DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

const STATUS_VALUE_TO_KEY = {
  [DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.unknown]: 'unknown',
  [DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.open]: 'open',
  [DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.assigned]: 'assigned',
  [DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.inProgress]: 'inProgress',
  [DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.closed]: 'closed',
};

export default {
  props: {
    value: {
      type: Number,
      default: DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.unknown,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const statusKey = computed(() => STATUS_VALUE_TO_KEY[props.value] ?? 'unknown');
    const statusColor = computed(() => (
      COLORS.declareTicketRuleTicketStatusChip?.[statusKey.value]
      ?? COLORS.declareTicketRuleTicketStatusChip?.unknown
    ));
    const statusText = computed(() => t(`declareTicket.status.${props.value}`));

    const textColor = computed(() => {
      if (props.value === DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.closed) {
        return 'var(--v-text-dark-primary, #FFFFFF)';
      }

      return 'var(--v-text-light-primary, rgba(0, 0, 0, 0.87))';
    });

    return {
      statusColor,
      statusText,
      textColor,
    };
  },
};
</script>

<style lang="scss" scoped>
.declare-ticket-rule-ticket-status-chip {
  border-radius: 12px;
}
</style>
