<template>
  <c-advanced-data-table
    :headers="headers"
    :items="tickets"
    :total-items="tickets.length"
    hide-actions
    disable-initial-sort
  >
    <template #ticket_url="{ item }">
      <template v-if="item.ticket_url">
        <a
          v-if="isValidTicketUrl(item.ticket_url)"
          :href="item.ticket_url"
          :title="item.ticket_url_title"
          target="_blank"
        >{{ item.ticket_url_title || item.ticket_url }}</a>
        <span v-else>{{ item.ticket_url }}</span>
      </template>
    </template>
    <template #t="{ item }">
      {{ item.t | date }}
    </template>
    <template #_t="{ item }">
      <c-help-icon
        v-bind="getIconProps(item)"
        top
      />
    </template>
    <template #metaalarm="{ item }">
      <v-icon
        v-if="item.ticket_meta_alarm_id === parentAlarmId"
        top
      >
        low_priority
      </v-icon>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { isValidUrl } from '@/plugins/validator/helpers/is-valid-url';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    tickets: {
      type: Array,
      required: true,
    },
    parentAlarmId: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const { t, tc } = useI18n();

    const headers = computed(() => [
      { text: t('declareTicket.ticketURL'), value: 'ticket_url' },
      { text: t('declareTicket.ticketID'), value: 'ticket' },
      { text: t('common.systemName'), value: 'ticket_system_name' },
      { text: t('declareTicket.ruleName'), value: 'ticket_rule_name' },
      { text: t('common.date'), value: 't' },
      { text: t('common.status'), value: '_t' },
      props.parentAlarmId && { text: t('alarm.metaAlarm'), value: 'metaalarm' },
      { text: t('common.author'), value: 'a' },
      { text: tc('common.comment'), value: 'ticket_comment' },
    ].filter(Boolean));

    const getIconProps = item => ({
      icon: 'check_circle', // TODO: hardcoded for #6000
      color: 'primary', // TODO: hardcoded for #6000
      text: item.m,
      maxWidth: 400,
    });

    return {
      headers,
      getIconProps,
      isValidTicketUrl: isValidUrl,
    };
  },
};
</script>
