export const isBroadcastMessageViewMatchingRoute = (activeMessageView, routeView, routeId, currentView) => (
  [routeView, routeId, currentView?.group?._id].includes(activeMessageView)
);

export const isActiveBroadcastMessageVisibleOnRoute = (
  { maintenance, views: messageViews },
  routeView,
  routeId,
  currentView,
) => maintenance || (messageViews || []).some(
  messageView => isBroadcastMessageViewMatchingRoute(messageView, routeView, routeId, currentView),
);
