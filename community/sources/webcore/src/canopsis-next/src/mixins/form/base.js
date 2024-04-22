import formComputedPropertiesMixin, {
  modelPropKeyComputed,
  modelEventKeyComputed,
} from './internal/computed-properties';

/**
 * @mixin Form mixin
 * @deprecated Should be used useModelField
 */
export const formBaseMixin = {
  mixins: [formComputedPropertiesMixin],
  methods: {
    /**
     * Update full model
     *
     * @param {*} model
     * @return {Array|Object}
     */
    updateModel(model) {
      this.$emit(this[modelEventKeyComputed], model);

      return model;
    },
  },
};

export {
  modelPropKeyComputed,
  modelEventKeyComputed,
};

export default formBaseMixin;
