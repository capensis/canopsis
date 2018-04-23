export default {
  common: {
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
  },
  tables: {
    alarmGeneral: {
      title: 'General',
      author: 'Author',
      connector: 'Connector',
      component: 'Component',
      resource: 'Resource',
    },
  },
};
