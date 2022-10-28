import { COLORS } from '@/config';

export const DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE = `<p><strong><span style="font-size: 18px;">{{entity.name}}</span></strong></p>
<hr id="null">
<p>{{ entity.output }}</p>
<p> Dernière mise à jour : {{ timestamp entity.last_update_date }}</p>`;

export const DEFAULT_SERVICE_WEATHER_MODAL_TEMPLATE = '{{ entities name="entity._id" }}';

export const DEFAULT_SERVICE_WEATHER_ENTITY_TEMPLATE = `<ul>
    <li><strong>Libellé</strong> : {{entity.name}}</li>
</ul>`;

export const DEFAULT_WIDGET_MARGIN = {
  top: 1,
  right: 1,
  bottom: 1,
  left: 1,
};

export const SERVICE_WEATHER_PATTERN_FIELDS = {
  grey: 'is_grey',
  primaryIcon: 'icon',
  secondaryIcon: 'secondary_icon',
  state: 'state.val',
};

export const SERVICE_WEATHER_STATE_COUNTERS = {
  all: 'all',
  active: 'active',
  depends: 'depends',
  ok: 'state.ok',
  minor: 'state.minor',
  major: 'state.major',
  critical: 'state.critical',
  acked: 'acked',
  unacked: 'unacked',
  underPbehavior: 'under_pbh',
  ackedUnderPbehavior: 'acked_under_pbh',
};

export const SERVICE_WEATHER_TEMPLATE_COUNTERS_BY_STATE_COUNTERS = {
  [SERVICE_WEATHER_STATE_COUNTERS.all]: '.All',
  [SERVICE_WEATHER_STATE_COUNTERS.active]: '.Active',
  [SERVICE_WEATHER_STATE_COUNTERS.depends]: '.Depends',
  [SERVICE_WEATHER_STATE_COUNTERS.ok]: '.State.Ok',
  [SERVICE_WEATHER_STATE_COUNTERS.minor]: '.State.Minor',
  [SERVICE_WEATHER_STATE_COUNTERS.major]: '.State.Major',
  [SERVICE_WEATHER_STATE_COUNTERS.critical]: '.State.Critical',
  [SERVICE_WEATHER_STATE_COUNTERS.acked]: '.Acknowledged',
  [SERVICE_WEATHER_STATE_COUNTERS.unacked]: '.NotAcknowledged',
  [SERVICE_WEATHER_STATE_COUNTERS.underPbehavior]: '.UnderPbehavior',
  [SERVICE_WEATHER_STATE_COUNTERS.ackedUnderPbehavior]: '.AcknowledgedUnderPbh',
};

export const SERVICE_WEATHER_STATE_COUNTERS_ICONS = {
  [SERVICE_WEATHER_STATE_COUNTERS.all]: 'notification_important',
  [SERVICE_WEATHER_STATE_COUNTERS.active]: '$vuetify.icons.notification_important_stroke',
  [SERVICE_WEATHER_STATE_COUNTERS.depends]: '$vuetify.icons.mediation',
  [SERVICE_WEATHER_STATE_COUNTERS.ok]: 'check_circle',
  [SERVICE_WEATHER_STATE_COUNTERS.minor]: '$vuetify.icons.warning_stroke',
  [SERVICE_WEATHER_STATE_COUNTERS.major]: '$vuetify.icons.warning_stroke',
  [SERVICE_WEATHER_STATE_COUNTERS.critical]: '$vuetify.icons.warning_stroke',
  [SERVICE_WEATHER_STATE_COUNTERS.acked]: 'playlist_add_check',
  [SERVICE_WEATHER_STATE_COUNTERS.unacked]: 'playlist_play',
  [SERVICE_WEATHER_STATE_COUNTERS.underPbehavior]: 'build',
  [SERVICE_WEATHER_STATE_COUNTERS.ackedUnderPbehavior]: '$vuetify.icons.playlist_build',
};

export const SERVICE_WEATHER_STATE_COUNTERS_COLORS = {
  [SERVICE_WEATHER_STATE_COUNTERS.all]: COLORS.error,
  [SERVICE_WEATHER_STATE_COUNTERS.active]: COLORS.error,
  [SERVICE_WEATHER_STATE_COUNTERS.minor]: COLORS.state.minor,
  [SERVICE_WEATHER_STATE_COUNTERS.major]: COLORS.state.major,
  [SERVICE_WEATHER_STATE_COUNTERS.critical]: COLORS.state.critical,
};

export const SERVICE_WEATHER_MAX_STATE_COUNTERS = 5;
