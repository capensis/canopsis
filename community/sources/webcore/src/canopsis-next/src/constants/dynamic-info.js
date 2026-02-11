export const DYNAMIC_INFO_PATTERNS_OPERATORS = ['>=', '>', '<', '<=', 'regex_match'];

export const DYNAMIC_INFO_INFORMATION_TYPES = {
  setToInfo: 'set_to_info',
  setToInfoFromTemplate: 'set_to_info_from_template',
  copyToInfo: 'copy_to_info',
};

export const DYNAMIC_INFO_FIELDS = {
  id: '_id',
  name: 'name',
  author: 'author.display_name',
  enabled: 'enabled',
  created: 'created',
  updated: 'updated',
};

export const DYNAMIC_INFO_FIELDS_TO_LABELS_KEYS = {
  [DYNAMIC_INFO_FIELDS.id]: 'common.id',
  [DYNAMIC_INFO_FIELDS.name]: 'common.name',
  [DYNAMIC_INFO_FIELDS.author]: 'common.author',
  [DYNAMIC_INFO_FIELDS.enabled]: 'common.enabled',
  [DYNAMIC_INFO_FIELDS.created]: 'common.created',
  [DYNAMIC_INFO_FIELDS.updated]: 'common.updated',
};
