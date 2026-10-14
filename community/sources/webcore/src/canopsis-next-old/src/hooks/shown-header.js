import { computed } from 'vue';
import { useRoute } from 'vue-router/composables';

import { useInfo } from './store/modules/info';
import { useActiveView } from './store/modules/active-view';

/**
 * Composable hook for determining header visibility based on kiosk mode and route configuration
 *
 * @returns {Object}
 * @returns {ComputedRef<boolean>} returns.showHeaderOnKioskMode - Whether to show header in kiosk mode
 * @returns {ComputedRef<boolean>} returns.activeViewIsKioskScreenMode - Whether the active view is in kiosk screen mode
 * @returns {ComputedRef<boolean>} returns.shownHeader - Whether the header should be shown
 */
export const useShownHeader = () => {
  const { showHeaderOnKioskMode } = useInfo();
  const { activeViewIsKioskScreenMode } = useActiveView();

  const route = useRoute();
  const shownHeader = computed(() => (
    activeViewIsKioskScreenMode.value
      ? showHeaderOnKioskMode.value
      : !route?.meta?.hideHeader
  ));

  return {
    showHeaderOnKioskMode,
    activeViewIsKioskScreenMode,
    shownHeader,
  };
};
