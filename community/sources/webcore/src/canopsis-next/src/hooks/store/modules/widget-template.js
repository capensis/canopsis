import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Custom hook to access the `widgetTemplate` store module.
 *
 * @returns {Object} An object containing the store, module, useGetters, and useActions
 *                   functions for the `widgetTemplate` namespace.
 */
export const useWidgetTemplateStoreModule = () => useStoreModuleHooks('widgetTemplate');

/**
 * Custom hook to interact with the `widgetTemplate` store module.
 *
 * This hook provides access to the `widgetTemplate` module's getters and actions,
 * allowing for easy retrieval and manipulation of widget template data.
 *
 * @returns {Object} An object containing:
 * - `widgetTemplates`: Getter for widget templates items.
 * - `widgetTemplatesMeta`: Getter for widget templates meta information.
 * - `widgetTemplatesPending`: Getter for widget templates pending state.
 * - `fetchWidgetTemplatesList`: Action to fetch the list of widget templates.
 * - `fetchWidgetTemplatesListWithPreviousParams`: Action to fetch widget templates with previous parameters.
 * - `createWidgetTemplate`: Action to create a new widget template.
 * - `updateWidgetTemplate`: Action to update an existing widget template.
 * - `removeWidgetTemplate`: Action to remove a widget template.
 */
export const useWidgetTemplate = () => {
  const { useGetters, useActions } = useWidgetTemplateStoreModule();

  const getters = useGetters({
    widgetTemplates: 'items',
    widgetTemplatesMeta: 'meta',
    widgetTemplatesPending: 'pending',
  });

  const actions = useActions({
    fetchWidgetTemplatesList: 'fetchList',
    fetchWidgetTemplatesListWithPreviousParams: 'fetchListWithPreviousParams',
    createWidgetTemplate: 'create',
    updateWidgetTemplate: 'update',
    removeWidgetTemplate: 'remove',
  });

  return {
    ...getters,
    ...actions,
  };
};
