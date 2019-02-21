export default {
  methods: {
    updateRecordsPerPage(limit) {
      this.updateLockedQuery({
        id: this.widget._id,
        query: { limit },
      });
    },
  },
};
