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

/**
 * Checks if an active broadcast message should be visible on the current route
 *
 * @param {Object} message - Active broadcast message
 * @param {boolean} [message.maintenance] - Whether the message is linked to maintenance mode
 * @param {string[]} [message.views] - Configured message views
 * @param {string} routeView - The current route's broadcast message view
 * @param {string} routeId - The current route's ID parameter
 * @param {Object|null} currentView - The current view object from store
 * @returns {boolean} Whether the message should be visible on the current route
 */
export const isActiveBroadcastMessageVisibleOnRoute = (
  { maintenance, views: messageViews },
  routeView,
  routeId,
  currentView,
) => maintenance || (messageViews || []).some(
  messageView => isBroadcastMessageViewMatchingRoute(messageView, routeView, routeId, currentView),
);
