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
 * @param {Object} options - Socket room configuration options
 * @param {string} options.room - The name of the socket room to join
 * @param {*} [options.data=null] - Optional data to pass when joining the room
 * @param {boolean} [options.needAuth=true] - Whether authentication is needed for the room
 * @param {Function} options.listener - The event listener function to add and remove
 * @returns {Object} The injected socket instance
 */
export const useSocketRoom = ({ room, data = null, needAuth = true, listener }) => {
  const socket = useSocket();

  onMounted(() => socket.join(room, data, needAuth).addListener(listener));
  onBeforeUnmount(() => socket.leave(room).removeListener(listener));

  return socket;
};
