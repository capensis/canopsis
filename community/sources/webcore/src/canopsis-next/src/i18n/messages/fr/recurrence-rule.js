import { PBEHAVIOR_RRULE_PERIODS_RANGES } from '@/constants';

export default {
  advancedHint: 'Séparer les nombres par une virgule',
  freq: 'Fréquence',
  until: 'Jusqu\'à',
  byweekday: 'Par jour de la semaine',
  count: 'Répéter',
  interval: 'Intervalle',
  wkst: 'Semaine de début',
  bymonth: 'Par mois',
  bysetpos: 'Par position',
  bymonthday: 'Par jour du mois',
  byyearday: 'Par jour de l\'année',
  byweekno: 'Par semaine n°',
  byhour: 'Par heure',
  byminute: 'Par minute',
  bysecond: 'Par seconde',
  tabs: {
    simple: 'Simple',
    advanced: 'Avancé',
  },
  errors: {
    main: 'La récurrence choisie n\'est pas valide. Nous vous recommandons de la modifier avant de sauvegarder',
  },
  periodsRanges: {
    [PBEHAVIOR_RRULE_PERIODS_RANGES.thisWeek]: 'Cette semaine',
    [PBEHAVIOR_RRULE_PERIODS_RANGES.nextWeek]: 'Semaine prochaine',
    [PBEHAVIOR_RRULE_PERIODS_RANGES.next2Weeks]: 'Deux prochaines semaines',
    [PBEHAVIOR_RRULE_PERIODS_RANGES.thisMonth]: 'Ce mois',
    [PBEHAVIOR_RRULE_PERIODS_RANGES.nextMonth]: 'Le mois prochain',
  },
  tooltips: {
    bysetpos: 'Si renseigné, doit être un ou plusieurs nombres entiers, positifs ou négatifs. Chaque entier correspondra à la ènième occurence de la règle dans l\'intervalle de fréquence. Par exemple, une \'bysetpos\' de -1 combinée à une fréquence mensuelle, et une \'byweekday\' de (lundi, mardi, mercredi, jeudi, vendredi), va nous donner le dernier jour travaillé de chaque mois',
    bymonthday: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours du mois auxquels s\'appliquera la récurrence.',
    byyearday: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours de l\'année auxquels  s\'appliquera la récurrence.',
    byweekno: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux numéros de semaine auxquelles s\'appliquera la récurrence. Les numéros de semaines sont ceux de ISO8601, la première semaine de l\'année étant celle contenant au moins 4 jours de cette année.',
    byhour: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux heures auxquelles s\'appliquera la récurrence.',
    byminute: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux minutes auxquelles s\'appliquera la récurrence.',
    bysecond: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux secondes auxquelles s\'appliquera la récurrence.',
  },
};
