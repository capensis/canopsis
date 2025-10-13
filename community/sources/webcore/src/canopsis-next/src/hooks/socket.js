import { onMounted, onBeforeUnmount } from 'vue';

import { useComponentInstance } from './vue';

/**
 * Provides access to the global $socket instance.
 *
 * @returns {Object} The socket instance.
 */
export const useSocket = () => {
  const instance = useComponentInstance();

  return instance.$socket;
};

/**
 * Manages joining and leaving a socket room with a listener on component mount/unmount.
 *
 * @param {string} room - The name of the socket room to join.
 * @param {Function} listener - The event listener function to add and remove.
 * @returns {Object} The injected socket instance.
 */
export const useSocketRoom = (room, listener) => {
  const socket = useSocket();

  onMounted(() => socket.join(room).addListener(listener));
  onBeforeUnmount(() => socket.leave(room).removeListener(listener));

  return socket;
};
