import formComputedPropertiesMixin, { modelPropKeyComputed, modelEventKeyComputed } from './internal/computed-properties';

/**
 * @mixin Form mixin
 */
export default {
  mixins: [formComputedPropertiesMixin],
  methods: {
    /**
     * Update full model
     *
     * @param {*} model
     */
    updateModel(model) {
      this.$emit(this[modelEventKeyComputed], model);
    },
  },
};

export {
  modelPropKeyComputed,
  modelEventKeyComputed,
};
