export const LINK_RULE_TYPES = {
  alarm: 'alarm',
  entity: 'entity',
};

export const LINK_RULE_DEFAULT_ALARM_SOURCE_CODE = `function generate(alarms, user, env) {
  for (const alarm of alarms) {

  }

  return [
    {
      label: "",
      category: "",
      icon_name: "",
      url: "",
      action: 0,
      single: true
    }
  ];
}`;

export const LINK_RULE_DEFAULT_ENTITY_SOURCE_CODE = `function generate(entities, user, env) {
  for (const entity of entities) {

  }

  return [
    {
      label: "",
      category: "",
      icon_name: "",
      url: "",
      action: 0,
      single: true
    }
  ];
}`;

export const LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES = {
  [LINK_RULE_TYPES.alarm]: LINK_RULE_DEFAULT_ALARM_SOURCE_CODE,
  [LINK_RULE_TYPES.entity]: LINK_RULE_DEFAULT_ENTITY_SOURCE_CODE,
};

export const DEFAULT_LINKS_INLINE_COUNT = 3;

export const LINK_RULE_ACTIONS = {
  open: 'open',
  copy: 'copy',
};

export const LINK_RULE_ADVANCED_COMPLETIONS = {
  alarm: {
    Value: {
      Component: '',
      Connector: '',
    },
  },
  entity: {
    Infos: {
      '%infos_name%': {
        Value: '',
      },
    },
  },
  user: {
    Username: '',
  },
  env: {
    var1: '',
  },
};
