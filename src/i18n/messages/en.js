export default {
  common: {
<<<<<<< HEAD
<<<<<<< HEAD
    actions: {
      close: 'Close',
      acknowledge: 'Acknowledge',
      acknowledgeAndReport: 'Acknowledge and report an incident',
      saveChanges: 'Save changes',
    },
    times: {
      second: 'second | seconds',
      minute: 'minute | minutes',
      hour: 'hour | hours',
      day: 'day | days',
      week: 'week | weeks',
      month: 'month | months',
      year: 'year | years',
    },
  },
  modals: {
    addAckEvent: {
      title: 'Add event type: ack',
      ticket: 'Ticket number',
      output: 'Note',
      ackResources: 'Ack resources',
      ackResourcesTooltip: 'Do you want to ack linked resources or not?',
    },
    addSnoozeEvent: {
      title: 'Add event type: snooze',
      duration: 'Duration',
    },
    addCancelEvent: {
      title: 'Add event type: cancel',
      output: 'Note',
    },
    addChangeStateEvent: {
      title: 'Add event type: change state',
      output: 'Note',
      states: {
        info: 'Info',
        minor: 'Minor',
        critical: 'Critical',
      },
    },
    addPbehavior: {
      title: 'Put a pbehavior on these elements ?',
      fields: {
        name: 'Name',
        start: 'Start',
        stop: 'End',
        reason: 'Reason',
        type: 'Type',
        rRuleQuestion: 'Put a rrule on this pbehavior ?',
      },
    },
  },
  tables: {
    alarmGeneral: {
      title: 'General',
      author: 'Author',
      connector: 'Connector',
      component: 'Component',
      resource: 'Resource',
    },
=======
    hello: 'Hello',
>>>>>>> parent of c7eef04... Add base modal structure and finish some modals
  },
  rRule: {
    advancedHint: 'Separate numbers with a comma',
    textLabel: 'Rrule',
    stringLabel: 'Summary',
    tabs: {
      simple: 'Simple',
      advanced: 'Advanced',
    },
    errors: {
      main: 'Please note that the Rrule you choose is not valid. We strongly advise you to modify it before saving changes to not causing trouble to Canopsis.',
    },
    fields: {
      freq: 'Frequency',
      until: 'Until',
      byweekday: 'By week day',
      count: 'Repeat',
      interval: 'Interval',
      wkst: 'Week start',
      bymonth: 'By month',
      bysetpos: {
        label: 'By set position',
        tooltip: 'If given, it must be one or many integers, positive or negative. Each given integer will specify an occurrence number, corresponding to the nth occurrence of the rule inside the frequency period. For example, a \'bysetpos\' of -1 if combined with a monthly frequency, and a \'byweekday\' of (Monday, Tuesday, Wednesday, Thursday, Friday), will result in the last work day of every month.',
      },
      bymonthday: {
        label: 'By month day',
        tooltip: 'If given, it must be one or many integers, meaning the month days to apply the recurrence to.',
      },
      byyearday: {
        label: 'By year day',
        tooltip: 'If given, it must be one or many integers, meaning the year days to apply the recurrence to.',
      },
      byweekno: {
        label: 'By week n°',
        tooltip: 'If given, it must be on or many integers, meaning the week numbers to apply the recurrence to. Week numbers have the meaning described in ISO8601, that is, the first week of the year is that containing at least four days of the new year.',
      },
      byhour: {
        label: 'By hour',
        tooltip: 'If given, it must be one or many integers, meaning the hours to apply the recurrence to.',
      },
      byminute: {
        label: 'By minute',
        tooltip: 'If given, it must be one or many integers, meaning the minutes to apply the recurrence to.',
      },
      bysecond: {
        label: 'By second',
        tooltip: 'If given, it must be one or many integers, meaning the seconds to apply the recurrence to.',
      },
    },
=======
    hello: 'Hello',
>>>>>>> parent of c7eef04... Add base modal structure and finish some modals
  },
};
