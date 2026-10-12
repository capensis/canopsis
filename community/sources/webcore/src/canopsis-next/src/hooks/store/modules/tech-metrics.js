import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Custom hook to use the tech metrics store module.
 *
 * This hook utilizes the `useStoreModuleHooks` function to access the tech metrics store module
 * and provides actions for creating and fetching tech metrics exports.
 *
 * @returns {Object} An object containing the actions for creating and fetching tech metrics exports.
 * @property {Function} createTechMetricsExport - Action to create a tech metrics export.
 * @property {Function} fetchTechMetricsExport - Action to fetch a tech metrics export.
 */
const useTechMetricsStoreModule = () => useStoreModuleHooks('techMetrics');

/**
 * Custom hook to use tech metrics actions.
 *
 * This hook provides actions for creating and fetching tech metrics exports, as well as fetching and updating tech
 * metrics settings,
 * by utilizing the `useTechMetricsStoreModule` hook.
 *
 * @returns {Object} An object containing the actions for creating and fetching tech metrics exports, and fetching and
 * updating tech metrics settings.
 * @property {Function} createTechMetricsExport - Action to create a tech metrics export.
 * @property {Function} fetchTechMetricsExport - Action to fetch a tech metrics export.
 * @property {Function} fetchTechMetricsSettings - Action to fetch tech metrics settings.
 * @property {Function} updateTechMetricsSettings - Action to update tech metrics settings.
 */
export const useTechMetrics = () => {
  const { useActions } = useTechMetricsStoreModule();

  const actions = useActions({
    createTechMetricsExport: 'createTechMetricsExport',
    fetchTechMetricsExport: 'fetchTechMetricsExport',
    fetchTechMetricsSettings: 'fetchTechMetricsSettings',
    updateTechMetricsSettings: 'updateTechMetricsSettings',
  });

  return {
    ...actions,
  };
};
