import { COLORS } from '@/config';

export const STATS_OPTIONS = {
  recursive: 'recursive',
  states: 'states',
  authors: 'authors',
  sla: 'sla',
};

export const STATS_TYPES = {
  alarmsCreated: {
    value: 'alarms_created',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  alarmsResolved: {
    value: 'alarms_resolved',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  alarmsCanceled: {
    value: 'alarms_canceled',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  alarmsAcknowledged: {
    value: 'alarms_acknowledged',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  ackTimeSla: {
    value: 'ack_time_sla',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors, STATS_OPTIONS.sla],
  },
  resolveTimeSla: {
    value: 'resolve_time_sla',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors, STATS_OPTIONS.sla],
  },
  timeInState: {
    value: 'time_in_state',
    options: [STATS_OPTIONS.states],
  },
  stateRate: {
    value: 'state_rate',
    options: [STATS_OPTIONS.states],
  },
  mtbf: {
    value: 'mtbf',
    options: [STATS_OPTIONS.recursive],
  },
  currentState: {
    value: 'current_state',
    options: [],
  },
  ongoingAlarms: {
    value: 'ongoing_alarms',
    options: [STATS_OPTIONS.states],
  },
  currentOngoingAlarms: {
    value: 'current_ongoing_alarms',
    options: [STATS_OPTIONS.states],
  },
  currentOngoingAlarmsWithAck: {
    value: 'current_ongoing_alarms_with_ack',
    options: [STATS_OPTIONS.states],
  },
  currentOngoingAlarmsWithoutAck: {
    value: 'current_ongoing_alarms_without_ack',
    options: [STATS_OPTIONS.states],
  },
};

export const STATS_CRITICITY = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
};

export const STATS_DEFAULT_COLOR = COLORS.statsDefault;

export const STATS_DISPLAY_MODE = {
  value: 'value',
  criticity: 'criticity',
};

export const STATS_DISPLAY_MODE_PARAMETERS = {
  criticityLevels: {
    ok: 0,
    minor: 10,
    major: 20,
    critical: 30,
  },
  colors: {
    ok: COLORS.state.ok,
    minor: COLORS.state.minor,
    major: COLORS.state.major,
    critical: COLORS.state.critical,
  },
};

export const STATS_CURVES_POINTS_STYLES = {
  circle: 'circle',
  cross: 'cross',
  crossRot: 'crossRot',
  dash: 'dash',
  line: 'line',
  rect: 'rect',
  rectRounded: 'rectRounded',
  rectRot: 'rectRot',
  star: 'star',
  triangle: 'triangle',
};
