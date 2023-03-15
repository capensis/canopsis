import { WIDGET_TEMPLATES_TYPES } from '@/constants';

export default {
  types: {
    [WIDGET_TEMPLATES_TYPES.alarmColumns]: 'General : Alarms columns',
    [WIDGET_TEMPLATES_TYPES.entityColumns]: 'General : Entities columns',
    [WIDGET_TEMPLATES_TYPES.alarmMoreInfos]: 'Alarm list : More Infos template',
    [WIDGET_TEMPLATES_TYPES.weatherItem]: 'Service Weather : Tile template',
    [WIDGET_TEMPLATES_TYPES.weatherModal]: 'Service Weather : Modal template',
    [WIDGET_TEMPLATES_TYPES.weatherEntity]: 'Service Weather : Entity template',
  },
  errors: {
    columnsRequired: 'You should add at least one column.',
  },
};
