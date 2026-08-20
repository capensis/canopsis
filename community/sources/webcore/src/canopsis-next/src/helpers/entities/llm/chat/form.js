import { LLM_SOCKET_CONTEXTS, PATTERNS_FIELDS, STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD } from '@/constants';

/**
 * Builds a single-key `patterns` map for corporate shared filters: `{ [field]: form }` when the block has groups.
 *
 * @param {Object} [form={}] - Single pattern editor block (`groups` array).
 * @param {string} field - API pattern field name (e.g. `PATTERNS_FIELDS.alarm`).
 * @returns {Object<string, Object>} Map with at most one entry, or `{}` if there are no groups.
 */
export const aiChatFormToPatternsField = (form = {}, field) => (
  form.groups.length ? { [field]: form } : {}
);

/**
 * Collects every pattern block on a filter-style form that uses `PATTERNS_FIELDS` keys and non-empty `groups`.
 *
 * @param {Object} [form={}] - Form whose keys may include alarm/entity/pbehavior/etc. pattern fields.
 * @returns {Object<string, Object>} Subset of `form` pattern fields that have at least one group.
 */
export const aiChatFormToPatternsDefault = (form = {}) => {
  const patternsForm = form.patterns ?? form;

  return Object.values(PATTERNS_FIELDS).reduce((acc, field) => {
    if (patternsForm[field]?.groups?.length) {
      acc[field] = patternsForm[field];
    }

    return acc;
  }, {});
};

/**
 * Maps a list of form rows (e.g. scenario actions) to a list of per-row pattern maps for the AI sidebar.
 *
 * @param {Object[]} [form=[]] - Array of form fragments, each passed through `aiChatFormToPatternsDefault`.
 * @returns {Object} Same length as `form`; each element is a `patterns`-shaped object.
 */
export const aiChatFormToPatternsItems = (form = []) => form.reduce((acc, item) => {
  acc[item.key] = aiChatFormToPatternsDefault(item);

  return acc;
}, {});

/**
 * Builds the AI chat sidebar `patterns` object for state-settings modals: one entry per block (`entity_pattern`
 * and inherited dependency pattern), each reduced via `aiChatFormToPatternsDefault` so only non-empty `groups`
 * are kept.
 *
 * @param {Object} [form={}] - State-setting form containing `entity_pattern` and
 *   `STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD` editors.
 * @returns {Object<string, Object>} Two keyed pattern maps aligned with `PATTERNS_FIELDS.entity` and
 *   `STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD`.
 */
export const aiChatFormToPatternsStateSettings = (form = {}) => ({
  [PATTERNS_FIELDS.entity]: aiChatFormToPatternsDefault(form[PATTERNS_FIELDS.entity]),
  [STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD]: aiChatFormToPatternsDefault(
    form[STATE_SETTINGS_INHERITED_ENTITY_PATTERN_FIELD],
  ),
});

/**
 * Normalizes modal form data into the `patterns` payload shape expected by the AI chat sidebar, by socket context.
 *
 * Scenario context returns an array of pattern maps; corporate pattern contexts use a single keyed block;
 * any other `context` uses `aiChatFormToPatternsDefault(form)`.
 *
 * @param {Object} [params={}]
 * @param {Object|Object[]} params.form - Filter object, or array of action rows for scenario.
 * @param {string} [params.field] - Pattern field name when using single-field corporate contexts.
 * @param {string} [params.context] - `LLM_SOCKET_CONTEXTS` value driving the mapping strategy.
 * @returns {Object|Object[]} `patterns` value for sidebar config (object or array of objects).
 */
export const aiChatFormToPatterns = ({
  form,
  field,
  context,
} = {}) => {
  const method = {
    [LLM_SOCKET_CONTEXTS.scenario]: () => aiChatFormToPatternsItems(form),
    [LLM_SOCKET_CONTEXTS.stateSettings]: () => aiChatFormToPatternsStateSettings(form),
    [LLM_SOCKET_CONTEXTS.corporateAlarmPattern]: () => aiChatFormToPatternsField(form, field),
    [LLM_SOCKET_CONTEXTS.corporateEntityPattern]: () => aiChatFormToPatternsField(form, field),
    [LLM_SOCKET_CONTEXTS.corporatePbehaviorPattern]: () => aiChatFormToPatternsField(form, field),
    [LLM_SOCKET_CONTEXTS.corporateWeatherServicePattern]: () => aiChatFormToPatternsField(form, field),
  }[context];

  return (method ?? (() => aiChatFormToPatternsDefault(form)))();
};
