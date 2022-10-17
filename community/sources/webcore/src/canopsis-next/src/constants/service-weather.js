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

/* TODO: Should be fixed after backend integration */
export const SERVICE_WEATHER_STATE_COUNTERS = {
  alarms: 'alarms',
  dependencies: 'depends',
  ok: 'state.info',
  underPbehavior: 'pbehavior_counters',
  minor: 'state.minor',
  major: 'state.major',
  critical: 'state.critical',
  acknowledged: 'acknowledged',
  notAcknowledged: 'not_acknowledged',
  acknowledgedUnderPbehavior: 'acknowledged_under_pbehavior',
};

export const SERVICE_WEATHER_MAX_STATE_COUNTERS = 5;
