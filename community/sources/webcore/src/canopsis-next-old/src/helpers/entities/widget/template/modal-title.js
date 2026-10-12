/**
 * Returns modal title for create/edit widget template modal.
 *
 * @param {Object} params
 * @param {Function} params.t - i18n translate function
 * @param {Function} params.te - i18n key exists function
 * @param {string} [params.type] - Widget template type
 * @param {boolean} [params.isEdit=false] - Whether the template is being edited
 * @returns {string}
 */
export const getWidgetTemplateModalTitle = ({ t, te, type, isEdit = false }) => {
  const baseKey = isEdit
    ? 'modals.createWidgetTemplate.edit.title'
    : 'modals.createWidgetTemplate.create.title';
  const base = t(baseKey);
  const typeKey = `modals.createWidgetTemplate.types.${type}`;

  if (type && te(typeKey)) {
    return `${base} - ${t(typeKey)}`;
  }

  return base;
};
