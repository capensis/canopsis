import isEqual from 'lodash/isEqual';
import { onBeforeUnmount, watch } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { mapIds } from '@/helpers/array';

import { useSocket } from '@/hooks/socket';

/**
 * Subscribes to live anomaly monitored connector states over websocket.
 *
 * Joins `SOCKET_ROOMS.anomalyMonitoredConnectorStates` with the current connector ids,
 * updates `connectors` and `meta.total_count` when the room pushes new data, and
 * reconnects whenever the connector id list changes. Leaves the room on unmount.
 *
 * @param {Object} options
 * @param {import('vue').Ref<Array>} options.connectors - Reactive list of monitored connectors.
 * @param {import('vue').Ref<{ total_count: number }>} options.meta - Pagination meta updated from socket.
 */
export const useHealthcheckConnectorsSocket = ({ connectors, meta }) => {
  const socket = useSocket();
  const connectorsRef = connectors;
  const metaRef = meta;

  /**
   * Applies incoming socket payload to `connectors` and `meta.total_count`.
   *
   * @param {Object} [data]
   * @param {Array} [data.data] - Next connectors list.
   * @param {number} [data.total_count] - Next total count.
   */
  const socketListener = (data = {}) => {
    const nextConnectors = data.data ?? [];

    if (!isEqual(connectorsRef.value, nextConnectors)) {
      connectorsRef.value = nextConnectors;
    }

    if (metaRef.value.total_count !== data.total_count) {
      metaRef.value = {
        ...metaRef.value,
        total_count: data.total_count,
      };
    }
  };

  /**
   * Joins the anomaly monitored connector states room for the given ids.
   *
   * @param {Array<string|number>} [ids] - Connector ids to subscribe to.
   */
  const joinToSocketRoom = (ids = []) => {
    if (!socket) {
      return;
    }

    socket.join(SOCKET_ROOMS.anomalyMonitoredConnectorStates, { ids }, true).addListener(socketListener);
  };

  /**
   * Leaves the anomaly monitored connector states room and removes the listener.
   */
  const leaveFromSocketRoom = () => {
    if (!socket) {
      return;
    }

    socket.leave(SOCKET_ROOMS.anomalyMonitoredConnectorStates).removeListener(socketListener);
  };

  /**
   * Rejoins the room with a new set of connector ids.
   *
   * @param {Array<string|number>} [ids] - Connector ids to subscribe to.
   */
  const reconnectToSocketRoom = (ids = []) => {
    leaveFromSocketRoom();
    joinToSocketRoom(ids);
  };

  watch(
    connectorsRef,
    (items, prevItems) => {
      const newIds = mapIds(items, 'id');
      const oldIds = mapIds(prevItems ?? [], 'id');

      if (isEqual(newIds, oldIds)) {
        return;
      }

      reconnectToSocketRoom(newIds);
    },
    { immediate: true },
  );

  onBeforeUnmount(leaveFromSocketRoom);
};
