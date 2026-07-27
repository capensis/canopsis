import { inject, onBeforeUnmount, onMounted } from 'vue';

export const useSocket = () => inject('$socket', null);

export const useSocketRoom = ({ room, payload, authNeeded, listener } = {}) => {
  const socket = useSocket();
  let socketRoom;

  onMounted(() => {
    if (!socket || !room) {
      return;
    }

    socketRoom = socket.join(room, payload, authNeeded);
    socketRoom.addListener(listener);
  });

  onBeforeUnmount(() => {
    if (!socket || !room) {
      return;
    }

    socket.leave(room).removeListener(listener);
  });

  return {
    socketRoom,
  };
};
