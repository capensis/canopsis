import {
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_TEMPLATE_FIELDS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_TEMPLATE_FIELDS,
  PBEHAVIOR_INFO_FIELDS,
  VARIABLES_STEP_FIELDS,
} from '@/constants';

const createStepVariable = (t, textKey, value) => ({
  text: t(textKey),
  value,
});

const createCommonVariables = t => ({
  stepValueVariable: createStepVariable(t, 'common.value', VARIABLES_STEP_FIELDS.value),
  stepTimestampVariable: createStepVariable(t, 'common.timestamp', VARIABLES_STEP_FIELDS.timestamp),
  stepMessageVariable: createStepVariable(t, 'common.message', VARIABLES_STEP_FIELDS.message),
  stepAuthorVariable: createStepVariable(t, 'common.author', VARIABLES_STEP_FIELDS.author),
});

export const prepareAlarmTemplateVariables = t => [
    { value: ENTITY_TEMPLATE_FIELDS.id },
    { value: ALARM_TEMPLATE_FIELDS.ack },
    { value: ALARM_TEMPLATE_FIELDS.state },
    { value: ALARM_TEMPLATE_FIELDS.status },
    { value: ALARM_TEMPLATE_FIELDS.ticket },
    { value: ALARM_TEMPLATE_FIELDS.component },
    { value: ALARM_TEMPLATE_FIELDS.connector },
    { value: ALARM_TEMPLATE_FIELDS.connectorName },
    { value: ALARM_TEMPLATE_FIELDS.resource },
    { value: ALARM_TEMPLATE_FIELDS.creationDate },
    { value: ALARM_TEMPLATE_FIELDS.displayName },
    { value: ALARM_TEMPLATE_FIELDS.output },
    { value: ALARM_TEMPLATE_FIELDS.lastUpdateDate },
    { value: ALARM_TEMPLATE_FIELDS.lastEventDate },
    { value: ALARM_TEMPLATE_FIELDS.pbehaviorInfo },
    { value: ALARM_TEMPLATE_FIELDS.duration },
    { value: ALARM_TEMPLATE_FIELDS.eventsCount },
  ].map(variable => ({
    value: `{{ ${variable.value} }}`,
    text: t(ALARM_FIELDS_TO_LABELS_KEYS[variable.value.replace('alarm.', '')], 2),
  }));

export const prepareEntityTemplateVariables = (t, entityInfos = []) => {
  const { stepValueVariable, stepTimestampVariable, stepMessageVariable, stepAuthorVariable } =
    createCommonVariables(t);
  const ackVariables = [stepTimestampVariable, stepAuthorVariable, stepMessageVariable];
  const ticketVariables = [stepValueVariable];
  const pbehaviorInfoVariables = [
    {
      text: t('pbehavior.pbehaviorType'),
      value: PBEHAVIOR_INFO_FIELDS.typeName,
    },
    {
      text: t('pbehavior.pbehaviorReason', 1),
      value: PBEHAVIOR_INFO_FIELDS.reason,
    },
    {
      text: t('pbehavior.pbehaviorName'),
      value: PBEHAVIOR_INFO_FIELDS.name,
    },
    {
      text: t('pbehavior.pbehaviorCanonicalType'),
      value: PBEHAVIOR_INFO_FIELDS.canonicalType,
    },
  ];
  const infosSubVariables = [
    stepValueVariable,
    createStepVariable(t, 'common.description', 'description'),
  ];
  const infosVariables = entityInfos.map(({ value }) => ({
    text: value,
    value,
    variables: infosSubVariables,
  }));
  const snoozeVariables = [stepTimestampVariable, stepAuthorVariable, stepMessageVariable];
  const lastCommentVariables = [stepTimestampVariable, stepAuthorVariable, stepMessageVariable];

  return [
    { value: ENTITY_TEMPLATE_FIELDS.id },
    { value: ENTITY_TEMPLATE_FIELDS.name },
    {
      value: ENTITY_TEMPLATE_FIELDS.infos,
      variables: infosVariables.length ? infosVariables : undefined,
    },
    { value: ENTITY_TEMPLATE_FIELDS.connector },
    { value: ENTITY_TEMPLATE_FIELDS.connectorName },
    { value: ENTITY_TEMPLATE_FIELDS.component },
    { value: ENTITY_TEMPLATE_FIELDS.resource },
    {
      value: ENTITY_TEMPLATE_FIELDS.state,
      variables: [stepValueVariable],
    },
    {
      value: ENTITY_TEMPLATE_FIELDS.status,
      variables: [stepValueVariable],
    },
    {
      value: ENTITY_TEMPLATE_FIELDS.snooze,
      variables: snoozeVariables,
    },
    {
      value: ENTITY_TEMPLATE_FIELDS.ack,
      variables: ackVariables,
    },
    { value: ENTITY_TEMPLATE_FIELDS.lastUpdateDate },
    { value: ENTITY_TEMPLATE_FIELDS.impactLevel },
    { value: ENTITY_TEMPLATE_FIELDS.impactState },
    { value: ENTITY_TEMPLATE_FIELDS.categoryName },
    {
      value: ENTITY_TEMPLATE_FIELDS.pbehaviorInfo,
      variables: pbehaviorInfoVariables,
    },
    { value: ENTITY_TEMPLATE_FIELDS.alarmCreationDate },
    {
      value: ENTITY_TEMPLATE_FIELDS.ticket,
      variables: ticketVariables,
    },
    { value: ENTITY_TEMPLATE_FIELDS.statsOk },
    { value: ENTITY_TEMPLATE_FIELDS.statsKo },
    { value: ENTITY_TEMPLATE_FIELDS.alarmDisplayName },
    { value: ENTITY_TEMPLATE_FIELDS.links },
    {
      value: ENTITY_TEMPLATE_FIELDS.alarmLastComment,
      variables: lastCommentVariables,
    },
    {
      value: ENTITY_TEMPLATE_FIELDS.lastComment,
      variables: lastCommentVariables,
    },
    { value: ENTITY_TEMPLATE_FIELDS.tags },
  ].map(variable => ({
    ...variable,
    text: t(ENTITY_FIELDS_TO_LABELS_KEYS[variable.value.replace('entity.', '')], 2),
  }));
};
