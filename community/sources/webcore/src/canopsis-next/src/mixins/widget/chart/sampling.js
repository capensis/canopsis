export const widgetSamplingFilterMixin = {
  methods: {
    updateSampling(sampling) {
      this.updateContentInUserPreference({ sampling });

      this.query = {
        ...this.query,
        sampling,
      };
    },
  },
};
