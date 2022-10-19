<template lang="pug">
  component(:is="component", v-bind="props")
</template>

<script>
export default {
  props: {
    template: {
      type: String,
      required: false,
    },
  },
  computed: {
    parent() {
      return this.$parent;
    },

    parentOptions() {
      return this.parent.$options;
    },

    methodProps() {
      return Object.entries(this.parentOptions.methods).reduce((acc, [key]) => {
        Object.defineProperty(acc, key, Object.getOwnPropertyDescriptor(this.$parent, key));

        return acc;
      }, {});
    },

    props() {
      return Object.entries({
        ...this.parent.$data,
        ...this.parent.$props,
        ...this.methodProps,
      }).reduce((acc, [key, value]) => {
        if (!this.isDeclaredByMixin(key)) {
          acc[key] = value;
        }

        return acc;
      }, {});
    },

    propsTypes() {
      return [
        ...Object.keys(this.parent.$data),
        ...Object.keys(this.parentOptions.props),
        ...Object.keys(this.parentOptions.methods),
      ].filter(key => !this.isDeclaredByMixin(key));
    },

    component() {
      return {
        template: this.template || '<div></div>',
        components: this.parentOptions.components,
        computed: this.parentOptions.computed,
        props: this.propsTypes,
        // eslint-disable-next-line no-underscore-dangle
        provide: this.parent._provided,
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
