import { escape } from 'lodash';
import { computed } from 'vue';

import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { linkifyHtml, sanitizeHtml } from '@/helpers/html';
import { isSuccessTicketDeclaration } from '@/helpers/entities/declare-ticket/event/entity';

import { useI18n } from '@/hooks/i18n';

export const useExtraDetailsAckTooltip = (props) => {
  const { t, tc } = useI18n();

  const date = computed(() => convertDateToStringWithFormatForToday(props.ack.t));

  const tooltipContent = computed(() => {
    let content = `<strong>${t('alarm.actions.iconsTitles.ack')}</strong>
        <div>${t('common.by')} : ${escape(props.ack.a)}</div>
        <div>${t('common.date')} : ${date.value}</div>`;

    if (props.ack.initiator) {
      content += `<div>${t('common.initiator')} : ${escape(props.ack.initiator)}</div>`;
    }

    if (props.ack.m) {
      content += `<div class="c-extra-details__message">${tc('common.comment')} : ${escape(props.ack.m)}</div>`;
    }

    return `<div class="text-md-center">${content}</div>`;
  });

  return { tooltipContent };
};

export const useExtraDetailsCanceledTooltip = (props) => {
  const { t, tc } = useI18n();

  const date = computed(() => convertDateToStringWithFormatForToday(props.canceled.t));

  const tooltipContent = computed(() => {
    let content = `<strong>${t('alarm.actions.iconsTitles.canceled')}</strong>
        <div>${t('common.by')} : ${escape(props.canceled.a)}</div>
        <div>${t('common.date')} : ${date.value}</div>`;

    if (props.canceled.m) {
      content += `<div class="c-extra-details__message">${tc('common.comment')} : ${escape(props.canceled.m)}</div>`;
    }

    return `<div class="text-md-center">${content}</div>`;
  });

  return { tooltipContent };
};

export const useExtraDetailsChildrenTooltip = (props) => {
  const { t } = useI18n();

  const ruleName = computed(() => props.rule?.name ?? '');

  const tooltipContent = computed(() => (`<div class="text-md-center">
        <strong>${t('alarm.actions.iconsTitles.grouping')}</strong>
        <div>${t('common.title')} : ${escape(ruleName.value)}</div>
        <div>${t('alarm.actions.iconsFields.children')} : ${props.total}</div>
        <div>${t('alarm.fields.openedChildren')} : ${props.opened}</div>
        <div>${t('alarm.fields.closedChildren')} : ${props.closed}</div>
      </div>`));

  return { tooltipContent };
};

export const useExtraDetailsLastCommentTooltip = (props) => {
  const { t, tc } = useI18n();

  const date = computed(() => convertDateToStringWithFormatForToday(props.lastComment.t));
  const sanitizedLastComment = computed(() => sanitizeHtml(linkifyHtml(String(props.lastComment?.m ?? ''))));

  const tooltipContent = computed(() => (`<div class="text-md-center">
          <strong>${t('alarm.actions.iconsTitles.comment')}</strong>
          <div>${t('common.by')}: ${escape(props.lastComment.a)}</div>
          <div>${t('common.date')}: ${date.value}</div>
          <div class="c-extra-details__message">${tc('common.comment')}:&nbsp;<div>${sanitizedLastComment.value}</div></div>
        </div>`));

  return { tooltipContent };
};

export const useExtraDetailsParentsTooltip = (props) => {
  const { t, tc } = useI18n();

  const tooltipContent = computed(() => {
    const rulesContent = props.rules.map(({ name }) => `<div class="rule-name">&nbsp;${escape(name)}</div>`).join('');

    return `<div class="text-md-center">
          <strong>${t('alarm.actions.iconsTitles.grouping')}</strong>
          <div class="v-layout column">
            <div>${tc('common.rule', props.rules.length)}&nbsp;:</div>
            ${rulesContent}
          </div>
          <div>${t('alarm.actions.iconsFields.parents')} : ${props.total}</div>
        </div>`;
  });

  return { tooltipContent };
};

export const useExtraDetailsPbehaviorTooltip = (props) => {
  const { t } = useI18n();

  const tooltipContent = computed(() => {
    let content = `<div class="mt-2 font-weight-bold">${escape(props.pbehaviorInfo.name)}</div>`;

    if (props.pbehaviorInfo.author) {
      content += `<div>${t('common.author')}: ${escape(props.pbehaviorInfo.author)}</div>`;
    }

    if (props.pbehaviorInfo.type_name) {
      content += `<div>${t('common.type')}: ${escape(props.pbehaviorInfo.type_name)}</div>`;
    }

    if (props.pbehaviorInfo.reason_name) {
      content += `<div>${t('common.reason')}: ${escape(props.pbehaviorInfo.reason_name)}</div>`;
    }

    if (props.pbehaviorInfo.last_comment) {
      const lastCommentAuthorContent = props.pbehaviorInfo.last_comment.author
        ? `${escape(props.pbehaviorInfo.last_comment.author)}:&nbsp;`
        : '';

      content += `<div>
        ${t('common.lastComment')}:
        <div class="ml-2">
          -&nbsp;
          ${lastCommentAuthorContent}
          ${escape(props.pbehaviorInfo.last_comment.message)}
        </div>
      </div>`;
    }

    return `<div>
          <strong>${t('alarm.actions.iconsTitles.pbehaviors')}</strong>
          <div>
            ${content}
          </div>
        </div>`;
  });

  return { tooltipContent };
};

export const useExtraDetailsSnoozeTooltip = (props) => {
  const { t, tc } = useI18n();

  const date = computed(() => convertDateToStringWithFormatForToday(props.snooze.t));
  const end = computed(() => convertDateToStringWithFormatForToday(props.snooze.val));

  const tooltipContent = computed(() => {
    let content = `<strong>${t('alarm.actions.iconsTitles.snooze')}</strong>
        <div>${t('common.by')} : ${escape(props.snooze.a)}</div>
        <div>${t('common.date')} : ${date.value}</div>
        <div>${t('common.end')} : ${end.value}</div>`;

    if (props.snooze.initiator) {
      content += `<div>${t('common.initiator')} : ${escape(props.snooze.initiator)}</div>`;
    }

    if (props.snooze.m) {
      content += `<div class="c-extra-details__message">${tc('common.comment')} : ${escape(props.snooze.m)}</div>`;
    }

    return `<div class="text-md-center">${content}</div>`;
  });

  return { tooltipContent };
};

export const useExtraDetailsTicketTooltip = (props) => {
  const { t, tc } = useI18n();

  const getTicketStatusText = ticket => t(`common.${isSuccessTicketDeclaration(ticket) ? 'ok' : 'failed'}`);
  const convertDateWithToday = date => convertDateToStringWithFormatForToday(date);

  const shownTickets = computed(() => props.tickets.slice(0, props.limit));

  const tooltipContent = computed(() => {
    const content = shownTickets.value.reduce((acc, ticket) => {
      let ticketContent = `<strong>${ticket.ticket_rule_name || ''} ${getTicketStatusText(ticket)}</strong>
          <div>${t('common.by')} : ${escape(ticket.a)}</div>
          <div>${t('common.date')} : ${convertDateWithToday(ticket.t)}</div>`;

      if (ticket.ticket) {
        ticketContent += `<div>${t('alarm.actions.iconsFields.ticketNumber')} : ${escape(ticket.ticket)}</div>`;
      }

      if (ticket.ticket_comment) {
        ticketContent += `<div>${tc('common.comment')} : ${escape(ticket.ticket_comment)}</div>`;
      }

      return `${acc}<div class="text-md-center">${ticketContent}</div>`;
    }, '');

    const otherTickets = props.tickets.length > props.limit ? `<i>${t('alarm.otherTickets')}</i>` : '';

    return `<div class="layout extra-details-ticket__list column">${content}</div>
              <div class="mt-2">${otherTickets}</div>`;
  });

  return { tooltipContent };
};
