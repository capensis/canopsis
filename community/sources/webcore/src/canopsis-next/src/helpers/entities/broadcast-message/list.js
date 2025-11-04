/**
 * Checks if a broadcast message view matches the current route context
 *
 * @param {string} activeMessageView - The view identifier from the broadcast message
 * @param {string} routeView - The current route's broadcast message view
 * @param {string} routeId - The current route's ID parameter
 * @param {Object|null} currentView - The current view object from store
 * @returns {boolean} Whether the message view matches the current route
 */
export const isBroadcastMessageViewMatchingRoute = (activeMessageView, routeView, routeId, currentView) => (
  [routeView, routeId, currentView?.group?._id].includes(activeMessageView)
);
