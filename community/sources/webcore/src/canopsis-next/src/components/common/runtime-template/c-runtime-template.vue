<template>
  <component
    :is="component"
    v-bind="props"
    v-on="$listeners"
  />
</template>

<script>
export default {
  props: {
    template: {
      type: String,
      required: false,
    },
    parent: {
      type: Object,
      required: false,
    },
    templateProps: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    parentNode() {
      return this.parent || this.$parent;
    },

    parentOptions() {
      return this.parentNode.$options;
    },

    methodProps() {
      return Object.entries(this.parentOptions.methods ?? {}).reduce((acc, [key]) => {
        acc[key] = this.parentNode[key];

        return acc;
      }, {});
    },

    props() {
      return Object.entries({
        ...this.parentNode.$data,
        ...this.parentNode.$props,
        ...this.methodProps,
        ...this.templateProps,
      }).reduce((acc, [key, value]) => {
        if (!this.isDeclaredByMixin(key)) {
          acc[key] = value;
        }

        return acc;
      }, {});
    },

    propsTypes() {
      return Object.keys({
        ...this.parentNode.$data,
        ...this.parentOptions.props,
        ...this.parentOptions.methods,
        ...this.templateProps,
      }).filter(key => !this.isDeclaredByMixin(key));
    },

    component() {
      return {
        template: this.template || '<div></div>',
        components: this.parentOptions.components,
        computed: this.parentOptions.computed,
        props: this.propsTypes,
        // eslint-disable-next-line no-underscore-dangle
        provide: this.parentNode._provided,
      };
    },
  },
  methods: {
    isDeclaredByMixin(key) {
      return this.$data[key] || this.$props[key] || this.$options.computed[key] || this.$options.methods[key];
    },
  },
};
</script>
