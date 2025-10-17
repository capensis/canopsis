import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the view Vuex store module
 *
 * @returns {Object} Store module hooks for view namespace
 * @property {import('vuex').Store} store - Vuex store instance
 * @property {import('vuex').Module} module - View module instance
 * @property {Function} useGetters - Hook for accessing view module getters
 * @property {Function} useActions - Hook for accessing view module actions
 */
const useViewStoreModule = () => useStoreModuleHooks('view');

/**
 * Hook for managing view-related operations and state
 * Replaces the entitiesViewMixin functionality with Composition API
 *
 * @returns {Object} An object containing view getters, actions, and methods
 * @property {Function} createView - Creates a new view
 * @property {Function} updateView - Updates an existing view
 * @property {Function} updateViewsPositions - Updates view positions
 * @property {Function} updateViewWithoutStore - Updates view without store
 * @property {Function} removeView - Removes a view
 * @property {Function} copyView - Copies a view
 * @property {Function} exportViewsWithoutStore - Exports views without store
 * @property {Function} importViewsWithoutStore - Imports views without store
 * @property {Function} createViewWithPopup - Creates view with success popup
 * @property {Function} updateViewWithPopup - Updates view with success popup
 * @property {Function} copyViewWithPopup - Copies view with success popup
 * @property {Function} removeViewWithPopup - Removes view with success popup
 */
export const useView = () => {
  const { t } = useI18n();
  const popups = usePopups();

  const { useActions } = useViewStoreModule();

  const actions = useActions({
    createView: 'createView',
    updateView: 'updateView',
    updateViewsPositions: 'updateViewPositions',
    updateViewWithoutStore: 'updateWithoutStoreView',
    removeView: 'removeView',
    copyView: 'copyView',
    exportViewsWithoutStore: 'exportViewWithoutStore',
    importViewsWithoutStore: 'importViewWithoutStore',
  });

  /**
   * Creates a new view and shows a success popup message
   * @async
   * @param {Object} params - The parameters object
   * @param {Object} params.data - The view data to create
   * @returns {Promise<void>}
   */
  const createViewWithPopup = async ({ data }) => {
    try {
      await actions.createView({ data });

      popups.success({ text: t('modals.view.success.create') });
    } catch (err) {
      popups.error({ text: t('modals.view.fail.create') });

      throw err;
    }
  };

  /**
   * Updates an existing view and shows a success popup message
   * @async
   * @param {Object} params - The parameters object
   * @param {string|number} params.id - The ID of the view to update
   * @param {Object} params.data - The view data to update
   * @returns {Promise<void>}
   */
  const updateViewWithPopup = async ({ id, data }) => {
    try {
      await actions.updateView({ id, data });

      popups.success({ text: t('modals.view.success.edit') });
    } catch (err) {
      popups.error({ text: t('modals.view.fail.edit') });

      throw err;
    }
  };

  /**
   * Copies an existing view and shows a success popup message
   * @async
   * @param {Object} params - The parameters object
   * @param {string|number} params.id - The ID of the view to copy
   * @param {Object} params.data - The view data for copying
   * @returns {Promise<void>}
   */
  const copyViewWithPopup = async ({ id, data }) => {
    try {
      await actions.copyView({ id, data });

      popups.success({ text: t('modals.view.success.duplicate') });
    } catch (err) {
      popups.error({ text: t('modals.view.fail.duplicate') });

      throw err;
    }
  };

  /**
   * Removes a view and shows a success popup message
   * @async
   * @param {Object} params - The parameters object
   * @param {string|number} params.id - The ID of the view to remove
   * @returns {Promise<void>}
   */
  const removeViewWithPopup = async ({ id }) => {
    try {
      await actions.removeView({ id });

      popups.success({ text: t('modals.view.success.delete') });
    } catch (err) {
      popups.error({ text: t('modals.view.fail.delete') });

      throw err;
    }
  };

  return {
    ...actions,
    createViewWithPopup,
    updateViewWithPopup,
    copyViewWithPopup,
    removeViewWithPopup,
  };
};
