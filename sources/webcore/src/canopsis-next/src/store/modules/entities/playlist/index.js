export default {
  namespaced: true,
  actions: {
    fetchItemWithoutStore() {
      return {
        _id: 'asd',
        name: 'Playlist #1',
        fullscreen: true,
        interval: {
          value: 10,
          unit: 'm',
        },
        tabs: [
          '875df4c2-027b-4549-8add-e20ed7ff7d4f', // Alarm default
          'view-tab_5a339b3a-0611-4d4c-b307-dc1b92aeb27d', // Meteo technic
          'view-tab_c02ae48e-7f0a-4ba4-9215-ba5662e1550c', // Meteo correct
        ],
      };
    },
  },
};
