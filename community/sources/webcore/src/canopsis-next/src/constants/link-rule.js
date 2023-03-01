export const LINK_RULE_TYPES = {
  alarm: 'alarm',
  entity: 'entity',
};

export const LINK_RULE_DEFAULT_ALARM_SOURCE_CODE = `function generate(alarms) {
  return [
    {
      label: "",
      category: "",
      icon_name: "",
      url: ""
    }
  ];
}`;

export const LINK_RULE_DEFAULT_ENTITY_SOURCE_CODE = `function generate(entities) {
  return [
    {
      label: "",
      category: "",
      icon_name: "",
      url: ""
    }
  ];
}`;

export const LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES = {
  [LINK_RULE_TYPES.alarm]: LINK_RULE_DEFAULT_ALARM_SOURCE_CODE,
  [LINK_RULE_TYPES.entity]: LINK_RULE_DEFAULT_ENTITY_SOURCE_CODE,
};
