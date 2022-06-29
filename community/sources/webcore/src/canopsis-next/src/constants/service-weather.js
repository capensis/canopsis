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
