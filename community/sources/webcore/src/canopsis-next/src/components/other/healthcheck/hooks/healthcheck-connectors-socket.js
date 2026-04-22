import { isEqual } from 'lodash';
import { watch, set, onBeforeUnmount } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { mapIds } from '@/helpers/array';

import { useSocket } from '@/hooks/socket';

/**
 * Subscribes to the anomaly monitored connector states socket room, applies push
 * updates to the paged `connectors` list and `meta.total_count`, and re-subscribes
 * when the set of ids on the current page changes. Cleans up the listener on unmount.
 *
 * @param {Object} params
 * @param {import('vue').Ref<Array>} params.connectors - Current page of connector state rows.
 * @param {import('vue').Ref<Object>} params.meta - List metadata (e.g. `total_count`).
 */
export const useHealthcheckConnectorsSocket = ({ connectors, meta }) => {
  const socket = useSocket();

  /**
   * Applies live anomaly monitored connector state updates from the socket to the
   * paged list and total count, skipping no-op per-index changes.
   *
   * @param {Object} data - Payload from the anomaly monitored connector states room.
   * @param {Array} data.data - Connectors in list order, aligned with the current page slice.
   * @param {number} data.total_count - Total number of connectors across all pages.
   */
  const socketListener = (data) => {
    data.data.forEach((connector, index) => {
      if (isEqual(connector, connectors.value[index])) {
        return;
      }

      set(connectors.value, index, connector);
    });

    if (meta.value.total_count !== data.total_count) {
      set(meta.value, 'total_count', data.total_count);
    }
  };

  /**
   * Subscribes to the anomaly monitored connector states room for the given connector ids
   * and registers the {@link socketListener} to receive push updates.
   *
   * @param {string[]} [ids=[]] - Connector ids to scope the subscription; empty subscribes without filtering.
   * @returns {Object} The socket instance (chainable from join).
   */
  const joinToSocketRoom = (ids = []) => socket
    .join(SOCKET_ROOMS.anomalyMonitoredConnectorStates, { ids }, true)
    .addListener(socketListener);

  /**
   * Unsubscribes from the anomaly monitored connector states room and detaches
   * {@link socketListener}.
   *
   * @returns {Object} The socket instance (chainable from leave).
   */
  const leaveFromSocketRoom = () => socket
    .leave(SOCKET_ROOMS.anomalyMonitoredConnectorStates)
    .removeListener(socketListener);

  /**
   * Re-subscribes the socket with a fresh set of connector ids (e.g. after the paged
   * list or filters change). Leaves the room first, then joins again.
   *
   * @param {string[]} [ids=[]] - Connector ids passed to {@link joinToSocketRoom}.
   */
  const reconnectToSocketRoom = (ids = []) => {
    leaveFromSocketRoom();
    joinToSocketRoom(ids);
  };

  watch(connectors, (items, prevItems) => {
    const newIds = mapIds(items, 'id');
    const oldIds = mapIds(prevItems ?? [], 'id');

    if (isEqual(newIds, oldIds)) {
      return;
    }

    reconnectToSocketRoom(newIds);
  }, { immediate: true });

  onBeforeUnmount(leaveFromSocketRoom);
};
