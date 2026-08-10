import { escape } from 'lodash';
import { computed } from 'vue';

import { DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES } from '@/constants';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { useI18n } from '@/hooks/i18n';

/**
 * Function to provide tooltip content for Ticket action in extra details.
 *
 * @param {Object} props - The props object containing ticket information.
 * @returns {Object} - The tooltip content for Ticket action.
 */
export const useExtraDetailsTicketTooltip = (props) => {
  const { t } = useI18n();

  const getTicketStatusText = ticket =>
    t(
      `declareTicket.status.${ticket.ticket_status ?? DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES.unknown}`,
    );
  const convertDateWithToday = date => convertDateToStringWithFormatForToday(date);

  const shownTickets = computed(() =>
    props.failedTicket ? [props.failedTicket] : props.tickets.slice(0, props.limit),
  );

  const tooltipContent = computed(() => {
    const content = shownTickets.value.reduce((acc, ticket) => {
      let ticketContent = `<strong>${ticket.ticket_rule_name ? `${ticket.ticket_rule_name} - ` : ''}${getTicketStatusText(ticket)}</strong>
          <div>${t('common.by')} : ${escape(ticket.a)}</div>
          <div>${t('common.date')} : ${convertDateWithToday(ticket.t)}</div>`;

      if (ticket.ticket) {
        ticketContent += `<div>${t('alarm.actions.iconsFields.ticketNumber')} : ${escape(ticket.ticket)}</div>`;
      }

      if (ticket.ticket_comment) {
        ticketContent += `<div>${t('common.comment')} : ${escape(ticket.ticket_comment)}</div>`;
      }

      return `${acc}<div class="text-md-center">${ticketContent}</div>`;
    }, '');

    const otherTickets =
      props.tickets.length > props.limit ? `<i>${t('alarm.otherTickets')}</i>` : '';

    return `<div class="layout extra-details-ticket__list column">${content}</div>
              <div class="mt-2">${otherTickets}</div>`;
  });

  return { tooltipContent };
};
