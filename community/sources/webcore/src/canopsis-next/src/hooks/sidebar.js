import { useSidebarStore } from '@/stores/modules/sidebar';

/**
 * Sidebar store facade.
 *
 * @returns {Object}
 */
export const useSidebar = () => {
  const sidebarStore = useSidebarStore();

  return {
    moduleName: 'sidebar',
    show: sidebarStore.show,
    hide: sidebarStore.hide,
    updateConfig: sidebarStore.updateConfig,
  };
};
